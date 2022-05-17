package upgrade

import (
	"context"
	"time"

	"github.com/weaveworks/eksctl/pkg/actions/cluster"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// updating from 1.15 to 1.16 has been observed to take longer than the default value of 25 minutes
// increased to 50 for flex fleet changes
const upgradeClusterTimeout = 65 * time.Minute

func upgradeCluster(cmd *cmdutils.Cmd) {
	upgradeClusterWithRunFunc(cmd, DoUpgradeCluster)
}

func upgradeClusterWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	// Reset version
	cfg.Metadata.Version = ""
	cmd.ClusterConfig = cfg

	cmd.SetDescription("cluster", "Upgrade control plane to the next version",
		"Upgrade control plane to the next Kubernetes version if available. Will also perform any updates needed in the cluster stack if resources are missing.")

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		// cmdutils.AddVersionFlag(fs, cfg.Metadata, `"next" and "latest" can be used to automatically increment version by one, or force latest`)

		cmdutils.AddApproveFlag(fs, cmd)

		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, upgradeClusterTimeout)
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
	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	if ok, err := ctl.CanUpdate(cfg); !ok {
		return err
	}

	if cmd.ClusterConfigFile != "" {
		logger.Warning("NOTE: cluster VPC (subnets, routing & NAT Gateway) configuration changes are not yet implemented")
	}

	c, err := cluster.New(ctx, cfg, ctl)
	if err != nil {
		return err
	}

	return c.Upgrade(ctx, cmd.Plan)
}
