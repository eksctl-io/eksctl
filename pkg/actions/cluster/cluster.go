package cluster

import (
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type Cluster interface {
	Upgrade(dryRun bool) error
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) (Cluster, error) {
	stackManager := ctl.NewStackManager(cfg)
	stacks, err := stackManager.DescribeStacks()
	if err != nil {
		return nil, err
	}

	if manager.IsClusterStack(stacks) {
		logger.Debug("Cluster %q was created by eksctl", cfg.Metadata.Name)
		return NewOwnedCluster(cfg, ctl, stackManager)
	}
	logger.Debug("Cluster %q was not created by eksctl", cfg.Metadata.Name)

	return NewUnownedCluster(cfg, ctl)
}
