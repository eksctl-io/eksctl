package update

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func updateClusterCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:   "cluster",
		Short: "Update cluster",
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doUpdateClusterCmd(cp); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cp.ClusterConfigFile)

		// cmdutils.AddVersionFlag(fs, cfg.Metadata, `"next" and "latest" can be used to automatically increment version by one, or force latest`)
		cmdutils.AddApproveFlag(fs, cp)
		fs.BoolVar(&cp.Plan, "dry-run", cp.Plan, "")
		fs.MarkDeprecated("dry-run", "see --aprove")

		cmdutils.AddWaitFlag(fs, &cp.Wait, "all update operations to complete")
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, false)

	group.AddTo(cp.Command)
	return cp.Command
}

func doUpdateClusterCmd(cp *cmdutils.CommonParams) error {
	if err := cmdutils.NewMetadataLoader(cp).Load(); err != nil {
		return err
	}

	cfg := cp.ClusterConfig
	meta := cp.ClusterConfig.Metadata

	printer := printers.NewJSONPrinter()
	ctl := eks.New(cp.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(cp.ProviderConfig)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	if cp.ClusterConfigFile != "" {
		logger.Warning("NOTE: config file is only used for finding cluster name and region, deep cluster configuration changes are not yet implemented")
	}

	currentVersion := ctl.ControlPlaneVersion()
	// determine next version based on what's currently deployed
	switch currentVersion {
	case "":
		return fmt.Errorf("unable to get control plane version")
	case api.Version1_10:
		cfg.Metadata.Version = api.Version1_11
	case api.Version1_11:
		cfg.Metadata.Version = api.Version1_12
	case api.Version1_12:
		cfg.Metadata.Version = api.Version1_12
	default:
		// version of control is not known to us, maybe we are just too old...
		return fmt.Errorf("control plane version version %q is known to this version of eksctl, try to upgrade eksctl first", currentVersion)
	}
	versionUpdateRequired := cfg.Metadata.Version != currentVersion

	if err := ctl.GetClusterVPC(cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	stackUpdateRequired, err := stackManager.AppendNewClusterStackResource(cp.Plan)
	if err != nil {
		return err
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	if versionUpdateRequired {
		msgNodeGroupsAndAddons := "you will need to follow the upgrade procedure for all of nodegroups and add-ons"
		cmdutils.LogIntendedAction(cp.Plan, "upgrade cluster %q control plane from current version %q to %q", cfg.Metadata.Name, currentVersion, cfg.Metadata.Version)
		if !cp.Plan {
			if cp.Wait {
				if err := ctl.UpdateClusterVersionBlocking(cfg); err != nil {
					return err
				}
				logger.Success("cluster %q control plan e has been upgraded to version %q", cfg.Metadata.Name, cfg.Metadata.Version)
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

	cmdutils.LogPlanModeWarning(cp.Plan && (stackUpdateRequired || versionUpdateRequired))

	return nil
}
