package utils

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

var (
	utilsKubeconfigInputPath  string
	utilsKubeconfigOutputPath string
	utilsSetContext           bool
	utilsAutoKubeconfigPath   bool
)

func writeKubeconfigCmd() *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "write-kubeconfig",
		Short: "Write kubeconfig file for a given cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doWriteKubeconfigCmd(p, cfg, ctl.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name (required)")

	fs.StringVarP(&p.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&p.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.BoolVar(&utilsAutoKubeconfigPath, "auto-kubeconfig", false, fmt.Sprintf("save kubeconfig file by cluster name – %q", kubeconfig.AutoPath("<name>")))
	fs.StringVar(&utilsKubeconfigOutputPath, "kubeconfig", kubeconfig.DefaultPath, "path to write kubeconfig")
	fs.BoolVar(&utilsSetContext, "set-kubeconfig-context", true, "if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten")

	return cmd
}

func doWriteKubeconfigCmd(p *api.ProviderConfig, cfg *api.ClusterConfig, name string) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && name != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.Metadata.Name, name)
	}

	if name != "" {
		cfg.Metadata.Name = name
	}

	if cfg.Metadata.Name == "" {
		return fmt.Errorf("--name must be set")
	}

	if utilsAutoKubeconfigPath {
		if utilsKubeconfigOutputPath != kubeconfig.DefaultPath {
			return fmt.Errorf("--kubeconfig and --auto-kubeconfig cannot be used at the same time")
		}
		utilsKubeconfigOutputPath = kubeconfig.AutoPath(cfg.Metadata.Name)
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
	filename, err := kubeconfig.Write(utilsKubeconfigOutputPath, config.Client, utilsSetContext)
	if err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}

	logger.Success("saved kubeconfig as %q", filename)

	return nil
}
