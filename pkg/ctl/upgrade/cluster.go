package upgrade

import (
	"fmt"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils"
)

func upgradeCluster(cmd *cmdutils.Cmd) {
	upgradeClusterWithRunFunc(cmd, DoUpgradeCluster)
}

func upgradeClusterWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("cluster", "Upgrade control plane to the next version",
		"Upgrade control plane to the next Kubernetes version if available. Will also perform any updates needed in the cluster stack if resources are missing.")

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		// cmdutils.AddVersionFlag(fs, cfg.Metadata, `"next" and "latest" can be used to automatically increment version by one, or force latest`)

		cmdutils.AddApproveFlag(fs, cmd)

		// updating from 1.15 to 1.16 has been observed to take longer than the default value of 25 minutes
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, 35*time.Minute)
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)

		if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
			return err
		}
		return runFunc(cmd)
	}
}

// DoUpgradeCluster made public so that it can be shared with update/cluster.go until this is deprecated
// TODO Once `eksctl update cluster` is officially deprecated this can be made package private again
func DoUpgradeCluster(cmd *cmdutils.Cmd) error {
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
		logger.Warning("NOTE: cluster VPC (subnets, routing & NAT Gateway) configuration changes are not yet implemented")
	}
	currentVersion := ctl.ControlPlaneVersion()
	versionUpdateRequired, err := requiresVersionUpgrade(cfg.Metadata, currentVersion)
	if err != nil {
		return err
	}

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

func requiresVersionUpgrade(clusterMeta *api.ClusterMeta, currentEKSVersion string) (bool, error) {
	nextVersion, err := getNextVersion(currentEKSVersion)
	if err != nil {
		return false, err
	}

	// If the version was not specified default to the next Kubernetes version and assume the user intended to upgrade if possible
	if clusterMeta.Version == "" {
		if api.IsSupportedVersion(nextVersion) {
			clusterMeta.Version = nextVersion
			return true, nil
		}

		// There is no new version, stay in the current one
		clusterMeta.Version = currentEKSVersion
		return false, nil
	}

	if c, err := utils.CompareVersions(clusterMeta.Version, currentEKSVersion); err != nil {
		return false, err
	} else if c < 0 {
		return false, fmt.Errorf("cannot upgrade to a lower version. Found given target version %q, current cluster version %q", clusterMeta.Version, currentEKSVersion)
	}

	if api.IsDeprecatedVersion(clusterMeta.Version) {
		return false, fmt.Errorf("control plane version %q has been deprecated", clusterMeta.Version)
	}

	if !api.IsSupportedVersion(clusterMeta.Version) {
		return false, fmt.Errorf("control plane version %q is not known to this version of eksctl, try to upgrade eksctl first", clusterMeta.Version)
	}

	if clusterMeta.Version == currentEKSVersion {
		return false, nil
	}

	if clusterMeta.Version == nextVersion {
		return true, nil
	}

	return false, fmt.Errorf(
		"upgrading more than one version at a time is not supported. Found upgrade from %q to %q. Please upgrade to %q first",
		currentEKSVersion,
		clusterMeta.Version,
		nextVersion)
}

func getNextVersion(currentVersion string) (string, error) {
	switch currentVersion {
	case "":
		return "", errors.New("unable to get control plane version")
	case api.Version1_12:
		return api.Version1_13, nil
	case api.Version1_13:
		return api.Version1_14, nil
	case api.Version1_14:
		return api.Version1_15, nil
	case api.Version1_15:
		return api.Version1_16, nil
	case api.Version1_16:
		return api.Version1_17, nil
	default:
		// version of control plane is not known to us, maybe we are just too old...
		return "", fmt.Errorf("control plane version %q is not known to this version of eksctl, try to upgrade eksctl first", currentVersion)
	}
}
