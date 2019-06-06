package drain

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"

	"github.com/weaveworks/eksctl/pkg/drain"
)

var (
	drainNodeGroupUndo bool

	includeNodeGroups []string
	excludeNodeGroups []string

	drainOnlyMissingNodeGroups bool
)

func drainNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:     "nodegroup",
		Short:   "Cordon and drain a nodegroup",
		Aliases: []string{"ng"},
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doDrainNodeGroup(cp, ng); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &cp.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cp)
		cmdutils.AddNodeGroupFilterFlags(fs, &includeNodeGroups, &excludeNodeGroups)
		fs.BoolVar(&drainOnlyMissingNodeGroups, "only-missing", false, "Only drain nodegroups that are not defined in the given config file")
		fs.BoolVar(&drainNodeGroupUndo, "undo", false, "Uncordone the nodegroup")
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, true)

	group.AddTo(cp.Command)
	return cp.Command
}

func doDrainNodeGroup(cp *cmdutils.CommonParams, ng *api.NodeGroup) error {
	ngFilter := cmdutils.NewNodeGroupFilter()

	if err := cmdutils.NewDeleteNodeGroupLoader(cp, ng, ngFilter, includeNodeGroups, excludeNodeGroups).Load(); err != nil {
		return err
	}

	cfg := cp.ClusterConfig

	ctl := eks.New(cp.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.GetCredentials(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	if cp.ClusterConfigFile != "" {
		logger.Info("comparing %d nodegroups defined in the given config (%q) against remote state", len(cfg.NodeGroups), cp.ClusterConfigFile)
		if err := ngFilter.SetIncludeOrExcludeMissingFilter(stackManager, drainOnlyMissingNodeGroups, &cfg.NodeGroups); err != nil {
			return err
		}
	}

	ngSubset, _ := ngFilter.MatchAll(cfg.NodeGroups)
	ngCount := ngSubset.Len()

	ngFilter.LogInfo(cfg.NodeGroups)
	verb := "drain"
	if drainNodeGroupUndo {
		verb = "uncordon"
	}
	cmdutils.LogIntendedAction(cp.Plan, "%s %d nodegroups in cluster %q", verb, ngCount, cfg.Metadata.Name)

	cmdutils.LogPlanModeWarning(cp.Plan && ngCount > 0)

	return ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
		if cp.Plan {
			return nil
		}
		if err := drain.NodeGroup(clientSet, ng, ctl.Provider.WaitTimeout(), drainNodeGroupUndo); err != nil {
			return err
		}
		return nil
	})
}
