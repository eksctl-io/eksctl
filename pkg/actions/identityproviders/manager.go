package identityproviders

import (
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type IdentityProviderManager struct {
	metadata api.ClusterMeta
	eksAPI   eksiface.EKSAPI
}

func NewIdentityProviderManager(metadata api.ClusterMeta, eksAPI eksiface.EKSAPI) IdentityProviderManager {
	return IdentityProviderManager{
		metadata: metadata,
		eksAPI:   eksAPI,
	}
}
