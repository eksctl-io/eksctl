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

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/credentials"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

// Client stores information about the client config
type Client struct {
	Config    *clientcmdapi.Config
	Generator TokenGenerator

	rawConfig *restclient.Config
}

// NewClient creates a new client config by embedding the STS token
func (c *ClusterProvider) NewClient(spec *api.ClusterConfig) (*Client, error) {
	config := kubeconfig.NewForUser(spec, c.GetUsername())
	generator := NewGenerator(c.Provider.STSPresigner(), &credentials.RealClock{})
	client := &Client{
		Config:    config,
		Generator: generator,
	}
	return client.new(spec)
}

// GetUsername extracts the username part from the IAM role ARN
func (c *ClusterProvider) GetUsername() string {
	usernameParts := strings.Split(c.Status.iamRoleARN, "/")
	if len(usernameParts) > 1 {
		return usernameParts[len(usernameParts)-1]
	}
	return "iam-root-account"
}

func (c *Client) new(spec *api.ClusterConfig) (*Client, error) {
	if err := c.useEmbeddedToken(spec); err != nil {
		return nil, err
	}

	rawConfig, err := clientcmd.NewDefaultClientConfig(*c.Config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client configuration from client config")
	}

	c.rawConfig = rawConfig
	c.rawConfig.QPS = float32(25)
	c.rawConfig.Burst = int(c.rawConfig.QPS) * 2

	return c, nil
}

func (c *Client) useEmbeddedToken(spec *api.ClusterConfig) error {
	tok, err := c.Generator.GetWithSTS(context.TODO(), spec.Metadata.Name)
	if err != nil {
		return errors.Wrap(err, "could not get token")
	}

	c.Config.AuthInfos[c.Config.CurrentContext].Token = tok.Token
	return nil
}

// NewClientSet creates a new API client
func (c *Client) NewClientSet() (*kubernetes.Clientset, error) {
	client, err := kubernetes.NewForConfig(c.rawConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client")
	}
	return client, nil
}

// NewStdClientSet creates a new API client in one go with an embedded STS token, this is most commonly used option
func (c *ClusterProvider) NewStdClientSet(spec *api.ClusterConfig) (*kubernetes.Clientset, error) {
	_, clientSet, err := c.newClientSetWithEmbeddedToken(spec)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func (c *ClusterProvider) newClientSetWithEmbeddedToken(spec *api.ClusterConfig) (*Client, *kubernetes.Clientset, error) {
	client, err := c.NewClient(spec)
	if err != nil {
		return nil, nil, errors.Wrap(err, "creating Kubernetes client config with embedded token")
	}

	clientSet, err := client.NewClientSet()
	if err != nil {
		return nil, nil, errors.Wrap(err, "creating Kubernetes client")
	}

	return client, clientSet, nil
}

// NewRawClient creates a new raw REST client in one go with an embedded STS token
func (c *ClusterProvider) NewRawClient(spec *api.ClusterConfig) (*kubewrapper.RawClient, error) {
	client, clientSet, err := c.newClientSetWithEmbeddedToken(spec)
	if err != nil {
		return nil, err
	}

	return kubewrapper.NewRawClient(clientSet, client.rawConfig)
}

// ServerVersion will use discovery API to fetch version of Kubernetes control plane
func (c *ClusterProvider) ServerVersion(rawClient *kubewrapper.RawClient) (string, error) {
	return rawClient.ServerVersion()
}

// UpdateAuthConfigMap creates or adds a nodegroup IAM role in the auth ConfigMap for the given nodegroup.
func (c *ClusterProvider) UpdateAuthConfigMap(nodeGroups []*api.NodeGroup, clientSet kubernetes.Interface) error {
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
func (c *ClusterProvider) WaitForNodes(clientSet kubernetes.Interface, ng KubeNodeGroup) error {
	minSize := ng.Size()
	if minSize == 0 {
		return nil
	}
	timeout := time.After(c.Provider.WaitTimeout())
	readyNodes := sets.NewString()
	watcher, err := clientSet.CoreV1().Nodes().Watch(context.TODO(), ng.ListOptions())
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}
	defer watcher.Stop()

	counter, err := getNodes(clientSet, ng)
	if err != nil {
		return errors.Wrap(err, "listing nodes")
	}

	logger.Info("waiting for at least %d node(s) to become ready in %q", minSize, ng.NameString())
	for {
		select {
		case event, ok := <-watcher.ResultChan():
			logger.Debug("event = %#v", event)
			if !ok {
				logger.Debug("the watcher channel was closed... stop processing events from it")
				return fmt.Errorf("the watcher channel for the nodes was closed by Kubernetes due to an unknown error")
			}
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
		case <-timeout:
			return fmt.Errorf("timed out (after %s) waiting for at least %d nodes to join the cluster and become ready in %q", c.Provider.WaitTimeout(), minSize, ng.NameString())
		}

		if counter >= minSize {
			break
		}
	}

	if _, err = getNodes(clientSet, ng); err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}

	return nil
}
