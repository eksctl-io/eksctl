package enable

import (
	"context"
	"time"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/gitops"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
)

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
		"Set up a repo for gitops, installing Flux in the cluster and initializing its manifests in the specified Git repository",
		"",
	)
	opts := configureRepositoryCmd(cmd)
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)

		if cmd.NameArg != "" {
			return cmdutils.ErrUnsupportedNameArg()
		}

		if err := cmdutils.NewGitOpsConfigLoader(cmd, &opts.cfg).WithRepoValidation().Load(); err != nil {
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
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &opts.timeout, gitops.DefaultPodReadyTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	return &opts
}

// doEnableRepository enables GitOps on the configured repository.
func doEnableRepository(cmd *cmdutils.Cmd) error {
	k8sClientSet, k8sRestConfig, err := cmdutils.KubernetesClientAndConfigFrom(cmd)
	if err != nil {
		return err
	}

	installer, err := flux.NewInstaller(k8sRestConfig, k8sClientSet, cmd.ClusterConfig, cmd.ProviderConfig.WaitTimeout)
	if err != nil {
		return err
	}

	fluxIsInstalled, err := installer.IsFluxInstalled()
	if err != nil {
		// Continue with installation
		logger.Warning(err.Error())
	} else if fluxIsInstalled {
		logger.Warning("found existing flux deployment in namespace %q. Skipping installation",
			cmd.ClusterConfig.Git.Operator.Namespace)
		return nil
	}

	userInstructions, err := installer.Run(context.Background())
	if err != nil {
		logger.Critical("unable to set up gitops repo: %s", err.Error())
		return err
	}
	logger.Info(userInstructions)
	return err
}
