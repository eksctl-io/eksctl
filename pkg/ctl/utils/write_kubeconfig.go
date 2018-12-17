package utils

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

var (
	writeKubeconfigOutputPath string
	writeKubeconfigSetContext bool
	writeKubeconfigAutoPath   bool
)

func writeKubeconfigCmd() *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "write-kubeconfig",
		Short: "Write kubeconfig file for a given cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doWriteKubeconfigCmd(p, cfg, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := &cmdutils.NamedFlagSetGroup{}

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name (required)")
		cmdutils.AddRegionFlag(fs, p)
	})

	group.InFlagSet("Output kubeconfig", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForKubeconfig(fs, &writeKubeconfigOutputPath, &writeKubeconfigSetContext, &writeKubeconfigAutoPath, "<name>")
	})

	cmdutils.AddCommonFlagsForAWS(group, p)

	group.AddTo(cmd)
	return cmd
}

func doWriteKubeconfigCmd(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && nameArg != "" {
		return cmdutils.ErrNameFlagAndArg(cfg.Metadata.Name, nameArg)
	}

	if nameArg != "" {
		cfg.Metadata.Name = nameArg
	}

	if cfg.Metadata.Name == "" {
		return fmt.Errorf("--name must be set")
	}

	if writeKubeconfigAutoPath {
		if writeKubeconfigOutputPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig %s", cmdutils.IncompatibleFlags)
		}
		writeKubeconfigOutputPath = kubeconfig.AutoPath(cfg.Metadata.Name)
	}

	cluster, err := ctl.DescribeControlPlane(cfg.Metadata)
	if err != nil {
		return err
	}

	logger.Debug("cluster = %#v", cluster)

	if err = ctl.GetCredentials(*cluster, cfg); err != nil {
		return err
	}

	clientConfigBase, err := ctl.NewClientConfig(cfg)
	if err != nil {
		return err
	}

	config := clientConfigBase.WithExecAuthenticator()
	filename, err := kubeconfig.Write(writeKubeconfigOutputPath, config.Client, writeKubeconfigSetContext)
	if err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}

	logger.Success("saved kubeconfig as %q", filename)

	return nil
}
