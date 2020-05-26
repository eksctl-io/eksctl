package create

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/ssh"
	"github.com/weaveworks/eksctl/pkg/utils/names"
	"github.com/weaveworks/eksctl/pkg/vpc"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

type createNodeGroupParams struct {
	updateAuthConfigMap       bool
	managed                   bool
	installNeuronDevicePlugin bool
}

func createNodeGroupCmd(cmd *cmdutils.Cmd) {
	createNodeGroupCmdWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, params createNodeGroupParams) error {
		return doCreateNodeGroups(cmd, ng, params)
	})
}

func createNodeGroupCmdWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, ng *api.NodeGroup, params createNodeGroupParams) error) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var params createNodeGroupParams

	cfg.Metadata.Version = "auto"

	cmd.SetDescription("nodegroup", "Create a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, ng, params)
	}

	exampleNodeGroupName := names.ForNodeGroup("", "")

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to add the nodegroup to")
		cmdutils.AddStringToStringVarPFlag(fs, &cfg.Metadata.Tags, "tags", "", map[string]string{}, "Used to tag the AWS resources")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddVersionFlag(fs, cfg.Metadata, `for nodegroups "auto" and "latest" can be used to automatically inherit version from the control plane or force latest`)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		cmdutils.AddUpdateAuthConfigMap(fs, &params.updateAuthConfigMap, "Add nodegroup IAM role to aws-auth configmap")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmd.FlagSetGroup.InFlagSet("New nodegroup", func(fs *pflag.FlagSet) {
		fs.StringVarP(&ng.Name, "name", "n", "", fmt.Sprintf("name of the new nodegroup (generated if unspecified, e.g. %q)", exampleNodeGroupName))
		cmdutils.AddCommonCreateNodeGroupFlags(fs, cmd, ng)
		fs.BoolVarP(&params.managed, "managed", "", false, "Create EKS-managed nodegroup")
	})

	cmd.FlagSetGroup.InFlagSet("Addons", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonCreateNodeGroupIAMAddonsFlags(fs, ng)
		fs.BoolVarP(&params.installNeuronDevicePlugin, "install-neuron-plugin", "", true, "Install Neuron plugin for Inferentia nodes")
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, true)
}

