package fargate

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	"k8s.io/client-go/kubernetes"
)

type Manager struct {
	ctl             *eks.ClusterProvider
	cfg             *api.ClusterConfig
	stackManager    manager.StackManager
	newStdClientSet func() (kubernetes.Interface, error)
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager) *Manager {
	return &Manager{
		ctl:             ctl,
		cfg:             cfg,
		stackManager:    stackManager,
		newStdClientSet: func() (kubernetes.Interface, error) { return ctl.NewStdClientSet(cfg) },
	}
}
