package scale

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

func scaleNodeGroupCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("nodegroup", "Scale a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doScaleNodeGroup(cmd, ng)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to scale")

		desiredCapacity := fs.IntP("nodes", "N", -1, "total number of nodes (scale to this number)")
		cmdutils.AddPreRun(cmd.CobraCommand, func(cobraCmd *cobra.Command, args []string) {
			if f := cobraCmd.Flag("nodes"); f.Changed {
				ng.DesiredCapacity = desiredCapacity
			}
		})

		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, true)

	cmd.FlagSetGroup.InFlagSet("Kubeconfig context", func(fs *pflag.FlagSet) {
		cmd.KubeconfigPath = ""
		cmd.KubeconfigContext = ""
		cmdutils.AddCommonFlagsForKubeconfigContext(fs, &cmd.KubeconfigPath, &cmd.KubeconfigContext)
	})
}

func doScaleNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup) error {
	cfg := cmd.ClusterConfig

	// TODO: move this into a loader when --config-file gets added to this command
	if cfg.Metadata.Name == "" && cmd.KubeconfigContext == "" {
		return cmdutils.ErrAtLeastOneMustBeSet(cmdutils.ClusterNameFlag(cmd),
			cmdutils.KubeconfigContextFlag(cmd))
	}

	if cmd.KubeconfigContext != "" {
		clusterName, err := kubeconfig.GetClusterNameFromContext(cmd.KubeconfigPath, cmd.KubeconfigContext)
		if err != nil {
			return err
		}
		cfg.Metadata.Name = clusterName
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

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ng.DesiredCapacity == nil || *ng.DesiredCapacity < 0 {
		return fmt.Errorf("number of nodes must be 0 or greater. Use the --nodes/-N flag")
	}

	stackManager := ctl.NewStackManager(cfg)
	err = stackManager.ScaleNodeGroup(ng)
	if err != nil {
		return fmt.Errorf("failed to scale nodegroup for cluster %q, error %v", cfg.Metadata.Name, err)
	}

	return nil
}
