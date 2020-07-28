package delete

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/spot"
)

func deleteNodeGroupCmd(cmd *cmdutils.Cmd) {
	deleteNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, params *cmdutils.DeleteNodeGroupCmdParams) error {
		return doDeleteNodeGroup(cmd, ng, params)
	})
}

func deleteNodeGroupWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, ng *api.NodeGroup, params *cmdutils.DeleteNodeGroupCmdParams) error) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	params := &cmdutils.DeleteNodeGroupCmdParams{}

	cmd.SetDescription("nodegroup", "Delete a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, ng, params)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&params.OnlyMissing, "only-missing", false, "Only delete nodegroups that are not defined in the given config file")
		cmdutils.AddUpdateAuthConfigMap(fs, &params.UpdateAuthConfigMap, "Remove nodegroup IAM role from aws-auth configmap")
		fs.BoolVar(&params.Drain, "drain", true, "Drain and cordon all nodes in the nodegroup before deletion")

		cmd.Wait = false
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "deletion of all resources")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, true)

	cmd.FlagSetGroup.InFlagSet("Spot", func(fs *pflag.FlagSet) {
		cmdutils.AddSpotOceanCommonFlags(fs, &params.SpotProfile)
		cmdutils.AddSpotOceanDeleteNodeGroupFlags(fs, &params.SpotRoll, &params.SpotRollBatchSize)
	})
}

func doDeleteNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup, params *cmdutils.DeleteNodeGroupCmdParams) error {
	ngFilter := filter.NewNodeGroupFilter()

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
		if params.OnlyMissing {
			err = ngFilter.SetOnlyRemote(stackManager, cfg)
			if err != nil {
				return err
			}
		}
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
					NodeGroupBase: &api.NodeGroupBase{
						Name: ng.Name,
					},
				},
			}
		}
	}

	// Spot Ocean.
	{
		// List all nodegroup stacks.
		stacks, err := stackManager.DescribeNodeGroupStacks()
		if err != nil {
			return err
		}

		// Filter nodegroups.
		nodeGroups := ngFilter.FilterMatching(cfg.NodeGroups)
		nodeGroupsDeleteFilter := spot.NewDeleteIncludedFilter(nodeGroups)

		// Execute pre-deletion actions.
		if err := spot.RunPreDeletion(ctl.Provider, cfg, nodeGroups, stacks,
			nodeGroupsDeleteFilter, params.SpotRoll, params.SpotRollBatchSize, cmd.Plan); err != nil {
			return err
		}

		// Recreate the API client to regenerate the embedded STS token.
		if params.SpotRoll {
			// By default, pre-signed STS URLs are valid for 15 minutes after
			// timestamp in x-amz-date header, which means the actual token
			// expiration is 14 minutes (aws-iam-authenticator sets the token
			// expiration to 1 minute before the pre-signed URL expires for
			// some cushion).  We have to regenerate the token here since
			// rolling one or more nodegroups may take longer to complete.
			clientSet, err = ctl.NewStdClientSet(cfg)
			if err != nil {
				return err
			}
		}

		// Explicitly append Ocean nodegroup to the include filter.
		if cmd.ClusterConfigFile == "" {
			for _, ng := range cfg.NodeGroups {
				if ng.Name == api.SpotOceanNodeGroupName {
					ngFilter.AppendIncludeNames(api.SpotOceanNodeGroupName)
					break
				}
			}
		}
	}

	logFiltered := cmdutils.ApplyFilter(cfg, ngFilter)
	logFiltered()

	if params.UpdateAuthConfigMap {
		for _, ng := range cfg.NodeGroups {
			if ng.IAM == nil || ng.IAM.InstanceRoleARN == "" {
				if err := ctl.GetNodeGroupIAM(stackManager, ng); err != nil {
					logger.Warning("error getting instance role ARN for nodegroup %q: %v", ng.Name, err)
					return nil
				}
			}
		}
	}

	allNodeGroups := cmdutils.ToKubeNodeGroups(cfg)

	if params.Drain && !params.SpotRoll {
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
	}

	if params.UpdateAuthConfigMap {
		cmdutils.LogIntendedAction(cmd.Plan, "delete %d nodegroups from auth ConfigMap in cluster %q", len(cfg.NodeGroups), cfg.Metadata.Name)
		if !cmd.Plan {
			for _, ng := range cfg.NodeGroups {
				if err := authconfigmap.RemoveNodeGroup(clientSet, ng); err != nil {
					logger.Warning(err.Error())
				}
			}
		}
	}

	cmdutils.LogCompletedAction(cmd.Plan, "deleted %d nodegroup(s) from cluster %q", len(allNodeGroups), cfg.Metadata.Name)

	cmdutils.LogPlanModeWarning(cmd.Plan && len(allNodeGroups) > 0)

	return nil
}
