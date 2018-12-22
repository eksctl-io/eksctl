package create

import (
	"fmt"
	"os"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"strings"
	"github.com/weaveworks/eksctl/pkg/utils"
)

func createNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:   "nodegroup",
		Short: "Create a nodegroup",
		Run: func(_ *cobra.Command, args []string) {
			name := cmdutils.GetNameArg(args)
			if name != "" {
				ng.Name = name
			}
			if err := doAddNodeGroup(p, cfg, ng); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "Name of the EKS cluster to add the nodegroup to")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddCFNRoleARNFlag(fs, p)
		fs.StringVar(&cfg.Metadata.Version, "version", api.LatestVersion, fmt.Sprintf("Kubernetes version (valid options: %s)", strings.Join(api.SupportedVersions(), ",")))
	})

	group.InFlagSet("Nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup. Defaults to \"ng-<ID>\"")
		cmdutils.AddCommonCreateNodeGroupFlags(fs, p, cfg, ng)
	})

	cmdutils.AddCommonFlagsForAWS(group, p)

	group.AddTo(cmd)

	return cmd
}

func doAddNodeGroup(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup) error {
	ctl := eks.New(p, cfg)
	meta := cfg.Metadata

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

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

	if err := ctl.EnsureAMI(meta.Version, ng); err != nil {
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
		if ng.Name == "" {
			ng.Name = utils.NodegroupName()
		}
		logger.Info("will create a Cloudformation stack for nodegroup %s for cluster %s", ng.Name, cfg.Metadata.Name)
		errs := stackManager.RunTask(stackManager.CreateNodeGroup, ng)
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and nodegroup hasn't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete nodegroup %s --region=%s --name=%s'", ng.Name, cfg.Metadata.Region, cfg.Metadata.Name)
			for _, err := range errs {
				if err != nil {
					logger.Critical("%s\n", err.Error())
				}
			}
			return fmt.Errorf("failed to create nodegroup %s for cluster %q", ng.Name, cfg.Metadata.Name)
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
	logger.Success("EKS cluster %q in %q region has a new nodegroup with name %d", cfg.Metadata.Name, cfg.Metadata.Region, ng.Name)

	return nil

}
