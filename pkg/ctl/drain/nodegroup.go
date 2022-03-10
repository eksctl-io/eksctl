package drain

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
)

func drainNodeGroupCmd(cmd *cmdutils.Cmd) {
	drainNodeGroupWithRunFunc(cmd, func(cmd *cmdutils.Cmd, ng *api.NodeGroup, undo, onlyMissing bool, maxGracePeriod, nodeDrainWaitPeriod time.Duration, disableEviction bool, parallel int) error {
		return doDrainNodeGroup(cmd, ng, undo, onlyMissing, maxGracePeriod, nodeDrainWaitPeriod, disableEviction, parallel)
	})
}

func drainNodeGroupWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, ng *api.NodeGroup, undo, onlyMissing bool, maxGracePeriod, nodeDrainWaitPeriod time.Duration, disableEviction bool, parallel int) error) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var (
		undo                bool
		onlyMissing         bool
		disableEviction     bool
		parallel            int
		maxGracePeriod      time.Duration
		nodeDrainWaitPeriod time.Duration
	)

	cmd.SetDescription("nodegroup", "Cordon and drain a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, ng, undo, onlyMissing, maxGracePeriod, nodeDrainWaitPeriod, disableEviction, parallel)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to drain")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only drain nodegroups that are not defined in the given config file")
		fs.BoolVar(&undo, "undo", false, "Uncordon the nodegroup")
		defaultMaxGracePeriod, _ := time.ParseDuration("10m")
		fs.DurationVar(&maxGracePeriod, "max-grace-period", defaultMaxGracePeriod, "Maximum pods termination grace period")
		defaultDisableEviction := false
		fs.BoolVar(&disableEviction, "disable-eviction", defaultDisableEviction, "Force drain to use delete, even if eviction is supported. This will bypass checking PodDisruptionBudgets, use with caution.")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		fs.DurationVar(&nodeDrainWaitPeriod, "node-drain-wait-period", 0, "Amount of time to wait between draining nodes in a nodegroup")
		fs.IntVar(&parallel, "parallel", 1, "Number of nodes to drain in parallel. Max 25")
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, true)
}

func doDrainNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup, undo, onlyMissing bool, maxGracePeriod, nodeDrainWaitPeriod time.Duration, disableEviction bool, parallel int) error {
	ngFilter := filter.NewNodeGroupFilter()

	if err := cmdutils.NewDeleteAndDrainNodeGroupLoader(cmd, ng, ngFilter).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewProviderForExistingCluster()
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
			err = ngFilter.SetOnlyRemote(ctl.Provider.EKS(), stackManager, cfg)
			if err != nil {
				return err
			}
		}
	} else {
		err := cmdutils.PopulateNodegroup(stackManager, ng.Name, cfg, ctl.Provider)
		if err != nil {
			return err
		}
	}

	logFiltered := cmdutils.ApplyFilter(cfg, ngFilter)

	verb := "drain"
	if undo {
		verb = "uncordon"
	}

	logAction := func(resource string, count int) {
		cmdutils.LogIntendedAction(cmd.Plan, "%s %d %s in cluster %q", verb, count, resource, cfg.Metadata.Name)
	}
	logFiltered()

	logAction("nodegroup(s)", len(cfg.NodeGroups))
	logAction("managed nodegroup(s)", len(cfg.ManagedNodeGroups))

	cmdutils.LogPlanModeWarning(cmd.Plan && (len(cfg.NodeGroups) > 0 || len(cfg.ManagedNodeGroups) > 0))

	if cmd.Plan {
		return nil
	}
	allNodeGroups := cmdutils.ToKubeNodeGroups(cfg)

	drainInput := &nodegroup.DrainInput{
		NodeGroups:          allNodeGroups,
		Plan:                cmd.Plan,
		MaxGracePeriod:      maxGracePeriod,
		NodeDrainWaitPeriod: nodeDrainWaitPeriod,
		Undo:                undo,
		DisableEviction:     disableEviction,
		Parallel:            parallel,
	}
	return nodegroup.New(cfg, ctl, clientSet).Drain(drainInput)
}
