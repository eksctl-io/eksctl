package utils

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func updateClusterStackCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "update-cluster-stack",
		Short: "Update cluster stack based on latest configuration",
		Run: func(_ *cobra.Command, args []string) {
			if err := doUpdateClusterStacksCmd(p, cfg, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name (required)")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddVersionFlag(fs, cfg.Metadata)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)
	return cmd
}

func doUpdateClusterStacksCmd(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && nameArg != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.Metadata.Name, nameArg)
	}

	if nameArg != "" {
		cfg.Metadata.Name = nameArg
	}

	if cfg.Metadata.Name == "" {
		return fmt.Errorf("--name must be set")
	}

	stackManager := ctl.NewStackManager(cfg)

	{
		stack, err := stackManager.DescribeClusterStack()
		if err != nil {
			return err
		}
		logger.Info("cluster = %#v", stack)
	}

	// if err := ctl.GetClusterVPC(cfg); err != nil {
	// 	return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	// }

	if err := stackManager.UpdateClusterForCompability(); err != nil {
		return err
	}

	{
		stack, err := stackManager.DescribeClusterStack()
		if err != nil {
			return err
		}
		logger.Info("cluster = %#v", stack)
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	return nil
}
