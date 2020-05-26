package enable

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/gitops"
)

func enableProfile(cmd *cmdutils.Cmd) {
	enableProfileWithRunFunc(cmd, doEnableProfile)
}

func enableProfileWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"profile",
		"Commits the components from the selected Quick Start profile to the destination repository.",
		"",
	)
	opts := configureProfileCmd(cmd)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)

		if err := cmdutils.NewGitOpsConfigLoader(cmd, opts).
			WithRepoValidation().
			WithProfileValidation().
			Load(); err != nil {
			return err
		}

		return runFunc(cmd)
	}
}

// configureProfileCmd configures the provided command object so that it can
// process CLI options and ClusterConfig file, to prepare for the installation
// of the configured profile on the configured repository & cluster.
func configureProfileCmd(cmd *cmdutils.Cmd) *api.Git {
	opts := api.NewGit()

	cmd.FlagSetGroup.InFlagSet("Enable profile", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForProfile(fs, opts.BootstrapProfile)
		cmdutils.AddCommonFlagsForGit(fs, opts.Repo)
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Metadata.Name, "cluster", "", "name of the EKS cluster to enable this Quick Start profile on")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, 20*time.Second)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	return opts
}

// doEnableProfile enables the configured profile on the configured repository.
func doEnableProfile(cmd *cmdutils.Cmd) error {
	return gitops.InstallProfile(cmd.ClusterConfig)
}
