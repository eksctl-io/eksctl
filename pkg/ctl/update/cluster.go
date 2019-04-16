package update

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func updateClusterCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	// cfg.Metadata.Version = "next"

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Update cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doUpdateClusterCmd(p, cfg, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)

		// cmdutils.AddVersionFlag(fs, cfg.Metadata, `"next" and "latest" can be used to automatically increment version by one, or force latest`)
		cmdutils.AddApproveFlag(&plan, cmd, fs)
		fs.BoolVar(&plan, "dry-run", plan, "")
		fs.MarkDeprecated("dry-run", "see --aprove")

		cmdutils.AddWaitFlag(&wait, fs, "all update operations to complete")
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)
	return cmd
}

func doUpdateClusterCmd(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string, cmd *cobra.Command) error {
	if err := cmdutils.NewMetadataLoader(p, cfg, clusterConfigFile, nameArg, cmd).Load(); err != nil {
		return err
	}

	ctl := eks.New(p, cfg)
	meta := cfg.Metadata

	printer := printers.NewJSONPrinter()

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	if clusterConfigFile != "" {
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

	stackUpdateRequired, err := stackManager.AppendNewClusterStackResource(plan)
	if err != nil {
		return err
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	if versionUpdateRequired {
		msgNodeGroupsAndAddons := "you will need to follow the upgrade procedure for all of nodegroups and add-ons"
		cmdutils.LogIntendedAction(plan, "upgrade cluster %q control plane from current version %q to %q", cfg.Metadata.Name, currentVersion, cfg.Metadata.Version)
		if !plan {
			if wait {
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

	cmdutils.LogPlanModeWarning(plan && (stackUpdateRequired || versionUpdateRequired))

	return nil
}
