package update

import (
	"time"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/upgrade"
)

func updateClusterCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	// Reset version before loading
	cfg.Metadata.Version = ""
	cmd.ClusterConfig = cfg

	cmd.SetDescription("cluster", "DEPRECATED: use 'upgrade cluster' instead. Upgrade control plane to the next version. ",
		"DEPRECATED: use 'upgrade cluster' instead. Upgrade control plane to the next Kubernetes version if available. Will also perform any updates needed in the cluster stack if resources are missing.")

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, "")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		// cmdutils.AddVersionFlag(fs, cfg.Metadata, `"next" and "latest" can be used to automatically increment version by one, or force latest`)

		cmdutils.AddApproveFlag(fs, cmd)
		fs.BoolVar(&cmd.Plan, "dry-run", cmd.Plan, "")
		_ = fs.MarkDeprecated("dry-run", "see --approve")

		cmdutils.AddWaitFlag(fs, &cmd.Wait, "all update operations to complete")
		_ = fs.MarkDeprecated("wait", "--wait is no longer respected; the cluster update always waits to complete")
		// updating from 1.15 to 1.16 has been observed to take longer than the default value of 25 minutes
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, 35*time.Minute)
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		logger.Warning("This command is to be deprecated. Please use 'eksctl upgrade cluster' instead")

		if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
			return err
		}

		return upgrade.DoUpgradeCluster(cmd)
	}

}
