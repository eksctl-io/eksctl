package delete

import (
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/drain"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func deleteNodeGroupCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	rc.ClusterConfig = cfg

	var updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool

	rc.SetDescription("nodegroup", "Delete a nodegroup", "", "ng")

	rc.SetRunFuncWithNameArg(func() error {
		return doDeleteNodeGroup(rc, ng, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup to delete")
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, rc)
		cmdutils.AddNodeGroupFilterFlags(fs, &rc.IncludeNodeGroups, &rc.ExcludeNodeGroups)
		fs.BoolVar(&onlyMissing, "only-missing", false, "Only delete nodegroups that are not defined in the given config file")
		cmdutils.AddUpdateAuthConfigMap(fs, &updateAuthConfigMap, "Remove nodegroup IAM role from aws-auth configmap")
		fs.BoolVar(&deleteNodeGroupDrain, "drain", true, "Drain and cordon all nodes in the nodegroup before deletion")

		rc.Wait = false
		cmdutils.AddWaitFlag(fs, &rc.Wait, "deletion of all resources")
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, true, true)
}

func doDeleteNodeGroup(rc *cmdutils.ResourceCmd, ng *api.NodeGroup, updateAuthConfigMap, deleteNodeGroupDrain, onlyMissing bool) error {
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

	if updateAuthConfigMap {
		cmdutils.LogIntendedAction(rc.Plan, "delete %d nodegroups from auth ConfigMap in cluster %q", ngCount, cfg.Metadata.Name)
		err := ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
			if rc.Plan {
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
		cmdutils.LogIntendedAction(rc.Plan, "drain %d nodegroups in cluster %q", ngCount, cfg.Metadata.Name)
		err := ngFilter.ForEach(cfg.NodeGroups, func(_ int, ng *api.NodeGroup) error {
			if rc.Plan {
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

	cmdutils.LogIntendedAction(rc.Plan, "delete %d nodegroups from cluster %q", ngCount, cfg.Metadata.Name)

	{
		tasks, err := stackManager.NewTasksToDeleteNodeGroups(ngSubset, rc.Wait, nil)
		if err != nil {
			return err
		}
		tasks.PlanMode = rc.Plan
		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			return handleErrors(errs, "nodegroup(s)")
		}
		cmdutils.LogCompletedAction(rc.Plan, "deleted %d nodegroups from cluster %q", ngCount, cfg.Metadata.Name)
	}

	cmdutils.LogPlanModeWarning(rc.Plan && ngCount > 0)

	return nil
}
