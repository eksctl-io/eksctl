package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func deleteNodeGroupCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool

	cmd.SetDescription("nodegroup", "Delete a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doDeleteNodeGroup(cmd, ng, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only delete nodegroups that are not defined in the given config file")
		cmdutils.AddUpdateAuthConfigMap(fs, &updateAuthConfigMap, "Remove nodegroup IAM role from aws-auth configmap when appropriate after stack deletion; requires 'wait' to be set")
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

	descriptions, err := stackManager.DescribeNodeGroupStacks()
	if err != nil {
		return err
	}

	stacks, err := stackManager.MakeNodeGroupStacksFromDescriptions(descriptions)
	if err != nil {
		return err
	}

	if cmd.ClusterConfigFile != "" {
		logger.Info("comparing %d nodegroups defined in the given config (%q) against remote state", len(cfg.NodeGroups), cmd.ClusterConfigFile)
		if err := ngFilter.SetIncludeOrExcludeMissingStackFilter(stacks, onlyMissing, cfg); err != nil {
			return err
		}
	} else {
		var nodeGroupType api.NodeGroupType
		for _, s := range stacks {
			if s.NodeGroupName == ng.Name {
				nodeGroupType = s.Type
				break
			}
		}
		nodeGroupType, err = stackManager.GetNodeGroupStackType(ng.Name)
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

		// NOTE: there is an issue in AWS where if a managed node group is using an IAM instance role also found in an
		//   unmanaged node group, it will delete all identity entries in the auth configmap and possibly orphaning the
		//   workers without having us call the removeARN function below to do so
		tasks, err := stackManager.MakeTasksToDeleteNodeGroupsFromDescriptions(shouldDelete, cmd.Wait, nil, descriptions)
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

	// needs to be ran after the nodegroup deletion since it marks them for deletion before we pass it to
	// removeARN; once it reaches here it should have been deleted already so that it's safe for us to delete
	//
	// if we did not add the wait check then node groups to be deleted can have their stacks fail to delete and we would
	// prematurely remove the identity from the auth configmap and lead to orphaned nodes until the auth is added back
	if updateAuthConfigMap && cmd.Wait {
		if err := removeARN(descriptions, stackManager, cfg, ctl, cmd.Plan, clientSet); err != nil {
			return err
		}
	}

	cmdutils.LogPlanModeWarning(cmd.Plan && len(allNodeGroups) > 0)

	return nil
}

// removeARN takes a look at what managed and unmanaged node groups were marked for deletion and removes their identity
// from the auth configmap
func removeARN(descriptions []*manager.Stack, stackManager *manager.StackCollection, cfgMarkedForDeletion *api.ClusterConfig,
	ctl *eks.ClusterProvider, cmdPlan bool, clientSet kubernetes.Interface) error {

	numToDelete := 0
	for _, n := range cfgMarkedForDeletion.NodeGroups {
		if n.IAM == nil || n.IAM.InstanceRoleARN == "" {
			if err := ctl.PopulateNodeGroupIAMFromDescriptions(stackManager, cfgMarkedForDeletion, n, descriptions); err != nil {
				// we want to return the error instead of logging as we don't want to delete something prematurely
				return err
			}
		}
		numToDelete++
	}

	// deletion of the identity role arn in the auth configmap for managed node groups are
	// currently handled by AWS when deleting the managed node group

	cmdutils.LogIntendedAction(cmdPlan, "delete %d identity role ARNs from auth ConfigMap in cluster %q", numToDelete, cfgMarkedForDeletion.Metadata.Name)

	tasks, err := stackManager.NewTasksToDeleteIdentityRoleARNFromAuthConfigMap(cfgMarkedForDeletion, clientSet)
	if err != nil {
		return err
	}
	tasks.PlanMode = cmdPlan
	logger.Info(tasks.Describe())
	if errs := tasks.DoAllSync(); len(errs) > 0 {
		return handleErrors(errs, "identities")
	}

	return nil
}
