package utils

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func describeStacksCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var all, events, trail bool

	cmd.SetDescription("describe-stacks", "Describe CloudFormation stack for a given cluster", "")

	cmd.SetRunFuncWithNameArg(func() error {
		return doDescribeStacksCmd(cmd, all, events, trail)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		fs.BoolVar(&all, "all", false, "include deleted stacks")
		fs.BoolVar(&events, "events", false, "include stack events")
		fs.BoolVar(&trail, "trail", false, "lookup CloudTrail events for the cluster")
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doDescribeStacksCmd(cmd *cmdutils.Cmd, all, events, trail bool) error {
	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(cfg.Metadata)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && cmd.NameArg != "" {
		return cmdutils.ErrFlagAndArg(cmdutils.ClusterNameFlag(cmd), cfg.Metadata.Name, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		cfg.Metadata.Name = cmd.NameArg
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	stackManager := ctl.NewStackManager(cfg)

	stacks, err := stackManager.DescribeStacks()
	if err != nil {
		return err
	}

	if len(stacks) < 2 {
		logger.Warning("only %d stacks found, for a ready-to-use cluster there should be at least 2", len(stacks))
	}

	for _, s := range stacks {
		if !all && *s.StackStatus == cloudformation.StackStatusDeleteComplete {
			continue
		}
		logger.Info("stack/%s = %#v", *s.StackName, s)
		if events {
			events, err := stackManager.DescribeStackEvents(s)
			if err != nil {
				logger.Critical(err.Error())
			}
			for i, e := range events {
				logger.Info("CloudFormation.events/%s[%d] = %#v", *s.StackName, i, e)
			}
		}
		if trail {
			events, err := stackManager.LookupCloudTrailEvents(s)
			if err != nil {
				logger.Critical(err.Error())
			}
			for i, e := range events {
				logger.Info("CloudTrail.events/%s[%d] = %#v", *s.StackName, i, e)
			}
		}
	}

	return nil
}
