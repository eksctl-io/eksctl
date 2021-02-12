package identityproviders

import (
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type AssociateProvidersTask struct {
	metadata  api.ClusterMeta
	providers []api.IdentityProvider
	eks       eksiface.EKSAPI
}

func NewAssociateProvidersTask(metadata api.ClusterMeta, providers []api.IdentityProvider, eks eksiface.EKSAPI) tasks.Task {
	return tasks.SynchronousTask{
		SynchronousTaskIface: &AssociateProvidersTask{
			metadata:  metadata,
			providers: providers,
			eks:       eks,
		},
	}
}

func (t *AssociateProvidersTask) Describe() string {
	return "associate identity providers with cluster"
}

func (t *AssociateProvidersTask) Do() error {
	m := NewManager(t.metadata, t.eks)
	return m.Associate(AssociateIdentityProvidersOptions{Providers: t.providers})
}
