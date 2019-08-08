package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func describeStacksCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	var all, events, trail bool

	rc.SetDescription("describe-stacks", "Describe CloudFormation stack for a given cluster", "")

	rc.SetRunFuncWithNameArg(func() error {
		return doDescribeStacksCmd(rc, all, events, trail)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		fs.BoolVar(&all, "all", false, "include deleted stacks")
		fs.BoolVar(&events, "events", false, "include stack events")
		fs.BoolVar(&trail, "trail", false, "lookup CloudTrail events for the cluster")
		cmdutils.AddTimeoutFlag(fs, &rc.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doDescribeStacksCmd(rc *cmdutils.ResourceCmd, all, events, trail bool) error {
	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && rc.NameArg != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.Metadata.Name, rc.NameArg)
	}

	if rc.NameArg != "" {
		cfg.Metadata.Name = rc.NameArg
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
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
