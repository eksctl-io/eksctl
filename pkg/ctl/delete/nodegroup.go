package delete

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
)

func deleteNodeGroupCmd(cmd *cmdutils.Cmd) {
	deleteNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool, maxGracePeriod, podEvictionWaitPeriod time.Duration, disableEviction bool, parallel int) error {
		return doDeleteNodeGroup(cmd, ng, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing, maxGracePeriod, podEvictionWaitPeriod, disableEviction, parallel)
	})
}

func deleteNodeGroupWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, ng *api.NodeGroup, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool, maxGracePeriod, podEvictionWaitPeriod time.Duration, disableEviction bool, parallel int) error) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var (
		updateAuthConfigMap   bool
		deleteNodeGroupDrain  bool
		onlyMissing           bool
		maxGracePeriod        time.Duration
		podEvictionWaitPeriod time.Duration
		disableEviction       bool
		parallel              int
	)

	cmd.SetDescription("nodegroup", "Delete a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, ng, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing, maxGracePeriod, podEvictionWaitPeriod, disableEviction, parallel)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only delete nodegroups that are not defined in the given config file")
		cmdutils.AddUpdateAuthConfigMap(fs, &updateAuthConfigMap, "Remove nodegroup IAM role from aws-auth configmap")
		fs.BoolVar(&deleteNodeGroupDrain, "drain", true, "Drain and cordon all nodes in the nodegroup before deletion")
		defaultMaxGracePeriod, _ := time.ParseDuration("10m")
		fs.DurationVar(&maxGracePeriod, "max-grace-period", defaultMaxGracePeriod, "Maximum pods termination grace period")
		defaultPodEvictionWaitPeriod, _ := time.ParseDuration("10s")
		fs.DurationVar(&podEvictionWaitPeriod, "pod-eviction-wait-period", defaultPodEvictionWaitPeriod, "Duration to wait after failing to evict a pod")
		defaultDisableEviction := false
		fs.BoolVar(&disableEviction, "disable-eviction", defaultDisableEviction, "Force drain to use delete, even if eviction is supported. This will bypass checking PodDisruptionBudgets, use with caution.")
		fs.IntVar(&parallel, "parallel", 1, "Number of nodes to drain in parallel. Max 25")

		cmd.Wait = false
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "deletion of all resources")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, true)
}

func doDeleteNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool, maxGracePeriod time.Duration, podEvictionWaitPeriod time.Duration, disableEviction bool, parallel int) error {
	ngFilter := filter.NewNodeGroupFilter()

	if err := cmdutils.NewDeleteAndDrainNodeGroupLoader(cmd, ng, ngFilter).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
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
		if onlyMissing {
			err = ngFilter.SetOnlyRemote(ctx, ctl.AWSProvider.EKS(), stackManager, cfg)
			if err != nil {
				return err
			}
		}
	} else {
		err := cmdutils.PopulateNodegroup(ctx, stackManager, ng.Name, cfg, ctl.AWSProvider)
		if err != nil {
			return err
		}
	}

	logFiltered := cmdutils.ApplyFilter(cfg, ngFilter)

	logFiltered()

	if updateAuthConfigMap {
		for _, ng := range cfg.NodeGroups {
			if ng.IAM == nil || ng.IAM.InstanceRoleARN == "" {
				if err := ctl.GetNodeGroupIAM(ctx, stackManager, ng); err != nil {
					err := fmt.Sprintf("error getting instance role ARN for nodegroup %q: %v", ng.Name, err)
					logger.Warning("continuing with deletion, error occurred: %s", err)
				}
			}
		}
	}
	allNodeGroups := cmdutils.ToKubeNodeGroups(cfg)

	nodeGroupManager := nodegroup.New(cfg, ctl, clientSet, selector.New(ctl.AWSProvider.Session()))
	if deleteNodeGroupDrain {
		cmdutils.LogIntendedAction(cmd.Plan, "drain %d nodegroup(s) in cluster %q", len(allNodeGroups), cfg.Metadata.Name)

		drainInput := &nodegroup.DrainInput{
			NodeGroups:            allNodeGroups,
			Plan:                  cmd.Plan,
			MaxGracePeriod:        maxGracePeriod,
			PodEvictionWaitPeriod: podEvictionWaitPeriod,
			DisableEviction:       disableEviction,
			Parallel:              parallel,
		}
		ctx, cancel := context.WithTimeout(ctx, cmd.ProviderConfig.WaitTimeout)
		defer cancel()
		err := nodeGroupManager.Drain(ctx, drainInput)
		if err != nil {
			logger.Warning("error occurred during drain, to skip drain use '--drain=false' flag")
			return err
		}
	}

	cmdutils.LogIntendedAction(cmd.Plan, "delete %d nodegroups from cluster %q", len(allNodeGroups), cfg.Metadata.Name)

	err = nodeGroupManager.Delete(ctx, cfg.NodeGroups, cfg.ManagedNodeGroups, cmd.Wait, cmd.Plan)
	if err != nil {
		return err
	}

	if updateAuthConfigMap {
		cmdutils.LogIntendedAction(cmd.Plan, "delete %d nodegroups from auth ConfigMap in cluster %q", len(cfg.NodeGroups), cfg.Metadata.Name)
		if !cmd.Plan {
			for _, ng := range cfg.NodeGroups {
				if ng.IAM != nil && ng.IAM.InstanceRoleARN != "" {
					if err := authconfigmap.RemoveNodeGroup(clientSet, ng); err != nil {
						logger.Warning(err.Error())
					}
				}
			}
		}
	}

	cmdutils.LogCompletedAction(cmd.Plan, "deleted %d nodegroup(s) from cluster %q", len(allNodeGroups), cfg.Metadata.Name)

	cmdutils.LogPlanModeWarning(cmd.Plan && len(allNodeGroups) > 0)

	return nil
}
