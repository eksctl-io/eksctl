package enable

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/flux"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/version"
)

func enableFlux2(cmd *cmdutils.Cmd) {
	configureAndRun(cmd, flux2Install)
}

func configureAndRun(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"flux",
		"Set up GitOps Toolkit - deploys FluxV2 and creates Git repo to store manifests",
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
	logger.Info("eksctl version %s", version.GetVersion())
	logger.Info("will install Flux v2 components on cluster %s", cmd.ClusterConfig.Metadata.Name)

	if kubeconfAndContextNotSet(cmd.ClusterConfig.GitOps.Flux.Flags) {
		ctl, err := cmd.NewProviderForExistingCluster()
		if err != nil {
			return err
		}
		if ok, err := ctl.CanOperate(cmd.ClusterConfig); !ok {
			return err
		}

		kubeCfgPath, err := os.CreateTemp("", cmd.ClusterConfig.Metadata.Name)
		if err != nil {
			return err
		}
		defer func() {
			if err := os.Remove(kubeCfgPath.Name()); err != nil {
				logger.Critical("failed to remove temporary kubeconfig", kubeCfgPath)
			}
		}()
		logger.Debug("writing temporary kubeconfig to %s", kubeCfgPath.Name())
		kubectlConfig := kubeconfig.NewForKubectl(cmd.ClusterConfig, ctl.GetUsername(), "", ctl.Provider.Profile())
		if _, err := kubeconfig.Write(kubeCfgPath.Name(), *kubectlConfig, true); err != nil {
			return err
		}
		cmd.ClusterConfig.GitOps.Flux.Flags["kubeconfig"] = kubeCfgPath.Name()
	}

	k8sClientSet, _, err := cmdutils.KubernetesClientAndConfigFrom(cmd)
	if err != nil {
		return err
	}

	installer, err := flux.New(k8sClientSet, cmd.ClusterConfig.GitOps)
	if err != nil {
		return err
	}

	return installer.Run()
}

func kubeconfAndContextNotSet(flags map[string]string) bool {
	_, cfg := flags["kubeconfig"]
	_, ctx := flags["context"]
	return !cfg && !ctx
}
