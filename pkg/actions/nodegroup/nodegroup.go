package nodegroup

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/utils/waiters"

	"github.com/aws/aws-sdk-go/aws/request"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	"k8s.io/client-go/kubernetes"
)

type Manager struct {
	stackManager manager.StackManager
	ctl          *eks.ClusterProvider
	cfg          *api.ClusterConfig
	clientSet    *kubernetes.Clientset
	wait         WaitFunc
}

type WaitFunc func(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet *kubernetes.Clientset) *Manager {
	return &Manager{
		stackManager: ctl.NewStackManager(cfg),
		ctl:          ctl,
		cfg:          cfg,
		clientSet:    clientSet,
		wait:         waiters.Wait,
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
