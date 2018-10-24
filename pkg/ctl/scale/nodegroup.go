package scale

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

func scaleNodeGroupCmd() *cobra.Command {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "nodegroup",
		Short: "Scale a nodegroup",
		Run: func(_ *cobra.Command, args []string) {
			if err := doScaleNodeGroup(cfg, ng); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "EKS cluster name")

	fs.IntVarP(&ng.DesiredCapacity, "nodes", "N", -1, "total number of nodes (scale to this number)")

	fs.StringVarP(&cfg.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&cfg.Profile, "profile", "p", "", "AWS creditials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.DurationVar(&cfg.WaitTimeout, "timeout", api.DefaultWaitTimeout, "max wait time in any polling operations")

	return cmd
}

func doScaleNodeGroup(cfg *api.ClusterConfig, ng *api.NodeGroup) error {
	ctl := eks.New(cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.ClusterName == "" {
		return fmt.Errorf("no cluster name supplied. Use the --name= flag")
	}

	if ng.DesiredCapacity < 0 {
		return fmt.Errorf("number of nodes must be 0 or greater. Use the --nodes/-N flag")
	}

	stackManager := ctl.NewStackManager()
	err := stackManager.ScaleInitialNodeGroup()
	if err != nil {
		return fmt.Errorf("failed to scale nodegroup for cluster %q, error %v", cfg.ClusterName, err)
	}

	return nil
}
