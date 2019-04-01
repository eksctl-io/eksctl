package create

import (
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils"
)

var (
	updateAuthConfigMap  bool
	nodeGroupOnlyFilters []string
)

func createNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cfg.Metadata.Version = "auto"

	cmd := &cobra.Command{
		Use:     "nodegroup",
		Short:   "Create a nodegroup",
		Aliases: []string{"ng"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := doCreateNodeGroups(p, cfg, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	exampleNodeGroupName := cmdutils.NodeGroupName("", "")

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to add the nodegroup to")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, `for nodegroups "auto" and "latest" can be used to automatically inherit version from the control plane or force latest`)
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)

		fs.StringSliceVar(&nodeGroupOnlyFilters, "only", nil,
			"select a subset of nodegroups via comma-separated list of globs, e.g.: 'ng-*,nodegroup?,N*group'")

		cmdutils.AddUpdateAuthConfigMap(&updateAuthConfigMap, fs, "Remove nodegroup IAM role from aws-auth configmap")
		cmdutils.AddCommonFlagsForCreateCmd(fs, &output)
	})

	group.InFlagSet("New nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVarP(&ng.Name, "name", "n", "", fmt.Sprintf("name of the new nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		cmdutils.AddCommonCreateNodeGroupFlags(cmd, fs, p, cfg, ng)
	})

	group.InFlagSet("IAM addons", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonCreateNodeGroupIAMAddonsFlags(fs, ng)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.AddTo(cmd)

	return cmd
}

func doCreateNodeGroups(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string, cmd *cobra.Command) error {
	ngFilter := cmdutils.NewNodeGroupFilter()

	if err := cmdutils.NewCreateNodeGroupLoader(p, cfg, clusterConfigFile, nameArg, cmd, ngFilter, nodeGroupOnlyFilters).Load(); err != nil {
		return err
	}

	if err := ngFilter.ValidateNodeGroupsAndSetDefaults(cfg.NodeGroups); err != nil {
		return err
	}

	meta := cfg.Metadata
	printer, err := printers.NewPrinter(output)
	if err != nil {
		return err
	}
	ctl := eks.New(p, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	if err := checkVersion(ctl, cfg.Metadata); err != nil {
		return err
	}

	if err := ctl.GetClusterVPC(cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	}

	stackManager := ctl.NewStackManager(cfg)

	if err := ngFilter.ApplyExistingFilter(stackManager); err != nil {
		return err
	}

	err = ngFilter.CheckEachNodeGroup(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
		// resolve AMI
		if err := ctl.EnsureAMI(meta.Version, ng); err != nil {
			return err
		}
		logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, cfg.Metadata.Version)

		if err := ctl.SetNodeLabels(ng, meta); err != nil {
			return err
		}

		// load or use SSH key - name includes cluster name and the
		// fingerprint, so if unique keys provided, each will get
		// loaded and used as intended and there is no need to have
		// nodegroup name in the key name
		if err := ctl.LoadSSHPublicKey(meta.Name, ng); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := printer.LogObj(logger.Debug, "cfg = \\\n%s\n", cfg); err != nil {
		return err
	}

	if err := ctl.ValidateClusterForCompatibility(cfg, stackManager); err != nil {
		return errors.Wrap(err, "cluster compatibility check failed")
	}

	ngSubset := ngFilter.MatchAll(cfg)
	ngCount := ngSubset.Len()

	{
		ngFilter.LogInfo(cfg)
		if ngCount > 0 {
			logger.Info("will create a CloudFormation stack for each of %d nodegroups in cluster %q", ngCount, cfg.Metadata.Name)
		}

		errs := stackManager.CreateAllNodeGroups(ngSubset, printer)
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and nodegroups haven't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete nodegroup --region=%s --cluster=%s --name=<name>' for each of the failed nodegroup", cfg.Metadata.Region, cfg.Metadata.Name)
			for _, err := range errs {
				if err != nil {
					logger.Critical("%s\n", err.Error())
				}
			}
			return fmt.Errorf("failed to create nodegroups for cluster %q", cfg.Metadata.Name)
		}
	}

	{ // post-creation action
		clientSet, err := ctl.NewStdClientSet(cfg)
		if err != nil {
			return err
		}

		err = ngFilter.CheckEachNodeGroup(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
			if updateAuthConfigMap {
				// authorise nodes to join
				if err = authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
					return err
				}

				// wait for nodes to join
				if err = ctl.WaitForNodes(clientSet, ng); err != nil {
					return err
				}
			}

			// if GPU instance type, give instructions
			if utils.IsGPUInstanceType(ng.InstanceType) {
				logger.Info("as you are using a GPU optimized instance type you will need to install NVIDIA Kubernetes device plugin.")
				logger.Info("\t see the following page for instructions: https://github.com/NVIDIA/k8s-device-plugin")
			}

			return nil
		})
		if err != nil {
			return err
		}
		logger.Success("created %d nodegroup(s) in cluster %q", ngCount, cfg.Metadata.Name)
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	return nil
}
