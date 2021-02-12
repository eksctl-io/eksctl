package identityproviders

import (
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type Manager struct {
	metadata api.ClusterMeta
	eksAPI   eksiface.EKSAPI
}

func NewManager(metadata api.ClusterMeta, eksAPI eksiface.EKSAPI) Manager {
	return Manager{
		metadata: metadata,
		eksAPI:   eksAPI,
	}
}
