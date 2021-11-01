package mockprovider

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	restclient "k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type MockKubeProvider struct {
	Client kubewrapper.Interface
}

func NewMockKubeProvider(client kubewrapper.Interface) *MockKubeProvider {
	return &MockKubeProvider{Client: client}
}

func (m MockKubeProvider) NewRawClient(spec *api.ClusterConfig) (kubewrapper.RawClientInterface, error) {
	return kubewrapper.NewRawClient(m.Client, &restclient.Config{})
}

func (m MockKubeProvider) NewClient(spec *api.ClusterConfig) (kubewrapper.ClientInterface, error) {
	return &MockKubeClient{Client: m.Client}, nil
}

func (m MockKubeProvider) NewStdClientSet(spec *api.ClusterConfig) (kubewrapper.Interface, error) {
	return m.Client, nil
}

type MockKubeClient struct {
	Client kubewrapper.Interface
}

func (m MockKubeClient) NewClientSet() (kubewrapper.Interface, error) {
	return m.Client, nil
}

func (m MockKubeClient) Config() *clientcmdapi.Config {
	return &clientcmdapi.Config{}
}
