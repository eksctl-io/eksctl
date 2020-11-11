package cluster

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

type EKSCTLCluster struct {
	cfg          *api.ClusterConfig
	ctl          *eks.ClusterProvider
	stackManager *manager.StackCollection
}

func newEKSCTLCluster(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager *manager.StackCollection) (*EKSCTLCluster, error) {
	return &EKSCTLCluster{
		cfg:          cfg,
		ctl:          ctl,
		stackManager: stackManager,
	}, nil
}

func (c *EKSCTLCluster) Upgrade(dryRun bool) error {
	currentVersion := c.ctl.ControlPlaneVersion()
	versionUpdateRequired, err := requiresVersionUpgrade(c.cfg.Metadata, currentVersion)

	if err := c.ctl.LoadClusterVPC(c.cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", c.cfg.Metadata.Name)
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

	if err := c.ctl.RefreshClusterStatus(c.cfg); err != nil {
		return err
	}

	supportsManagedNodes, err := c.ctl.SupportsManagedNodes(c.cfg)
	if err != nil {
		return err
	}

	stackUpdateRequired, err := c.stackManager.AppendNewClusterStackResource(dryRun, supportsManagedNodes)
	if err != nil {
		return err
	}

	if err := c.ctl.ValidateExistingNodeGroupsForCompatibility(c.cfg, c.stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	cmdutils.LogPlanModeWarning(dryRun && (stackUpdateRequired || versionUpdateRequired))
	return nil

}
