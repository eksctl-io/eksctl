package drain

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"

	"github.com/weaveworks/eksctl/pkg/drain"
)

func drainNodeGroupCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	rc.ClusterConfig = cfg

	var undo, onlyMissing bool

	rc.SetDescription("nodegroup", "Cordon and drain a nodegroup", "", "ng")

	rc.SetRunFuncWithNameArg(func() error {
		return doDrainNodeGroup(rc, ng, undo, onlyMissing)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, rc)
		cmdutils.AddNodeGroupFilterFlags(fs, &rc.IncludeNodeGroups, &rc.ExcludeNodeGroups)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only drain nodegroups that are not defined in the given config file")
		fs.BoolVar(&undo, "undo", false, "Uncordone the nodegroup")
		cmdutils.AddTimeoutFlag(fs, &rc.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, true)
}

func doDrainNodeGroup(rc *cmdutils.ResourceCmd, ng *api.NodeGroup, undo, onlyMissing bool) error {
	ngFilter := cmdutils.NewNodeGroupFilter()

	if err := cmdutils.NewDeleteNodeGroupLoader(rc, ng, ngFilter).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if err := ctl.RefreshClusterConfig(cfg); err != nil {
		return errors.Wrapf(err, "getting credentials for cluster %q", cfg.Metadata.Name)
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	if rc.ClusterConfigFile != "" {
		logger.Info("comparing %d nodegroups defined in the given config (%q) against remote state", len(cfg.NodeGroups), rc.ClusterConfigFile)
		if err := ngFilter.SetIncludeOrExcludeMissingFilter(stackManager, onlyMissing, &cfg.NodeGroups); err != nil {
			return err
		}
	}

	ngSubset, _ := ngFilter.MatchAll(cfg.NodeGroups)
	ngCount := ngSubset.Len()

	ngFilter.LogInfo(cfg.NodeGroups)
	verb := "drain"
	if undo {
		verb = "uncordon"
	}
	cmdutils.LogIntendedAction(rc.Plan, "%s %d nodegroups in cluster %q", verb, ngCount, cfg.Metadata.Name)

	cmdutils.LogPlanModeWarning(rc.Plan && ngCount > 0)

	return ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
		if rc.Plan {
			return nil
		}
		if err := drain.NodeGroup(clientSet, ng, ctl.Provider.WaitTimeout(), undo); err != nil {
			return err
		}
		return nil
	})
}
