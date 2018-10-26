package create

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils"

	awseks "github.com/aws/aws-sdk-go/service/eks"
)

func createNodeGroupCmd() *cobra.Command {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "nodegroup",
		Short: "Create a nodegroup",
		Run: func(_ *cobra.Command, args []string) {
			if err := doAddNodeGroup(cfg, ng, ctl.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	addCommonCreateFlags(fs, cfg, ng)
	fs.StringVarP(&cfg.ClusterName, "name", "n", "", "Name of the EKS cluster to add the nodegroup to")

	return cmd
}

func doAddNodeGroup(cfg *api.ClusterConfig, ng *api.NodeGroup, name string) error {
	ctl := eks.New(cfg)

	if cfg.Region != api.EKSRegionUSWest2 && cfg.Region != api.EKSRegionUSEast1 && cfg.Region != api.EKSRegionEUWest1 {
		return fmt.Errorf("%s is not supported only %s, %s and %s are supported", cfg.Region, api.EKSRegionUSWest2, api.EKSRegionUSEast1, api.EKSRegionEUWest1)
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if utils.ClusterName(cfg.ClusterName, name) == "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.ClusterName, name)
	}
	cfg.ClusterName = utils.ClusterName(cfg.ClusterName, name)

	if ng.SSHPublicKeyPath == "" {
		return fmt.Errorf("--ssh-public-key must be non-empty string")
	}

	//TODO:  do we need to do the AZ stuff from create????

	if err := ctl.EnsureAMI(ng); err != nil {
		return err
	}

	if err := ctl.LoadSSHPublicKey(ng); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	//TODO: is this check needed????
	// Check the cluster exists and is active
	eksCluster, err := ctl.DescribeControlPlane()
	if err != nil {
		return err
	}
	if *eksCluster.Status != awseks.ClusterStatusActive {
		return fmt.Errorf("cluster %s status is %s, its needs to be active to add a nodegroup", *eksCluster.Name, *eksCluster.Status)
	}
	logger.Info("found cluster %s", eksCluster.Name)
	logger.Debug("cluster = %#v", eksCluster)

	{
		stackManager := ctl.NewStackManager()
		maxSeq, err := stackManager.GetMaxNodeGroupSeq()
		if err != nil {
			return err
		}
		ng.ID = maxSeq + 1
		logger.Info("will create a Cloudformation stack for nodegroup %d for cluster %s", ng.ID, cfg.ClusterName)
		errs := stackManager.CreateNodeGroups()
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and nodegroup hasn't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete nodegroup %d --region=%s --name=%s'", ng.ID, cfg.Region, cfg.ClusterName)
			for _, err := range errs {
				logger.Critical("%s\n", err.Error())
			}
			return fmt.Errorf("failed to create nodegroup %d for cluster %q", ng.ID, cfg.ClusterName)
		}
	}

	{ // post-creation action
		clientConfigBase, err := ctl.NewClientConfig()
		if err != nil {
			return err
		}

		clientSet, err := clientConfigBase.NewClientSetWithEmbeddedToken()
		if err != nil {
			return err
		}

		// authorise nodes to join
		if err = ctl.CreateNodeGroupAuthConfigMap(clientSet, ng); err != nil {
			return err
		}

		// wait for nodes to join
		if err = ctl.WaitForNodes(clientSet, ng); err != nil {
			return err
		}
	}
	logger.Success("EKS cluster %q in %q region has a new nodegroup with with id %d", cfg.ClusterName, cfg.Region, ng.ID)

	return nil

}
