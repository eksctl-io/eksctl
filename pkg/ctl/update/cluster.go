package update

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func updateClusterCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("cluster", "Upgrade control plane to the next version",
		"Upgrade control plane to the next Kubernetes version if available. Will also perform any updates needed in the cluster stack if resources are missing.")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doUpdateClusterCmd(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		// cmdutils.AddVersionFlag(fs, cfg.Metadata, `"next" and "latest" can be used to automatically increment version by one, or force latest`)

		cmdutils.AddApproveFlag(fs, cmd)
		fs.BoolVar(&cmd.Plan, "dry-run", cmd.Plan, "")
		_ = fs.MarkDeprecated("dry-run", "see --approve")

		cmdutils.AddWaitFlag(fs, &cmd.Wait, "all update operations to complete")
		_ = fs.MarkDeprecated("wait", "--wait is no longer respected; the cluster update always waits to complete")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)

}

func doUpdateClusterCmd(cmd *cmdutils.Cmd) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	printer := printers.NewJSONPrinter()

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(meta)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanUpdate(cfg); !ok {
		return err
	}

	if cmd.ClusterConfigFile != "" {
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
		cfg.Metadata.Version = api.Version1_14
	case api.Version1_14:
		cfg.Metadata.Version = api.Version1_15
	case api.Version1_15:
		cfg.Metadata.Version = api.Version1_15
	case api.Version1_16:
		cfg.Metadata.Version = api.Version1_16	
	default:
		// version of control plane is not known to us, maybe we are just too old...
		return fmt.Errorf("control plane version %q is not known to this version of eksctl, try to upgrade eksctl first", currentVersion)
	}
	versionUpdateRequired := cfg.Metadata.Version != currentVersion

	if err := ctl.LoadClusterVPC(cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	if versionUpdateRequired {
		msgNodeGroupsAndAddons := "you will need to follow the upgrade procedure for all of nodegroups and add-ons"
		cmdutils.LogIntendedAction(cmd.Plan, "upgrade cluster %q control plane from current version %q to %q", cfg.Metadata.Name, currentVersion, cfg.Metadata.Version)
		if !cmd.Plan {
			if err := ctl.UpdateClusterVersionBlocking(cfg); err != nil {
				return err
			}
			logger.Success("cluster %q control plane has been upgraded to version %q", cfg.Metadata.Name, cfg.Metadata.Version)
			logger.Info(msgNodeGroupsAndAddons)
		}
	}

	if err := ctl.RefreshClusterStatus(cfg); err != nil {
		return err
	}

	supportsManagedNodes, err := ctl.SupportsManagedNodes(cfg)
	if err != nil {
		return err
	}

	stackUpdateRequired, err := stackManager.AppendNewClusterStackResource(cmd.Plan, supportsManagedNodes)
	if err != nil {
		return err
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && (stackUpdateRequired || versionUpdateRequired))

	return nil
}
