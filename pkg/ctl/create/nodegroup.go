package create

import (
	"fmt"
	"os"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

func createNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
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

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "cluster", "n", "", "Name of the EKS cluster to add the nodegroup to")
		cmdutils.AddRegionFlag(fs, p)
		fs.StringVar(&p.Version, "version", "1.11", "Kubernetes version (valid options: 1.10, 1.11)")
	})

	group.InFlagSet("Nodegroup", func(fs *pflag.FlagSet) {
		addCommonCreateFlags(fs, p, cfg, ng)
	})

	cmdutils.AddCommonFlagsForAWS(group, p)

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

	if cfg.Metadata.Name == "" {
		return errors.New("--cluster must be specified. run `eksctl get cluster` to show existing clusters")
	}

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
	logger.Info("found cluster %s", *eksCluster.Name)
	logger.Debug("cluster = %#v", eksCluster)

	// Populate cfg with the endopoint, CA data, and so on obtained from the described control-plane
	// So that we won't end up rendering a incomplete useradata missing those things
	if err = ctl.GetCredentials(*eksCluster, cfg); err != nil {
		return err
	}

	{
		stackManager := ctl.NewStackManager(cfg)
		maxSeq, err := stackManager.GetMaxNodeGroupSeq()
		if err != nil {
			return err
		}
		ng.ID = maxSeq + 1
		logger.Info("will create a Cloudformation stack for nodegroup %d for cluster %s", ng.ID, cfg.Metadata.Name)
		errs := stackManager.RunTask(stackManager.CreateNodeGroup, ng)
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and nodegroup hasn't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete nodegroup %d --region=%s --name=%s'", ng.ID, cfg.Metadata.Region, cfg.Metadata.Name)
			for _, err := range errs {
				if err != nil {
					logger.Critical("%s\n", err.Error())
				}
			}
			return fmt.Errorf("failed to create nodegroup %d for cluster %q", ng.ID, cfg.Metadata.Name)
		}
	}

	{ // post-creation action
		clientConfigBase, err := ctl.NewClientConfig(cfg)
		if err != nil {
			return err
		}

		clientConfig := clientConfigBase.WithExecAuthenticator()

		clientSet, err := clientConfig.NewClientSet()
		if err != nil {
			return err
		}

		// authorise nodes to join
		if err = ctl.AddNodeGroupToAuthConfigMap(clientSet, ng); err != nil {
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
