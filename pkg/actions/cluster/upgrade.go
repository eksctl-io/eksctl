package cluster

import (
	"context"

	"github.com/weaveworks/eksctl/pkg/printers"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func upgrade(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, dryRun bool) (bool, error) {
	cvm, err := eks.NewClusterVersionsManager(ctl.AWSProvider.EKS())
	if err != nil {
		return false, err
	}

	upgradeVersion, err := cvm.ResolveUpgradeVersion(
		/* desiredVersion */ cfg.Metadata.Version,
		/* currentVersion */ ctl.ControlPlaneVersion())
	if err != nil {
		return false, err
	}

	printer := printers.NewJSONPrinter()
	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return false, err
	}

	if upgradeVersion != "" {
		msgNodeGroupsAndAddons := "you will need to follow the upgrade procedure for all of nodegroups and add-ons"
		cmdutils.LogIntendedAction(dryRun, "upgrade cluster %q control plane from current version %q to %q", cfg.Metadata.Name, ctl.ControlPlaneVersion(), cfg.Metadata.Version)
		if !dryRun {
			cfg.Metadata.Version = upgradeVersion
			if err := ctl.UpdateClusterVersionBlocking(ctx, cfg); err != nil {
				return false, err
			}
			logger.Success("cluster %q control plane has been upgraded to version %q", cfg.Metadata.Name, cfg.Metadata.Version)
			logger.Info(msgNodeGroupsAndAddons)
		}
	} else {
		logger.Info("no cluster version update required")
	}
	return upgradeVersion != "", nil
}
