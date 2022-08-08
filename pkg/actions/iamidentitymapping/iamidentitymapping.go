package iamidentitymapping

import (
	kubeclient "k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

type Manager struct {
	clusterConfig *api.ClusterConfig
	clientSet     kubeclient.Interface
	rawClient     *kubernetes.RawClient
	region        string
}

func New(clusterConfig *api.ClusterConfig, clientSet kubeclient.Interface, rawClient *kubernetes.RawClient, region string) (*Manager, error) {
	return &Manager{
		clusterConfig: clusterConfig,
		clientSet:     clientSet,
		rawClient:     rawClient,
		region:        region,
	}, nil
}
