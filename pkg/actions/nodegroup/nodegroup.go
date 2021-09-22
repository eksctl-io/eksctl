package nodegroup

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/utils/waiters"

	"github.com/aws/aws-sdk-go/aws/request"
	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type Manager struct {
	stackManager manager.StackManager
	ctl          *eks.ClusterProvider
	cfg          *api.ClusterConfig
	clientSet    kubernetes.Interface
	wait         WaitFunc
	init         eks.NodeGroupInitialiser
	kubeProvider eks.KubeProvider
}

type WaitFunc func(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error

// New creates a new manager.
func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) *Manager {
	return &Manager{
		stackManager: ctl.NewStackManager(cfg),
		ctl:          ctl,
		cfg:          cfg,
		clientSet:    clientSet,
		wait:         waiters.Wait,
		init: &eks.NodeGroupService{
			Provider: ctl.Provider,
		},
		kubeProvider: ctl,
	}
}

func (m *Manager) hasStacks(name string) (bool, error) {
	stacks, err := m.stackManager.ListNodeGroupStacks()
	if err != nil {
		return false, err
	}
	for _, stack := range stacks {
		if stack.NodeGroupName == name {
			return true, nil
		}
	}
	return false, nil
}
