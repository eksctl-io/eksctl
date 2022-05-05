package identityproviders

import (
	"context"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type AssociateProvidersTask struct {
	metadata  api.ClusterMeta
	providers []api.IdentityProvider
	eks       awsapi.EKS
	ctx       context.Context
}

func NewAssociateProvidersTask(ctx context.Context, metadata api.ClusterMeta, providers []api.IdentityProvider, eks awsapi.EKS) tasks.Task {
	return tasks.SynchronousTask{
		SynchronousTaskIface: &AssociateProvidersTask{
			metadata:  metadata,
			providers: providers,
			eks:       eks,
			ctx:       ctx,
		},
	}
}

func (t *AssociateProvidersTask) Describe() string {
	return "associate identity providers with cluster"
}

func (t *AssociateProvidersTask) Do() error {
	m := NewManager(t.metadata, t.eks)
	return m.Associate(t.ctx, AssociateIdentityProvidersOptions{Providers: t.providers})
}
