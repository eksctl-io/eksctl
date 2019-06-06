package scale

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func scaleNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:     "nodegroup",
		Short:   "Scale a nodegroup",
		Aliases: []string{"ng"},
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doScaleNodeGroup(cp, ng); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to scale")

		desiredCapacity := fs.IntP("nodes", "N", -1, "total number of nodes (scale to this number)")
		cmdutils.AddPreRun(cp.Command, func(cmd *cobra.Command, args []string) {
			if f := cmd.Flag("nodes"); f.Changed {
				ng.DesiredCapacity = desiredCapacity
			}
		})

		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, true)

	group.AddTo(cp.Command)
	return cp.Command
}

func doScaleNodeGroup(cp *cmdutils.CommonParams, ng *api.NodeGroup) error {
	cfg := cp.ClusterConfig

	ctl := eks.New(cp.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--cluster")
	}

	if ng.Name != "" && cp.NameArg != "" {
		return cmdutils.ErrNameFlagAndArg(ng.Name, cp.NameArg)
	}

	if cp.NameArg != "" {
		ng.Name = cp.NameArg
	}

	if ng.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}

	if ng.DesiredCapacity == nil || *ng.DesiredCapacity < 0 {
		return fmt.Errorf("number of nodes must be 0 or greater. Use the --nodes/-N flag")
	}

	stackManager := ctl.NewStackManager(cfg)
	err := stackManager.ScaleNodeGroup(ng)
	if err != nil {
		return fmt.Errorf("failed to scale nodegroup for cluster %q, error %v", cfg.Metadata.Name, err)
	}

	return nil
}