func doCreateNodeGroups(cmd *cmdutils.Cmd, ng *api.NodeGroup, params createNodeGroupParams) error {
	ngFilter := cmdutils.NewNodeGroupFilter()

	if err := cmdutils.NewCreateNodeGroupLoader(cmd, ng, ngFilter, params.managed).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	printer := printers.NewJSONPrinter()

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(meta)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	if err := checkVersion(cmd, ctl, cfg.Metadata); err != nil {
		return err
	}

	if err := ctl.LoadClusterVPC(cfg); err != nil {
		return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
	}

	stackManager := ctl.NewStackManager(cfg)

	if err := ngFilter.SetExcludeExistingFilter(stackManager); err != nil {
		return err
	}

	logFiltered := cmdutils.ApplyFilter(cfg, ngFilter)

	// EKS 1.14 clusters created with prior versions of eksctl may not support Managed Nodes
	supportsManagedNodes, err := ctl.SupportsManagedNodes(cfg)
	if err != nil {
		return err
	}

	if len(cfg.ManagedNodeGroups) > 0 && !supportsManagedNodes {
		return errors.New("Managed Nodegroups are not supported for this cluster version. Please update the cluster before adding managed nodegroups")
	}

	if err := eks.ValidateBottlerocketSupport(ctl.ControlPlaneVersion(), cmdutils.ToKubeNodeGroups(cfg)); err != nil {
		return err
	}

	for _, ng := range cfg.NodeGroups {
		// resolve AMI
		if err := eks.EnsureAMI(ctl.Provider, meta.Version, ng); err != nil {
			return err
		}
		logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, cfg.Metadata.Version)

		// load or use SSH key - name includes cluster name and the
		// fingerprint, so if unique keys provided, each will get
		// loaded and used as intended and there is no need to have
		// nodegroup name in the key name
		publicKeyName, err := ssh.LoadKey(ng.SSH, meta.Name, ng.Name, ctl.Provider.EC2())
		if err != nil {
			return err
		}
		if publicKeyName != "" {
			ng.SSH.PublicKeyName = &publicKeyName
		}
	}

	managedService := eks.NewNodeGroupService(cfg, ctl.Provider.EC2())
	if err := managedService.NormalizeManaged(cfg.ManagedNodeGroups); err != nil {
		return err
	}

	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	// TODO
	if err := ctl.ValidateClusterForCompatibility(cfg, stackManager); err != nil {
		return errors.Wrap(err, "cluster compatibility check failed")
	}

	if err := vpc.ValidateLegacySubnetsForNodeGroups(cfg, ctl.Provider); err != nil {
		return err
	}

	{
		logFiltered()
		logMsg := func(resource string, count int) {
			logger.Info("will create a CloudFormation stack for each of %d %s in cluster %q", count, resource, cfg.Metadata.Name)
		}
		if len(cfg.NodeGroups) > 0 {
			logMsg("nodegroups", len(cfg.NodeGroups))
		}

		if len(cfg.ManagedNodeGroups) > 0 {
			logMsg("managed nodegroups", len(cfg.ManagedNodeGroups))
		}

		tasks := &manager.TaskTree{
			Parallel: false,
		}
		if supportsManagedNodes {
			tasks.Append(stackManager.NewClusterCompatTask())
		}

		allNodeGroupTasks := &manager.TaskTree{
			Parallel: true,
		}
		nodeGroupTasks := stackManager.NewTasksToCreateNodeGroups(cfg.NodeGroups, supportsManagedNodes)
		if nodeGroupTasks.Len() > 0 {
			allNodeGroupTasks.Append(nodeGroupTasks)
		}
		managedTasks := stackManager.NewManagedNodeGroupTask(cfg.ManagedNodeGroups)
		if managedTasks.Len() > 0 {
			allNodeGroupTasks.Append(managedTasks)
		}

		tasks.Append(allNodeGroupTasks)
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

		tasks := ctl.ClusterTasksForNodeGroups(cfg, params.installNeuronDevicePlugin)
		logger.Info(tasks.Describe())
		errs := tasks.DoAllSync()
		if len(errs) > 0 {
			logger.Info("%d error(s) occurred and nodegroups haven't been created properly, you may wish to check CloudFormation console", len(errs))
			logger.Info("to cleanup resources, run 'eksctl delete nodegroup --region=%s --cluster=%s --name=<name>' for each of the failed nodegroups", cfg.Metadata.Region, cfg.Metadata.Name)
			for _, err := range errs {
				if err != nil {
					logger.Critical("%s\n", err.Error())
				}
			}
			return fmt.Errorf("failed to create nodegroups for cluster %q", cfg.Metadata.Name)
		}

		for _, ng := range cfg.NodeGroups {
			if params.updateAuthConfigMap {
				// authorise nodes to join
				if err = authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
					return err
				}

				// wait for nodes to join
				if err = ctl.WaitForNodes(clientSet, ng); err != nil {
					return err
				}
			}

			showDevicePluginMessageForNodeGroup(ng, params.installNeuronDevicePlugin)
		}
		logger.Success("created %d nodegroup(s) in cluster %q", len(cfg.NodeGroups), cfg.Metadata.Name)

		for _, ng := range cfg.ManagedNodeGroups {
			if err := ctl.WaitForNodes(clientSet, ng); err != nil {
				return err
			}
		}

		logger.Success("created %d managed nodegroup(s) in cluster %q", len(cfg.ManagedNodeGroups), cfg.Metadata.Name)
	}

	if err := ctl.ValidateExistingNodeGroupsForCompatibility(cfg, stackManager); err != nil {
		logger.Critical("failed checking nodegroups", err.Error())
	}

	return nil
}
