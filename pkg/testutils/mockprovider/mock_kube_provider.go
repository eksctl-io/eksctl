package mockprovider

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type MockKubeProvider struct {
	Client *kubernetes.Clientset
}

func (m MockKubeProvider) NewRawClient(spec *api.ClusterConfig) (kubewrapper.RawClientInterface, error) {
	return kubewrapper.NewRawClient(m.Client, &restclient.Config{})
}

func (m MockKubeProvider) NewClient(spec *api.ClusterConfig) (eks.ClientInterface, error) {
	return &MockKubeClient{Client: m.Client}, nil
}

func (m MockKubeProvider) NewStdClientSet(spec *api.ClusterConfig) (kubernetes.Interface, error) {
	return m.Client, nil
}

type MockKubeClient struct {
	Client *kubernetes.Clientset
}

func (m MockKubeClient) NewClientSet() (*kubernetes.Clientset, error) {
	return m.Client, nil
}

func (m MockKubeClient) Config() *clientcmdapi.Config {
	return &clientcmdapi.Config{}
}
