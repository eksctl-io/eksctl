package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/drain"
)

func deleteNodeGroupCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool

	cmd.SetDescription("nodegroup", "Delete a nodegroup", "", "ng")

	cmd.SetRunFuncWithNameArg(func() error {
		return doDeleteNodeGroup(cmd, ng, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only delete nodegroups that are not defined in the given config file")
		cmdutils.AddUpdateAuthConfigMap(fs, &updateAuthConfigMap, "Remove nodegroup IAM role from aws-auth configmap")
		fs.BoolVar(&deleteNodeGroupDrain, "drain", true, "Drain and cordon all nodes in the nodegroup before deletion")

		cmd.Wait = false
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "deletion of all resources")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, true)
}

func doDeleteNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool) error {
	ngFilter := cmdutils.NewNodeGroupFilter()

	if err := cmdutils.NewDeleteNodeGroupLoader(cmd, ng, ngFilter).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(cfg.Metadata)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	if cmd.ClusterConfigFile != "" {
		logger.Info("comparing %d nodegroups defined in the given config (%q) against remote state", len(cfg.NodeGroups), cmd.ClusterConfigFile)
		if err := ngFilter.SetIncludeOrExcludeMissingFilter(stackManager, onlyMissing, &cfg.NodeGroups); err != nil {
			return err
		}
		// TODO apply filters to managed nodegroups
	} else {
		nodeGroupType, err := stackManager.GetNodeGroupStackType(ng.Name)
		if err != nil {
			return err
		}
		switch nodeGroupType {
		case api.NodeGroupTypeUnmanaged:
			cfg.NodeGroups = []*api.NodeGroup{ng}
		case api.NodeGroupTypeManaged:
			cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					Name: ng.Name,
				},
			}
		}
	}

	logFiltered := cmdutils.ApplyFilter(cfg, ngFilter)

	logFiltered()

	if updateAuthConfigMap {
		cmdutils.LogIntendedAction(cmd.Plan, "delete %d nodegroups from auth ConfigMap in cluster %q", len(cfg.NodeGroups), cfg.Metadata.Name)
		if !cmd.Plan {
			for _, ng := range cfg.NodeGroups {
				if ng.IAM == nil || ng.IAM.InstanceRoleARN == "" {
					if err := ctl.GetNodeGroupIAM(stackManager, cfg, ng); err != nil {
						logger.Warning("error getting instance role ARN for nodegroup %q: %v", ng.Name, err)
						return nil
					}
				}
				if err := authconfigmap.RemoveNodeGroup(clientSet, ng); err != nil {
					logger.Warning(err.Error())
				}
			}
		}
	}

	allNodeGroups := cmdutils.ToKubeNodeGroups(cfg)

	if deleteNodeGroupDrain {
		cmdutils.LogIntendedAction(cmd.Plan, "drain %d nodegroup(s) in cluster %q", len(allNodeGroups), cfg.Metadata.Name)

		if !cmd.Plan {
			for _, ng := range allNodeGroups {
				if err := drain.NodeGroup(clientSet, ng, ctl.Provider.WaitTimeout(), false); err != nil {
					return err
				}
			}
		}
	}

	cmdutils.LogIntendedAction(cmd.Plan, "delete %d nodegroups from cluster %q", len(allNodeGroups), cfg.Metadata.Name)

	{
		shouldDelete := func(ngName string) bool {
			for _, ng := range allNodeGroups {
				if ng.NameString() == ngName {
					return true
				}
			}
			return false
		}

		tasks, err := stackManager.NewTasksToDeleteNodeGroups(shouldDelete, cmd.Wait, nil)
		if err != nil {
			return err
		}
		tasks.PlanMode = cmd.Plan
		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			return handleErrors(errs, "nodegroup(s)")
		}
		cmdutils.LogCompletedAction(cmd.Plan, "deleted %d nodegroup(s) from cluster %q", len(allNodeGroups), cfg.Metadata.Name)
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && len(allNodeGroups) > 0)

	return nil
}
