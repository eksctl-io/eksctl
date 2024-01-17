package eks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/transport"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	"github.com/weaveworks/eksctl/pkg/credentials"
	"github.com/weaveworks/eksctl/pkg/eks/auth"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

// Client stores information about the client config
type Client struct {
	Config *clientcmdapi.Config

	rawConfig *restclient.Config
}

// NewClient creates a new client config.
func (c *KubernetesProvider) NewClient(clusterInfo kubeconfig.ClusterInfo) (*Client, error) {
	config := kubeconfig.NewForUser(clusterInfo, GetUsername(c.RoleARN))
	client := &Client{
		Config: config,
	}
	tokenSource := &auth.TokenSource{
		ClusterID:      clusterInfo.ID(),
		TokenGenerator: auth.NewGenerator(c.Signer, &credentials.RealClock{}),
		Leeway:         1 * time.Minute,
	}
	return client.new(tokenSource)
}

// GetUsername extracts the username part from the IAM role ARN
func GetUsername(roleArn string) string {
	usernameParts := strings.Split(roleArn, "/")
	if len(usernameParts) > 1 {
		return usernameParts[len(usernameParts)-1]
	}
	return "iam-root-account"
}

func (c *Client) new(tokenSource oauth2.TokenSource) (*Client, error) {
	rawConfig, err := clientcmd.NewDefaultClientConfig(*c.Config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client configuration from client config")
	}
	rawConfig.WrapTransport = transport.TokenSourceWrapTransport(transport.NewCachedTokenSource(tokenSource))

	c.rawConfig = rawConfig
	c.rawConfig.QPS = float32(25)
	c.rawConfig.Burst = int(c.rawConfig.QPS) * 2

	return c, nil
}

// NewClientSet creates a new API client
func (c *Client) NewClientSet() (*kubernetes.Clientset, error) {
	client, err := kubernetes.NewForConfig(c.rawConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client")
	}
	return client, nil
}

// NewStdClientSet creates a new API client.
func (c *KubernetesProvider) NewStdClientSet(clusterInfo kubeconfig.ClusterInfo) (kubernetes.Interface, error) {
	_, clientSet, err := c.newClientSet(clusterInfo)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func (c *KubernetesProvider) newClientSet(clusterInfo kubeconfig.ClusterInfo) (*Client, *kubernetes.Clientset, error) {
	client, err := c.NewClient(clusterInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("creating Kubernetes client config: %w", err)
	}

	clientSet, err := client.NewClientSet()
	if err != nil {
		return nil, nil, fmt.Errorf("creating Kubernetes client: %w", err)
	}

	return client, clientSet, nil
}

// NewRawClient creates a new raw REST client.
func (c *KubernetesProvider) NewRawClient(clusterInfo kubeconfig.ClusterInfo) (*kubewrapper.RawClient, error) {
	client, clientSet, err := c.newClientSet(clusterInfo)
	if err != nil {
		return nil, err
	}

	return kubewrapper.NewRawClient(clientSet, client.rawConfig)
}

// ServerVersion will use discovery API to fetch version of Kubernetes control plane
func (c *KubernetesProvider) ServerVersion(rawClient *kubewrapper.RawClient) (string, error) {
	return rawClient.ServerVersion()
}

// WaitForControlPlane waits till the control plane is ready
func (c *KubernetesProvider) WaitForControlPlane(meta *api.ClusterMeta, clientSet *kubewrapper.RawClient, waitTimeout time.Duration) error {
	successCount := 0
	operation := func() (bool, error) {
		_, err := c.ServerVersion(clientSet)
		if err == nil {
			if successCount >= 5 {
				return true, nil
			}
			successCount++
			return false, nil
		}
		logger.Debug("control plane not ready yet – %s", err.Error())
		return false, nil
	}

	w := waiter.Waiter{
		Operation: operation,
		NextDelay: func(_ int) time.Duration {
			return 20 * time.Second
		},
	}

	if err := w.WaitWithTimeout(waitTimeout); err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("timed out waiting for control plane %q after %s", meta.Name, waitTimeout)
		}
		return err
	}
	return nil
}

// UpdateAuthConfigMap creates or adds a nodegroup IAM role in the auth ConfigMap for the given nodegroup.
func UpdateAuthConfigMap(nodeGroups []*api.NodeGroup, clientSet kubernetes.Interface) error {
	for _, ng := range nodeGroups {
		// authorise nodes to join
		if err := authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
			return err
		}
	}
	return nil
}

// WaitForNodes waits till the nodes are ready
func WaitForNodes(ctx context.Context, clientSet kubernetes.Interface, ng KubeNodeGroup) error {
	minSize := ng.Size()
	if minSize == 0 {
		return nil
	}

	readyNodes := sets.NewString()
	watcher, err := clientSet.CoreV1().Nodes().Watch(context.TODO(), ng.ListOptions())
	if err != nil {
		return errors.Wrap(err, "creating node watcher")
	}
	defer watcher.Stop()

	counter, err := GetNodes(clientSet, ng)
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
			if event.Type == watch.Error {
				logger.Debug("received an error event type from watcher: %+v", event.Object)
				msg := "unexpected error event type from node watcher"
				if statusErr, ok := event.Object.(*metav1.Status); ok {
					return fmt.Errorf("%s: %s", msg, statusErr.String())
				}
				return fmt.Errorf("%s: %+v", msg, event.Object)
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
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for at least %d nodes to join the cluster and become ready in %q: %w", minSize, ng.NameString(), ctx.Err())
		}

		if counter >= minSize {
			break
		}
	}

	if _, err = GetNodes(clientSet, ng); err != nil {
		return errors.Wrap(err, "re-listing nodes")
	}

	return nil
}
