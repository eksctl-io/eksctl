package nodegroup

import (
	"k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type Manager struct {
	stackManager          manager.StackManager
	ctl                   *eks.ClusterProvider
	cfg                   *api.ClusterConfig
	clientSet             kubernetes.Interface
	instanceSelector      eks.InstanceSelector
	launchTemplateFetcher *builder.LaunchTemplateFetcher
	accessEntry           *accessentry.Service
}

// New creates a new manager.
func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface, instanceSelector eks.InstanceSelector) *Manager {
	return &Manager{
		stackManager:          ctl.NewStackManager(cfg),
		ctl:                   ctl,
		cfg:                   cfg,
		clientSet:             clientSet,
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
