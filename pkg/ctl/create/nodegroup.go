package create

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils"
)

func createNodeGroupCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	rc.ClusterConfig = cfg

	var updateAuthConfigMap bool

	cfg.Metadata.Version = "auto"

	rc.SetDescription("nodegroup", "Create a nodegroup", "", "ng")

	rc.SetRunFuncWithNameArg(func() error {
		return doCreateNodeGroups(rc, updateAuthConfigMap)
	})

	exampleNodeGroupName := cmdutils.NodeGroupName("", "")

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to add the nodegroup to")
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, `for nodegroups "auto" and "latest" can be used to automatically inherit version from the control plane or force latest`)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
		cmdutils.AddNodeGroupFilterFlags(fs, &rc.IncludeNodeGroups, &rc.ExcludeNodeGroups)
		cmdutils.AddUpdateAuthConfigMap(fs, &updateAuthConfigMap, "Remove nodegroup IAM role from aws-auth configmap")
		cmdutils.AddTimeoutFlag(fs, &rc.ProviderConfig.WaitTimeout)
	})

	rc.FlagSetGroup.InFlagSet("New nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVarP(&ng.Name, "name", "n", "", fmt.Sprintf("name of the new nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		cmdutils.AddCommonCreateNodeGroupFlags(fs, rc, ng)
	})

	rc.FlagSetGroup.InFlagSet("IAM addons", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonCreateNodeGroupIAMAddonsFlags(fs, ng)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, true)
}

func doCreateNodeGroups(rc *cmdutils.ResourceCmd, updateAuthConfigMap bool) error {
	ngFilter := cmdutils.NewNodeGroupFilter()

	if err := cmdutils.NewCreateNodeGroupLoader(rc, ngFilter).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig
	meta := rc.ClusterConfig.Metadata

	api.SetClusterConfigDefaults(cfg)

	if err := api.ValidateClusterConfig(cfg); err != nil {
		return err
	}

	if err := ngFilter.ValidateNodeGroupsAndSetDefaults(cfg.NodeGroups); err != nil {
		return err
	}

	printer := printers.NewJSONPrinter()
	ctl := eks.New(rc.ProviderConfig, cfg)

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(rc.ProviderConfig)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.RefreshClusterConfig(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	if err := checkVersion(rc, ctl, cfg.Metadata); err != nil {
		return err
	}

	if err := ctl.GetClusterVPC(cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	}

	stackManager := ctl.NewStackManager(cfg)

	if err := ngFilter.SetExcludeExistingFilter(stackManager); err != nil {
		return err
	}

	err := ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
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
		if err := loadSSHKey(ng, meta.Name, ctl.Provider); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	if err := ctl.ValidateClusterForCompatibility(cfg, stackManager); err != nil {
		return errors.Wrap(err, "cluster compatibility check failed")
	}

	ngSubset, _ := ngFilter.MatchAll(cfg.NodeGroups)
	ngCount := ngSubset.Len()

	{
		ngFilter.LogInfo(cfg.NodeGroups)
		if ngCount > 0 {
			logger.Info("will create a CloudFormation stack for each of %d nodegroups in cluster %q", ngCount, cfg.Metadata.Name)
		}

		tasks := stackManager.NewTasksToCreateNodeGroups(ngSubset)
		logger.Info(tasks.Describe())
		errs := tasks.DoAllSync()
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

		err = ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
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
			if utils.IsGPUInstanceType(ng.InstanceType) || (ng.InstancesDistribution != nil && utils.HasGPUInstanceType(ng.InstancesDistribution.InstanceTypes)) {
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
