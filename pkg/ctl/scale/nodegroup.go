package scale

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func scaleNodeGroupCmd(cmd *cmdutils.Cmd) {
	scaleNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup) error {
		return doScaleNodeGroup(cmd, ng)
	})
}

func scaleNodeGroupWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, ng *api.NodeGroup) error) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("nodegroup", "Scale a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, ng)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to scale")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		desiredCapacity := fs.IntP("nodes", "N", -1, "desired number of nodes (required)")
		maxCapacity := fs.IntP("nodes-max", "M", -1, "maximum number of nodes")
		minCapacity := fs.IntP("nodes-min", "m", -1, "minimum number of nodes")

		cmdutils.AddPreRun(cmd.CobraCommand, func(cobraCmd *cobra.Command, args []string) {
			if f := cobraCmd.Flag("nodes"); f.Changed {
				ng.DesiredCapacity = desiredCapacity
			}
			if f := cobraCmd.Flag("nodes-max"); f.Changed {
				ng.MaxSize = maxCapacity
			}
			if f := cobraCmd.Flag("nodes-min"); f.Changed {
				ng.MinSize = minCapacity
			}
		})

		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, true)
}

func doScaleNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup) error {
	if err := cmdutils.NewScaleNodeGroupLoader(cmd, ng).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)
	err = stackManager.ScaleNodeGroup(ng)
	if err != nil {
		return fmt.Errorf("failed to scale nodegroup for cluster %q, error %v", cfg.Metadata.Name, err)
	}

	return nil
}
