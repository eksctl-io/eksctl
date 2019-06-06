package utils

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

var (
	describeStacksAll    bool
	describeStacksEvents bool
	describeStacksTrail  bool
)

func describeStacksCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	cp.Command = &cobra.Command{
		Use:   "describe-stacks",
		Short: "Describe CloudFormation stack for a given cluster",
		Run: func(_ *cobra.Command, args []string) {
			cp.NameArg = cmdutils.GetNameArg(args)
			if err := doDescribeStacksCmd(cp); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
		fs.BoolVar(&describeStacksAll, "all", false, "include deleted stacks")
		fs.BoolVar(&describeStacksEvents, "events", false, "include stack events")
		fs.BoolVar(&describeStacksTrail, "trail", false, "lookup CloudTrail events for the cluster")
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, false)

	group.AddTo(cp.Command)
	return cp.Command
}

func doDescribeStacksCmd(cp *cmdutils.CommonParams) error {
	cfg := cp.ClusterConfig

	ctl := eks.New(cp.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && cp.NameArg != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.Metadata.Name, cp.NameArg)
	}

	if cp.NameArg != "" {
		cfg.Metadata.Name = cp.NameArg
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
		if !describeStacksAll && *s.StackStatus == cloudformation.StackStatusDeleteComplete {
			continue
		}
		logger.Info("stack/%s = %#v", *s.StackName, s)
		if describeStacksEvents {
			events, err := stackManager.DescribeStackEvents(s)
			if err != nil {
				logger.Critical(err.Error())
			}
			for i, e := range events {
				logger.Info("CloudFormation.events/%s[%d] = %#v", *s.StackName, i, e)
			}
		}
		if describeStacksTrail {
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
