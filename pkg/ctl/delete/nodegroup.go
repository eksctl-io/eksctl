package delete

import (
	"os"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/eks"
)

var (
	updateAuthConfigMap  bool
	deleteNodeGroupDrain bool

	includeNodeGroups []string
	excludeNodeGroups []string

	deleteOnlyMissingNodeGroups bool
)

func deleteNodeGroupCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:     "nodegroup",
		Short:   "Delete a nodegroup",
		Aliases: []string{"ng"},
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doDeleteNodeGroup(cp, ng); err != nil {
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
		fs.BoolVar(&deleteOnlyMissingNodeGroups, "only-missing", false, "Only delete nodegroups that are not defined in the given config file")
		cmdutils.AddUpdateAuthConfigMap(fs, &updateAuthConfigMap, "Remove nodegroup IAM role from aws-auth configmap")
		fs.BoolVar(&deleteNodeGroupDrain, "drain", true, "Drain and cordon all nodes in the nodegroup before deletion")
		cmdutils.AddWaitFlag(fs, &cp.Wait, "deletion of all resources")
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, true)

	group.AddTo(cp.Command)
	return cp.Command
}

func doDeleteNodeGroup(cp *cmdutils.CommonParams, ng *api.NodeGroup) error {
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
		if err := ngFilter.SetIncludeOrExcludeMissingFilter(stackManager, deleteOnlyMissingNodeGroups, &cfg.NodeGroups); err != nil {
			return err
		}
	}

	ngSubset, _ := ngFilter.MatchAll(cfg.NodeGroups)
	ngCount := ngSubset.Len()

	ngFilter.LogInfo(cfg.NodeGroups)

	if updateAuthConfigMap {
		cmdutils.LogIntendedAction(cp.Plan, "delete %d nodegroups from auth ConfigMap in cluster %q", ngCount, cfg.Metadata.Name)
		err := ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
			if cp.Plan {
				return nil
			}
			if ng.IAM == nil || ng.IAM.InstanceRoleARN == "" {
				if err := ctl.GetNodeGroupIAM(stackManager, cfg, ng); err != nil {
					logger.Warning("error getting instance role ARN for nodegroup %q", ng.Name)
					return nil
				}
			}
			if err := authconfigmap.RemoveNodeGroup(clientSet, ng); err != nil {
				logger.Warning(err.Error())
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if deleteNodeGroupDrain {
		cmdutils.LogIntendedAction(cp.Plan, "drain %d nodegroups in cluster %q", ngCount, cfg.Metadata.Name)
		err := ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
			if cp.Plan {
				return nil
			}
			if err := drain.NodeGroup(clientSet, ng, ctl.Provider.WaitTimeout(), false); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	cmdutils.LogIntendedAction(cp.Plan, "delete %d nodegroups from cluster %q", ngCount, cfg.Metadata.Name)

	{
		tasks, err := stackManager.NewTasksToDeleteNodeGroups(ngSubset, cp.Wait, nil)
		if err != nil {
			return err
		}
		tasks.PlanMode = cp.Plan
		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			return handleErrors(errs, "nodegroup(s)")
		}
		cmdutils.LogCompletedAction(cp.Plan, "deleted %d nodegroups from cluster %q", ngCount, cfg.Metadata.Name)
	}

	cmdutils.LogPlanModeWarning(cp.Plan && ngCount > 0)

	return nil
}
