package utils

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

func writeKubeconfigCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var (
		outputPath           string
		authenticatorRoleARN string
		setContext, autoPath bool
	)

	cmd.SetDescription("write-kubeconfig", "Write kubeconfig file for a given cluster", "")

	cmd.SetRunFuncWithNameArg(func() error {
		return doWriteKubeconfigCmd(cmd, outputPath, authenticatorRoleARN, setContext, autoPath)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmd.FlagSetGroup.InFlagSet("Output kubeconfig", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubeconfig(fs, &outputPath, &authenticatorRoleARN, &setContext, &autoPath, "<name>")
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doWriteKubeconfigCmd(cmd *cmdutils.Cmd, outputPath, roleARN string, setContext, autoPath bool) error {
	cfg := cmd.ClusterConfig

	// TODO: move this into a loader when --config-file gets added to this command
	if cfg.Metadata.Name != "" && cmd.NameArg != "" {
		return cmdutils.ErrClusterFlagAndArg(cmd, cfg.Metadata.Name, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		cfg.Metadata.Name = cmd.NameArg
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if autoPath {
		if outputPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig %s", cmdutils.IncompatibleFlags)
		}
		outputPath = kubeconfig.AutoPath(cfg.Metadata.Name)
	}

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(cfg.Metadata)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	kubectlConfig := kubeconfig.NewForKubectl(cfg, ctl.GetUsername(), roleARN, ctl.Provider.Profile())
	filename, err := kubeconfig.Write(outputPath, *kubectlConfig, setContext)
	if err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}

	logger.Success("saved kubeconfig as %q", filename)

	return nil
}
