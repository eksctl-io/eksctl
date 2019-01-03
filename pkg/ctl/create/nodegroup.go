package create

import (
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
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
			if err := doCreateNodeGroup(p, cfg, ng, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	exampleNodeGroupName := utils.NodeGroupName("", "")

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to add the nodegroup to")
		cmdutils.AddRegionFlag(fs, p)
		fs.StringVar(&cfg.Metadata.Version, "version", api.LatestVersion, fmt.Sprintf("Kubernetes version (valid options: %s)", strings.Join(api.SupportedVersions(), ",")))
	})

	group.InFlagSet("New nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVarP(&ng.Name, "name", "n", "", fmt.Sprintf("name of the new nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		cmdutils.AddCommonCreateNodeGroupFlags(fs, p, cfg, ng)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.AddTo(cmd)

	return cmd
}

func doCreateNodeGroup(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup, nameArg string) error {
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
		return errors.New("--cluster must be set")
	}

	if utils.NodeGroupName(ng.Name, nameArg) == "" {
		return cmdutils.ErrNameFlagAndArg(ng.Name, nameArg)
	}
	ng.Name = utils.NodeGroupName(ng.Name, nameArg)

	if ng.SSHPublicKeyPath == "" {
		return fmt.Errorf("--ssh-public-key must be non-empty string")
	}

	if err := ctl.EnsureAMI(meta.Version, ng); err != nil {
		return err
	}
	logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, cfg.Metadata.Version)

	if err := ctl.SetNodeLabels(ng, meta); err != nil {
		return err
	}

	if err := ctl.LoadSSHPublicKey(cfg.Metadata.Name, ng); err != nil {
		return err
	}

	logger.Debug("cfg = %#v", cfg)

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	{
		stackManager := ctl.NewStackManager(cfg)
		logger.Info("will create a Cloudformation stack for nodegroup %s in cluster %s", ng.Name, cfg.Metadata.Name)
		errs := stackManager.CreateOneNodeGroup(ng)
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
		if err = ctl.CreateOrUpdateNodeGroupAuthConfigMap(clientSet, ng); err != nil {
			return err
		}

		// wait for nodes to join
		if err = ctl.WaitForNodes(clientSet, ng); err != nil {
			return err
		}
	}
	logger.Success("created nodegroup %q in cluster %q", ng.Name, cfg.Metadata.Name)

	return nil

}
