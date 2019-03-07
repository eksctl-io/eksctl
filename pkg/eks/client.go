package eks

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/cli-runtime/pkg/genericclioptions/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/kubernetes-sigs/aws-iam-authenticator/pkg/token"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

// Client stores information about the client config
type Client struct {
	Config      *clientcmdapi.Config
	ContextName string

	rawConfig *restclient.Config
}

// RawClient stores information about the client config
type RawClient struct {
	mapper    meta.RESTMapper
	config    *restclient.Config
	ClientSet kubernetes.Interface
}

// NewClient creates a new client config, if withEmbeddedToken is true
// it will embed the STS token, otherwise it will use authenticator exec plugin
// and ensures that AWS_PROFILE environment variable gets set also
func (c *ClusterProvider) NewClient(spec *api.ClusterConfig, withEmbeddedToken bool) (*Client, error) {
	clientConfig, _, contextName := kubeconfig.New(spec, c.getUsername(), "")

	config := &Client{
		Config:      clientConfig,
		ContextName: contextName,
	}

	return config.new(spec, withEmbeddedToken, c.Provider.STS(), c.Provider.Profile())
}

func (c *Client) new(spec *api.ClusterConfig, withEmbeddedToken bool, stsClient stsiface.STSAPI, profile string) (*Client, error) {
	if withEmbeddedToken {
		if err := c.useEmbeddedToken(spec, stsClient); err != nil {
			return nil, err
		}
	} else {
		kubeconfig.AppendAuthenticator(c.Config, spec, utils.DetectAuthenticator(), profile)
	}

	rawConfig, err := clientcmd.NewDefaultClientConfig(*c.Config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client configuration from client config")
	}
	c.rawConfig = rawConfig

	return c, nil
}

func (c *Client) useEmbeddedToken(spec *api.ClusterConfig, stsclient stsiface.STSAPI) error {
	gen, err := token.NewGenerator(true)
	if err != nil {
		return errors.Wrap(err, "could not get token generator")
	}

	tok, err := gen.GetWithSTS(spec.Metadata.Name, stsclient.(*sts.STS))
	if err != nil {
		return errors.Wrap(err, "could not get token")
	}

	c.Config.AuthInfos[c.ContextName].Token = tok
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
	client, err := c.NewClient(spec, true)
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
func (c *ClusterProvider) NewRawClient(spec *api.ClusterConfig) (*RawClient, error) {
	client, clientSet, err := c.newClientSetWithEmbeddedToken(spec)
	if err != nil {
		return nil, err
	}
	rawClient := &RawClient{
		config:    client.rawConfig,
		ClientSet: clientSet,
	}
	return rawClient.new()
}

func (c *RawClient) new() (*RawClient, error) {
	apiGroupResources, err := restmapper.GetAPIGroupResources(c.ClientSet.Discovery())
	if err != nil {
		return nil, errors.Wrap(err, "getting list of API resources for raw REST client")
	}

	for i, r := range apiGroupResources {
		logger.Debug("apiGroupResources[%d] = %#v", i, *r)
	}

	c.mapper = restmapper.NewDiscoveryRESTMapper(apiGroupResources)

	if c.config.APIPath == "" {
		c.config.APIPath = "/api"
	}
	if c.config.NegotiatedSerializer == nil {
		c.config.NegotiatedSerializer = &serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	}
	if err := restclient.SetKubernetesDefaults(c.config); err != nil {
		return nil, errors.Wrap(err, "applying defaults for REST client")
	}
	return c, nil
}

// NewFor construct a resource type-specific client for a give gvk
// (it's based on k8s.io/kubernetes/pkg/kubectl/cmd/util/factory_client_access.go)
func (c *RawClient) NewFor(gvk schema.GroupVersionKind) (*resource.Helper, error) {
	mapping, err := c.mapper.RESTMapping(gvk.GroupKind(), gvk.GroupVersion().Version, "")
	if err != nil {
		return nil, errors.Wrapf(err, "constructing REST client mapping for %s", gvk.String())
	}

	switch gvk.Group {
	case corev1.GroupName:
		c.config.APIPath = "/api"
	default:
		c.config.APIPath = "/apis"
	}
	gv := gvk.GroupVersion()
	c.config.GroupVersion = &gv

	client, err := restclient.RESTClientFor(c.config)
	if err != nil {
		return nil, errors.Wrapf(err, "constructing REST client for %s", gvk.String())
	}

	return resource.NewHelper(client, mapping), nil
}
