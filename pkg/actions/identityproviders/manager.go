package identityproviders

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
)

type Manager struct {
	metadata api.ClusterMeta
	eksAPI   awsapi.EKS
}

func NewManager(metadata api.ClusterMeta, eksAPI awsapi.EKS) Manager {
	return Manager{
		metadata: metadata,
		eksAPI:   eksAPI,
	}
}
