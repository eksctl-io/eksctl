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
	cfg := cmd.ClusterConfig

	// TODO: move this into a loader when --config-file gets added to this command
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if ng.Name != "" && cmd.NameArg != "" {
		return cmdutils.ErrFlagAndArg("--name", ng.Name, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		ng.Name = cmd.NameArg
	}

	if ng.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}

	if ng.DesiredCapacity == nil || *ng.DesiredCapacity < 0 {
		return fmt.Errorf("number of nodes must be 0 or greater. Use the --nodes/-N flag")
	}

	if ng.MaxSize != nil && *ng.MaxSize < 0 {
		return fmt.Errorf("maximum number of nodes must be 0 or greater. Use the --nodes-max flag")
	}

	if ng.MaxSize != nil && ng.MinSize != nil && (*ng.MinSize > *ng.DesiredCapacity || *ng.MaxSize < *ng.DesiredCapacity) {
		return fmt.Errorf("number of nodes must be within range of min nodes and max nodes")
	}

	if ng.MaxSize != nil && *ng.MaxSize < *ng.DesiredCapacity {
		return fmt.Errorf("maximum number of nodes must be greater than or equal to number of nodes")
	}

	if ng.MinSize != nil && *ng.MinSize < 0 {
		return fmt.Errorf("minimum number of nodes must be 0 or greater. Use the --nodes-min flag")
	}

	if ng.MinSize != nil && *ng.MinSize > *ng.DesiredCapacity {
		return fmt.Errorf("minimum number of nodes must be less than or equal to number of nodes")
	}

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
