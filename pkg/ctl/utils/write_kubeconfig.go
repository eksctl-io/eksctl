package utils

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

func writeKubeconfigCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	var (
		outputPath           string
		authenticatorRoleARN string
		setContext, autoPath bool
	)

	rc.SetDescription("write-kubeconfig", "Write kubeconfig file for a given cluster", "")

	rc.SetRunFuncWithNameArg(func() error {
		return doWriteKubeconfigCmd(rc, outputPath, authenticatorRoleARN, setContext, autoPath)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
	})

	rc.FlagSetGroup.InFlagSet("Output kubeconfig", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubeconfig(fs, &outputPath, &authenticatorRoleARN, &setContext, &autoPath, "<name>")
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doWriteKubeconfigCmd(rc *cmdutils.ResourceCmd, outputPath, roleARN string, setContext, autoPath bool) error {
	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && rc.NameArg != "" {
		return cmdutils.ErrNameFlagAndArg(cfg.Metadata.Name, rc.NameArg)
	}

	if rc.NameArg != "" {
		cfg.Metadata.Name = rc.NameArg
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}

	if autoPath {
		if outputPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig %s", cmdutils.IncompatibleFlags)
		}
		outputPath = kubeconfig.AutoPath(cfg.Metadata.Name)
	}

	if err := ctl.GetCredentials(cfg); err != nil {
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
