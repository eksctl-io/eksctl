package create

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils"

	awseks "github.com/aws/aws-sdk-go/service/eks"
)

func createNodeGroupCmd() *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "nodegroup",
		Short: "Create a nodegroup",
		Run: func(_ *cobra.Command, args []string) {
			if err := doAddNodeGroup(p, cfg, ng, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	addCommonCreateFlags(fs, p, cfg, ng)
	fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "Name of the EKS cluster to add the nodegroup to")

	return cmd
}

func doAddNodeGroup(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup, name string) error {
	ctl := eks.New(p, cfg)

	if cfg.Metadata.Region != api.EKSRegionUSWest2 && cfg.Metadata.Region != api.EKSRegionUSEast1 && cfg.Metadata.Region != api.EKSRegionEUWest1 {
		return fmt.Errorf("%s is not supported only %s, %s and %s are supported", cfg.Metadata.Region, api.EKSRegionUSWest2, api.EKSRegionUSEast1, api.EKSRegionEUWest1)
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if utils.ClusterName(cfg.Metadata.Name, name) == "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.Metadata.Name, name)
	}
	cfg.Metadata.Name = utils.ClusterName(cfg.Metadata.Name, name)

	if ng.SSHPublicKeyPath == "" {
		return fmt.Errorf("--ssh-public-key must be non-empty string")
	}

	//TODO:  do we need to do the AZ stuff from create????

	if err := ctl.EnsureAMI(ng); err != nil {
		return err
	}

	if err := ctl.LoadSSHPublicKey(cfg.Metadata.Name, ng); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	//TODO: is this check needed????
	// Check the cluster exists and is active
	eksCluster, err := ctl.DescribeControlPlane(cfg.Metadata)
	if err != nil {
		return err
	}
	if *eksCluster.Status != awseks.ClusterStatusActive {
		return fmt.Errorf("cluster %s status is %s, its needs to be active to add a nodegroup", *eksCluster.Name, *eksCluster.Status)
	}
	logger.Info("found cluster %s", eksCluster.Name)
	logger.Debug("cluster = %#v", eksCluster)

	{
		stackManager := ctl.NewStackManager(cfg)
		maxSeq, err := stackManager.GetMaxNodeGroupSeq()
		if err != nil {
			return err
		}
		ng.ID = maxSeq + 1
		logger.Info("will create a Cloudformation stack for nodegroup %d for cluster %s", ng.ID, cfg.Metadata.Name)
		errs := stackManager.CreateNodeGroups()
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and nodegroup hasn't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete nodegroup %d --region=%s --name=%s'", ng.ID, cfg.Metadata.Region, cfg.Metadata.Name)
			for _, err := range errs {
				logger.Critical("%s\n", err.Error())
			}
			return fmt.Errorf("failed to create nodegroup %d for cluster %q", ng.ID, cfg.Metadata.Name)
		}
	}

	{ // post-creation action
		clientConfigBase, err := ctl.NewClientConfig(cfg)
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
	logger.Success("EKS cluster %q in %q region has a new nodegroup with with id %d", cfg.Metadata.Name, cfg.Metadata.Region, ng.ID)

	return nil

}
