package iamidentitymapping

import (
	kubeclient "k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type Manager struct {
	clusterConfig   *api.ClusterConfig
	clientSet       kubeclient.Interface
	clusterProvider *eks.ClusterProvider
	region          string
}

func New(clusterConfig *api.ClusterConfig, clientSet kubeclient.Interface, clusterProvider *eks.ClusterProvider, region string) (*Manager, error) {
	return &Manager{
		clusterConfig:   clusterConfig,
		clientSet:       clientSet,
		clusterProvider: clusterProvider,
		region:          region,
	}, nil
}

