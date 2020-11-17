package cluster

import (
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
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
	currentVersion := c.ctl.ControlPlaneVersion()
	versionUpdateRequired, err := requiresVersionUpgrade(c.cfg.Metadata, currentVersion)
	if err != nil {
		return err
	}

	printer := printers.NewJSONPrinter()
	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", c.cfg); err != nil {
		return err
	}

	if versionUpdateRequired {
		if err := updateVersion(dryRun, c.cfg, currentVersion, c.ctl); err != nil {
			return err
		}
	}
	return nil
}
