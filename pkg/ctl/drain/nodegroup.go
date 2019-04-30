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
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:     "nodegroup",
		Short:   "Cordon and drain a nodegroup",
		Aliases: []string{"ng"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := doDrainNodeGroup(p, cfg, ng, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, p)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)
		cmdutils.AddApproveFlag(&plan, cmd, fs)
		cmdutils.AddNodeGroupFilterFlags(&includeNodeGroups, &excludeNodeGroups, fs)
		fs.BoolVar(&drainOnlyMissingNodeGroups, "only-missing", false, "Only drain nodegroups that are not defined in the given config file")
		fs.BoolVar(&drainNodeGroupUndo, "undo", false, "Uncordone the nodegroup")
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.AddTo(cmd)

	return cmd
}

func doDrainNodeGroup(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup, nameArg string, cmd *cobra.Command) error {
	ngFilter := cmdutils.NewNodeGroupFilter()

	if err := cmdutils.NewDeleteNodeGroupLoader(p, cfg, ng, clusterConfigFile, nameArg, cmd, ngFilter, includeNodeGroups, excludeNodeGroups, &plan).Load(); err != nil {
		return err
	}

	ctl := eks.New(p, cfg)

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

	if clusterConfigFile != "" {
		logger.Info("comparing %d nodegroups defined in the given config (%q) against remote state", len(cfg.NodeGroups), clusterConfigFile)
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
	cmdutils.LogIntendedAction(plan, "%s %d nodegroups in cluster %q", verb, ngCount, cfg.Metadata.Name)

	cmdutils.LogPlanModeWarning(plan && ngCount > 0)

	return ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
		if plan {
			return nil
		}
		if err := drain.NodeGroup(clientSet, ng, ctl.Provider.WaitTimeout(), drainNodeGroupUndo); err != nil {
			return err
		}
		return nil
	})
}
