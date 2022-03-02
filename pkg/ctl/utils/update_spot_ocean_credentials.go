package utils

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spotinst/spotinst-sdk-go/spotinst/featureflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/spot"
)

func updateSpotOceanCredentials(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	ng := spot.NewOceanVirtualNodeGroup()
	cmd.ClusterConfig = cfg

	var onlyMissing bool

	cmd.SetDescription("update-spot-ocean-credentials", "Update Spot Ocean Credentials", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doUpdateSpotOceanCredentials(cmd, ng, onlyMissing)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		cmdutils.AddNodeGroupFilterFlags(fs, &cmd.Include, &cmd.Exclude)
		fs.BoolVar(&onlyMissing, "only-missing", false, "only update nodegroups that are not defined in the given config file")
		fs.StringVarP(&ng.Name, "name", "n", "", "name of the nodegroup to update")
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func doUpdateSpotOceanCredentials(cmd *cmdutils.Cmd, ng *api.NodeGroup, onlyMissing bool) error {
	ngFilter := filter.NewNodeGroupFilter()
	if err := cmdutils.NewUtilsSpotOceanUpdateCredentials(cmd, ng, ngFilter).Load(); err != nil {
		return err
	}

	ctx := context.TODO()
	cfg := cmd.ClusterConfig
	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)
	stacks, err := stackManager.ListNodeGroupStacks(ctx)
	if err != nil {
		return err
	}

	if cmd.ClusterConfigFile != "" {
		logger.Info("comparing %d nodegroups defined in the given config (%q) "+
			"against remote state", len(cfg.NodeGroups), cmd.ClusterConfigFile)
		if onlyMissing {
			if err = ngFilter.SetOnlyRemote(ctx, ctl.AWSProvider.EKS(), stackManager, cfg); err != nil {
				return err
			}
		}
		for _, ng := range cfg.NodeGroups {
			if ng.SpotOcean != nil {
				cfg.NodeGroups = append(cfg.NodeGroups,
					spot.NewOceanClusterNodeGroup(cfg))
				break
			}
		}
	} else {
		cfg.NodeGroups = []*api.NodeGroup{ng}
	}

	logFiltered := cmdutils.ApplyFilter(cfg, ngFilter)
	logFiltered()

	featureflag.Set(fmt.Sprintf("%s=%t", spot.AllowCredentialsChanges.Name(), true))
	for _, ng := range cfg.NodeGroups {
		if spot.IsNodeGroupManagedByOcean(ng, stacks) && !cmd.Plan {
			if err := spot.UpdateCredentials(ctx, ctl.AWSProvider, ng, stacks); err != nil {
				return err
			}
		} else {
			logger.Debug("ocean: skipping credentials update for "+
				"nodegroup %q (reason: nodegroup isn't managed by ocean)", ng.Name)
		}
	}

	cmdutils.LogPlanModeWarning(cmd.Plan)
	return nil
}
