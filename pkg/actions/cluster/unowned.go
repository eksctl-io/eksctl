package cluster

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

type UnownedCluster struct {
	cfg *api.ClusterConfig
	ctl *eks.ClusterProvider
}

func NewUnownedCluster(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) (*UnownedCluster, error) {
	return &UnownedCluster{
		cfg: cfg,
		ctl: ctl,
	}, nil
}

func (c *UnownedCluster) Upgrade(dryRun bool) error {
	versionUpdateRequired, err := upgrade(c.cfg, c.ctl, dryRun)
	if err != nil {
		return err
	}

	// if no version update is required, don't log asking them to rerun with --approve
	cmdutils.LogPlanModeWarning(dryRun && versionUpdateRequired)
	return nil
}
