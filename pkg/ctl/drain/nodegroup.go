package drain

import (
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	"github.com/weaveworks/eksctl/pkg/drain"
)

func drainNodeGroupCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cmd.ClusterConfig = cfg

	var undo, onlyMissing bool

	cmd.SetDescription("nodegroup", "Cordon and drain a nodegroup", "", "ng")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doDrainNodeGroup(cmd, ng, undo, onlyMissing)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only drain nodegroups that are not defined in the given config file")
		fs.BoolVar(&undo, "undo", false, "Uncordone the nodegroup")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, true)
}

func doDrainNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup, undo, onlyMissing bool) error {
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
		if err := ngFilter.SetIncludeOrExcludeMissingFilter(stackManager, onlyMissing, cfg); err != nil {
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
	for _, ng := range allNodeGroups {
		if err := drain.NodeGroup(clientSet, ng, ctl.Provider.WaitTimeout(), undo); err != nil {
			return err
		}
	}
	return nil
}
