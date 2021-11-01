package eks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

// Client stores information about the client config
type Client struct {
	config    *clientcmdapi.Config
	rawConfig *restclient.Config
}

// AWSKubeProvider encapsulates kubernetes client building logic
type AWSKubeProvider struct {
	AWSProvider     api.AWSProvider
	ClusterProvider ClusterProvider // TODO: resolve circular reference, only required for GetUsername()
}

// NewClient creates a new client config by embedding the STS token
func (p *AWSKubeProvider) NewClient(spec *api.ClusterConfig) (kubewrapper.ClientInterface, error) {
	config := kubeconfig.NewForUser(spec, p.ClusterProvider.GetUsername())
	client := &Client{
		config: config,
	}
	return client.new(spec, p.AWSProvider.STS())
}

// NewRawClient creates a new raw REST client in one go with an embedded STS token
func (p *AWSKubeProvider) NewRawClient(spec *api.ClusterConfig) (kubewrapper.RawClientInterface, error) {
	client, clientSet, err := p.newClientSetWithEmbeddedToken(spec)
	if err != nil {
		return nil, err
	}

	return kubewrapper.NewRawClient(clientSet, client.rawConfig)
}

// NewStdClientSet creates a new API client in one go with an embedded STS token, this is most commonly used option
func (p *AWSKubeProvider) NewStdClientSet(spec *api.ClusterConfig) (kubewrapper.Interface, error) {
	_, clientSet, err := p.newClientSetWithEmbeddedToken(spec)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func (p *AWSKubeProvider) newClientSetWithEmbeddedToken(spec *api.ClusterConfig) (*Client, kubernetes.Interface, error) {
	client, err := p.NewClient(spec)
	if err != nil {
		return nil, nil, errors.Wrap(err, "creating Kubernetes client config with embedded token")
	}

	clientSet, err := client.NewClientSet()
	if err != nil {
		return nil, nil, errors.Wrap(err, "creating Kubernetes client")
	}

	return client.(*Client), clientSet, nil
}

// GetUsername extracts the username part from the IAM role ARN
func (c *ClusterProviderImpl) GetUsername() string {
	usernameParts := strings.Split(c.status.iamRoleARN, "/")
	if len(usernameParts) > 1 {
		return usernameParts[len(usernameParts)-1]
	}
	return "iam-root-account"
}

func (c *Client) new(spec *api.ClusterConfig, stsClient stsiface.STSAPI) (*Client, error) {
	if err := c.useEmbeddedToken(spec, stsClient); err != nil {
		return nil, err
	}

	rawConfig, err := clientcmd.NewDefaultClientConfig(*c.config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client configuration from client config")
	}
	c.rawConfig = rawConfig

	return c, nil
}

func (c *Client) useEmbeddedToken(spec *api.ClusterConfig, stsclient stsiface.STSAPI) error {
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return errors.Wrap(err, "could not get token generator")
	}

	tok, err := gen.GetWithSTS(spec.Metadata.Name, stsclient)
	if err != nil {
		return errors.Wrap(err, "could not get token")
	}

	c.config.AuthInfos[c.config.CurrentContext].Token = tok.Token
	return nil
}

// NewClientSet creates a new API client
func (c *Client) NewClientSet() (kubernetes.Interface, error) {
	client, err := kubernetes.NewForConfig(c.rawConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client")
	}
	return client, nil
}

func (c *Client) Config() *clientcmdapi.Config {
	return c.config
}

// ServerVersion will use discovery API to fetch version of Kubernetes control plane
func (c *ClusterProviderImpl) ServerVersion(rawClient kubewrapper.RawClientInterface) (string, error) {
	return rawClient.ServerVersion()
}

// UpdateAuthConfigMap creates or adds a nodegroup IAM role in the auth ConfigMap for the given nodegroup.
func (c *ClusterProviderImpl) UpdateAuthConfigMap(nodeGroups []*api.NodeGroup, clientSet kubernetes.Interface) error {
	for _, ng := range nodeGroups {
		// authorise nodes to join
		if err := authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
			return err
		}

		// wait for nodes to join
		if err := c.WaitForNodes(clientSet, ng); err != nil {
			return err
		}
	}
	return nil
}

// WaitForNodes waits till the nodes are ready
func (c *ClusterProviderImpl) WaitForNodes(clientSet kubernetes.Interface, ng KubeNodeGroup) error {
	minSize := ng.Size()
	if minSize == 0 {
		return nil
	}
	timer := time.After(c.awsProvider.WaitTimeout())
	timeout := false
	readyNodes := sets.NewString()
	watcher, err := clientSet.CoreV1().Nodes().Watch(context.TODO(), ng.ListOptions())
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}

	counter, err := getNodes(clientSet, ng)
	if err != nil {
		return errors.Wrap(err, "listing nodes")
	}

	logger.Info("waiting for at least %d node(s) to become ready in %q", minSize, ng.NameString())
	for !timeout && counter < minSize {
		select {
		case event := <-watcher.ResultChan():
			logger.Debug("event = %#v", event)
			if event.Object != nil && event.Type != watch.Deleted {
				if node, ok := event.Object.(*corev1.Node); ok {
					if isNodeReady(node) {
						readyNodes.Insert(node.Name)
						counter = readyNodes.Len()
						logger.Debug("node %q is ready in %q", node.Name, ng.NameString())
					} else {
						logger.Debug("node %q seen in %q, but not ready yet", node.Name, ng.NameString())
						logger.Debug("node = %#v", *node)
					}
				}
			}
		case <-timer:
			timeout = true
		}
	}
	watcher.Stop()
	if timeout {
		return fmt.Errorf("timed out (after %s) waiting for at least %d nodes to join the cluster and become ready in %q", c.awsProvider.WaitTimeout(), minSize, ng.NameString())
	}

	if _, err = getNodes(clientSet, ng); err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}

	return nil
}
