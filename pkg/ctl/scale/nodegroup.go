package scale

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func scaleNodeGroupCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	rc.ClusterConfig = cfg

	rc.SetDescription("nodegroup", "Scale a nodegroup", "", "ng")

	rc.SetRunFuncWithNameArg(func() error {
		return doScaleNodeGroup(rc, ng)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to scale")

		desiredCapacity := fs.IntP("nodes", "N", -1, "total number of nodes (scale to this number)")
		cmdutils.AddPreRun(rc.Command, func(cmd *cobra.Command, args []string) {
			if f := cmd.Flag("nodes"); f.Changed {
				ng.DesiredCapacity = desiredCapacity
			}
		})

		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &rc.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, true)
}

func doScaleNodeGroup(rc *cmdutils.ResourceCmd, ng *api.NodeGroup) error {
	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--cluster")
	}

	if ng.Name != "" && rc.NameArg != "" {
		return cmdutils.ErrNameFlagAndArg(ng.Name, rc.NameArg)
	}

	if rc.NameArg != "" {
		ng.Name = rc.NameArg
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
