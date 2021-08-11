// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
package enable

import (
	"time"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/repo"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/gitops"
)

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
type options struct {
	cfg     api.Git
	timeout time.Duration
}

func enableRepo(cmd *cmdutils.Cmd) {
	enableRepoWithRunFunc(cmd, doEnableRepository)
}

func enableRepoWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"repo",
		"DEPRECATED: https://github.com/weaveworks/eksctl/issues/2963. Set up a repo for gitops, installing Flux in the cluster and initializing its manifests in the specified Git repository",
		"",
	)
	cmd.CobraCommand.Hidden = true
	opts := configureRepositoryCmd(cmd)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)

		if cmd.NameArg != "" {
			return cmdutils.ErrUnsupportedNameArg()
		}

		if err := cmdutils.NewGitConfigLoader(cmd, &opts.cfg).WithRepoValidation().Load(); err != nil {
			return err
		}

		cmd.ProviderConfig.WaitTimeout = opts.timeout
		return runFunc(cmd)
	}
}

// configureRepositoryCmd configures the provided command object so that it can
// process CLI options and ClusterConfig file, to prepare for the "enablement"
// of the configured repository & cluster.
func configureRepositoryCmd(cmd *cmdutils.Cmd) *options {
	opts := options{
		cfg: *api.NewGit(),
	}
	cmd.FlagSetGroup.InFlagSet("Enable repository", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForFlux(fs, &opts.cfg)
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Metadata.Name, "cluster", "", "name of the EKS cluster to enable gitops on")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &opts.timeout, gitops.DefaultPodReadyTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
	return &opts
}

func doEnableRepository(cmd *cmdutils.Cmd) error {
	k8sClientSet, k8sRestConfig, err := cmdutils.KubernetesClientAndConfigFrom(cmd)
	if err != nil {
		return err
	}

	installer, err := repo.New(k8sRestConfig, k8sClientSet, cmd.ClusterConfig, cmd.ProviderConfig.WaitTimeout)
	if err != nil {
		return err
	}

	if err := installer.Run(); err != nil {
		logger.Critical("unable to set up gitops repo: %s", err.Error())
		return err
	}

	return nil
}
