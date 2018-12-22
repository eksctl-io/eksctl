package scale

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

func scaleNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "nodegroup NAME",
		Short: "Scale a nodegroup",
		RunE: func(_ *cobra.Command, args []string) error {
			name := cmdutils.GetNameArg(args)
			if name != "" {
				ng.Name = name
			}
			if err := doScaleNodeGroup(p, cfg, ng); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")

		fs.IntVarP(&ng.DesiredCapacity, "nodes", "N", -1, "total number of nodes (scale to this number)")

		fs.StringVarP(&p.Region, "region", "r", "", "AWS region")
		fs.StringVarP(&p.Profile, "profile", "p", "", "AWS creditials profile to use (overrides the AWS_PROFILE environment variable)")

		fs.DurationVar(&p.WaitTimeout, "timeout", api.DefaultWaitTimeout, "max wait time in any polling operations")
	})

	group.InFlagSet("Nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup. Generated if unset, e.g. \"ng-a345f4\"")
	})

	group.AddTo(cmd)

	return cmd
}

func doScaleNodeGroup(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return fmt.Errorf("no cluster name supplied. Use the --name= flag")
	}

	if ng.DesiredCapacity < 0 {
		return fmt.Errorf("number of nodes must be 0 or greater. Use the --nodes/-N flag")
	}

	stackManager := ctl.NewStackManager(cfg)
	err := stackManager.ScaleNodeGroup(ng)
	if err != nil {
		return fmt.Errorf("failed to scale nodegroup for cluster %q, error %v", cfg.Metadata.Name, err)
	}

	return nil
}
