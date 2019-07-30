package update

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func updateClusterCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	rc.SetDescription("cluster", "Update cluster", "")

	rc.SetRunFuncWithNameArg(func() error {
		return doUpdateClusterCmd(rc)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)

		// cmdutils.AddVersionFlag(fs, cfg.Metadata, `"next" and "latest" can be used to automatically increment version by one, or force latest`)

		cmdutils.AddApproveFlag(fs, rc)
		fs.BoolVar(&rc.Plan, "dry-run", rc.Plan, "")
		fs.MarkDeprecated("dry-run", "see --aprove")

		rc.Wait = true
		cmdutils.AddWaitFlag(fs, &rc.Wait, "all update operations to complete")
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)

}

func doUpdateClusterCmd(rc *cmdutils.ResourceCmd) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig
	meta := rc.ClusterConfig.Metadata

	if err := api.SetClusterConfigDefaults(cfg); err != nil {
		return err
	}

	printer := printers.NewJSONPrinter()
	ctl := eks.New(rc.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(rc.ProviderConfig)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.RefreshClusterConfig(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	if rc.ClusterConfigFile != "" {
		logger.Warning("NOTE: config file is used for finding cluster name and region")
		logger.Warning("NOTE: cluster VPC (subnets, routing & NAT Gateway) configuration changes are not yet implemented")
	}

	currentVersion := ctl.ControlPlaneVersion()
	// determine next version based on what's currently deployed
	switch currentVersion {
	case "":
		return errors.New("unable to get control plane version")
	case api.Version1_11:
		cfg.Metadata.Version = api.Version1_12
	case api.Version1_12:
		cfg.Metadata.Version = api.Version1_13
	case api.Version1_13:
		cfg.Metadata.Version = api.Version1_13
	default:
		// version of control is not known to us, maybe we are just too old...
		return fmt.Errorf("control plane version %q is not known to this version of eksctl, try to upgrade eksctl first", currentVersion)
	}
	versionUpdateRequired := cfg.Metadata.Version != currentVersion

	if err := ctl.GetClusterVPC(cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	stackUpdateRequired, err := stackManager.AppendNewClusterStackResource(rc.Plan)
	if err != nil {
		return err
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	if versionUpdateRequired {
		msgNodeGroupsAndAddons := "you will need to follow the upgrade procedure for all of nodegroups and add-ons"
		cmdutils.LogIntendedAction(rc.Plan, "upgrade cluster %q control plane from current version %q to %q", cfg.Metadata.Name, currentVersion, cfg.Metadata.Version)
		if !rc.Plan {
			if rc.Wait {
				if err := ctl.UpdateClusterVersionBlocking(cfg); err != nil {
					return err
				}
				logger.Success("cluster %q control plane has been upgraded to version %q", cfg.Metadata.Name, cfg.Metadata.Version)
				logger.Info(msgNodeGroupsAndAddons)
			} else {
				if _, err := ctl.UpdateClusterVersion(cfg); err != nil {
					return err
				}
				logger.Success("a version update operation has been requested for cluster %q", cfg.Metadata.Name)
				logger.Info("once it has been updated, %s", cfg.Metadata.Name, msgNodeGroupsAndAddons)
			}
		}
	}

	cmdutils.LogPlanModeWarning(rc.Plan && (stackUpdateRequired || versionUpdateRequired))

	return nil
}
