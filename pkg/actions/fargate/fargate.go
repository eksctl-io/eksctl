package fargate

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type Manager struct {
	ctl *eks.ClusterProvider
	cfg *api.ClusterConfig
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) *Manager {
	return &Manager{
		ctl: ctl,
		cfg: cfg,
	}
}
