package nodegroup

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	"k8s.io/client-go/kubernetes"
)

type Manager struct {
	manager   *manager.StackCollection
	ctl       *eks.ClusterProvider
	cfg       *api.ClusterConfig
	clientSet *kubernetes.Clientset
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet *kubernetes.Clientset) *Manager {
	return &Manager{
		manager:   ctl.NewStackManager(cfg),
		ctl:       ctl,
		cfg:       cfg,
		clientSet: clientSet,
	}
}
