package eks

import (
	"strings"

	"github.com/pkg/errors"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

// Client stores information about the client config
type Client struct {
	Config *clientcmdapi.Config

	rawConfig *restclient.Config
}

// NewClient creates a new client config by embedding the STS token
func (c *ClusterProvider) NewClient(spec *api.ClusterConfig) (*Client, error) {
	clientConfig := kubeconfig.
		NewBuilder(spec.Metadata, spec.Status, c.GetUsername()).
		Build()

	config := &Client{
		Config: clientConfig,
	}

	return config.new(spec, c.Provider.STS())
}

// GetUsername extracts the username part from the IAM role ARN
func (c *ClusterProvider) GetUsername() string {
	usernameParts := strings.Split(c.Status.iamRoleARN, "/")
	if len(usernameParts) > 1 {
		return usernameParts[len(usernameParts)-1]
	}
	return "iam-root-account"
}

func (c *Client) new(spec *api.ClusterConfig, stsClient stsiface.STSAPI) (*Client, error) {
	if err := c.useEmbeddedToken(spec, stsClient); err != nil {
		return nil, err
	}

	rawConfig, err := clientcmd.NewDefaultClientConfig(*c.Config, &clientcmd.ConfigOverrides{}).ClientConfig()
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

	tok, err := gen.GetWithSTS(spec.Metadata.Name, stsclient.(*sts.STS))
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
