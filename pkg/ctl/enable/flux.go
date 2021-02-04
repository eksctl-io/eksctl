package enable

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/flux"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/version"
)

func enableFlux2(cmd *cmdutils.Cmd) {
	configureAndRun(cmd, flux2Install)
}

func configureAndRun(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"flux",
		"EXPERIMENTAL. Set up GitOps Toolkit - deploys FluxV2 and creates Git repo to store manifests",
		"",
	)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)

		if cmd.NameArg != "" {
			return cmdutils.ErrUnsupportedNameArg()
		}

		if err := cmdutils.NewGitOpsConfigLoader(cmd).Load(); err != nil {
			return err
		}

		return runFunc(cmd)
	}
}

func flux2Install(cmd *cmdutils.Cmd) error {
	logger.Warning("running experimental command")
	logger.Info("eksctl version %s", version.GetVersion())
	logger.Info("will install Flux v2 components on cluster %s", cmd.ClusterConfig.Metadata.Name)

	k8sClientSet, _, err := cmdutils.KubernetesClientAndConfigFrom(cmd)
	if err != nil {
		return err
	}

	installer, err := flux.New(k8sClientSet, cmd.ClusterConfig.GitOps)
	if err != nil {
		return err
	}

	if err := installer.Run(); err != nil {
		return err
	}

	return nil
}
