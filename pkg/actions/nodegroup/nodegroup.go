package nodegroup

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
)

type Manager struct {
	stackManager          manager.StackManager
	ctl                   *eks.ClusterProvider
	cfg                   *api.ClusterConfig
	clientSet             kubernetes.Interface
	wait                  WaitFunc
	instanceSelector      eks.InstanceSelector
	launchTemplateFetcher *builder.LaunchTemplateFetcher
	accessEntry           *accessentry.Service
}

type WaitFunc func(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error

// New creates a new manager.
func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface, instanceSelector eks.InstanceSelector) *Manager {
	return &Manager{
		stackManager:          ctl.NewStackManager(cfg),
		ctl:                   ctl,
		cfg:                   cfg,
		clientSet:             clientSet,
		wait:                  waiters.Wait,
		instanceSelector:      instanceSelector,
		launchTemplateFetcher: builder.NewLaunchTemplateFetcher(ctl.AWSProvider.EC2()),
		accessEntry: &accessentry.Service{
			ClusterStateGetter: ctl,
		},
	}
}

func findStack(stacks []manager.NodeGroupStack, name string) *manager.NodeGroupStack {
	for _, stack := range stacks {
		if stack.NodeGroupName == name {
			return &stack
		}
	}
	return nil
}
