package delete

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
)

type deleteNodeGroupOptions struct {
	updateAuthConfigMap   *bool
	deleteNodeGroupDrain  bool
	onlyMissing           bool
	maxGracePeriod        time.Duration
	podEvictionWaitPeriod time.Duration
	disableEviction       bool
	parallel              int
}

func deleteNodeGroupCmd(cmd *cmdutils.Cmd) {
	deleteNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, options deleteNodeGroupOptions) error {
		return doDeleteNodeGroup(cmd, ng, options)
	})
}

func deleteNodeGroupWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, ng *api.NodeGroup, options deleteNodeGroupOptions) error) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var options deleteNodeGroupOptions

	cmd.SetDescription("nodegroup", "Delete a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, ng, options)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&options.onlyMissing, "only-missing", false, "Only delete nodegroups that are not defined in the given config file")
		options.updateAuthConfigMap = cmdutils.AddUpdateAuthConfigMap(fs, "Remove nodegroup IAM role from aws-auth configmap")
		fs.BoolVar(&options.deleteNodeGroupDrain, "drain", true, "Drain and cordon all nodes in the nodegroup before deletion")
		defaultMaxGracePeriod, _ := time.ParseDuration("10m")
		fs.DurationVar(&options.maxGracePeriod, "max-grace-period", defaultMaxGracePeriod, "Maximum pods termination grace period")
		defaultPodEvictionWaitPeriod, _ := time.ParseDuration("10s")
		fs.DurationVar(&options.podEvictionWaitPeriod, "pod-eviction-wait-period", defaultPodEvictionWaitPeriod, "Duration to wait after failing to evict a pod")
		fs.BoolVar(&options.disableEviction, "disable-eviction", false, "Force drain to use delete, even if eviction is supported. This will bypass checking PodDisruptionBudgets, use with caution.")
		fs.IntVar(&options.parallel, "parallel", 1, "Number of nodes to drain in parallel. Max 25")

		cmd.Wait = false
		cmdutils.AddWaitFlag(fs, &cmd.Wait, "deletion of all resources")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, true)
}

type authConfigMapUpdater struct {
	clientSet kubernetes.Interface
}

func (a *authConfigMapUpdater) RemoveNodeGroup(ng *api.NodeGroup) error {
	return authconfigmap.RemoveNodeGroup(a.clientSet, ng)
}

func doDeleteNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup, options deleteNodeGroupOptions) error {
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
		if options.onlyMissing {
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

	if api.IsEnabled(options.updateAuthConfigMap) {
		for _, ng := range cfg.NodeGroups {
			if ng.IAM == nil || ng.IAM.InstanceRoleARN == "" {
				if err := ctl.GetNodeGroupIAM(ctx, stackManager, ng); err != nil {
					err := fmt.Sprintf("error getting instance role ARN for nodegroup %q: %v", ng.Name, err)
					logger.Warning("continuing with deletion, error occurred: %s", err)
				}
			}
		}
	}
	allNodeGroups := cmdutils.ToKubeNodeGroups(cfg.NodeGroups, cfg.ManagedNodeGroups)

	if options.deleteNodeGroupDrain {
		cmdutils.LogIntendedAction(cmd.Plan, "drain %d nodegroup(s) in cluster %q", len(allNodeGroups), cfg.Metadata.Name)

		drainInput := &nodegroup.DrainInput{
			NodeGroups:            allNodeGroups,
			Plan:                  cmd.Plan,
			MaxGracePeriod:        options.maxGracePeriod,
			PodEvictionWaitPeriod: options.podEvictionWaitPeriod,
			DisableEviction:       options.disableEviction,
			Parallel:              options.parallel,
		}
		drainCtx, cancel := context.WithTimeout(ctx, cmd.ProviderConfig.WaitTimeout)
		defer cancel()

		drainer := &nodegroup.Drainer{
			ClientSet: clientSet,
		}
		if err := drainer.Drain(drainCtx, drainInput); err != nil {
			logger.Warning("error occurred during drain, to skip drain use '--drain=false' flag")
			return err
		}
	}

	cmdutils.LogIntendedAction(cmd.Plan, "delete %d nodegroups from cluster %q", len(allNodeGroups), cfg.Metadata.Name)

	deleter := &nodegroup.Deleter{
		StackHelper:      stackManager,
		NodeGroupDeleter: ctl.AWSProvider.EKS(),
		ClusterName:      cfg.Metadata.Name,
		AuthConfigMapUpdater: &authConfigMapUpdater{
			clientSet: clientSet,
		},
	}
	if err := deleter.Delete(ctx, cfg.NodeGroups, cfg.ManagedNodeGroups, nodegroup.DeleteOptions{
		Wait:                cmd.Wait,
		Plan:                cmd.Plan,
		UpdateAuthConfigMap: !api.IsDisabled(options.updateAuthConfigMap),
	}); err != nil {
		return err
	}

	cmdutils.LogCompletedAction(cmd.Plan, "deleted %d nodegroup(s) from cluster %q", len(allNodeGroups), cfg.Metadata.Name)
	cmdutils.LogPlanModeWarning(cmd.Plan && len(allNodeGroups) > 0)
	return nil
}
