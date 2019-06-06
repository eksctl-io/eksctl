package utils

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

var (
	writeKubeconfigOutputPath string
	writeKubeconfigSetContext bool
	writeKubeconfigAutoPath   bool
)

func writeKubeconfigCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:   "write-kubeconfig",
		Short: "Write kubeconfig file for a given cluster",
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doWriteKubeconfigCmd(cp); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
	})

	group.InFlagSet("Output kubeconfig", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubeconfig(fs, &writeKubeconfigOutputPath, &writeKubeconfigSetContext, &writeKubeconfigAutoPath, "<name>")
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, false)

	group.AddTo(cp.Command)
	return cp.Command
}

func doWriteKubeconfigCmd(cp *cmdutils.CommonParams) error {
	cfg := cp.ClusterConfig

	ctl := eks.New(cp.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && cp.NameArg != "" {
		return cmdutils.ErrNameFlagAndArg(cfg.Metadata.Name, cp.NameArg)
	}

	if cp.NameArg != "" {
		cfg.Metadata.Name = cp.NameArg
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}

	if writeKubeconfigAutoPath {
		if writeKubeconfigOutputPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig %s", cmdutils.IncompatibleFlags)
		}
		writeKubeconfigOutputPath = kubeconfig.AutoPath(cfg.Metadata.Name)
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return err
	}

	client, err := ctl.NewClient(cfg, false)
	if err != nil {
		return err
	}

	filename, err := kubeconfig.Write(writeKubeconfigOutputPath, *client.Config, writeKubeconfigSetContext)
	if err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}

	logger.Success("saved kubeconfig as %q", filename)

	return nil
}
