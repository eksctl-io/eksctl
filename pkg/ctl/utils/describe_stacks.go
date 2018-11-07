package utils

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

var (
	utilsDescribeStackAll    bool
	utilsDescribeStackEvents bool
)

func describeStacksCmd() *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "describe-stacks",
		Short: "Describe CloudFormation stack for a given cluster",
		Run: func(_ *cobra.Command, args []string) {
			if err := doDescribeStacksCmd(p, cfg, ctl.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name (required)")

	fs.StringVarP(&p.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&p.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.BoolVar(&utilsDescribeStackAll, "all", false, "include deleted stacks")
	fs.BoolVar(&utilsDescribeStackEvents, "events", false, "include stack events")

	return cmd
}

func doDescribeStacksCmd(p *api.ProviderConfig, cfg *api.ClusterConfig, name string) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name != "" && name != "" {
		return fmt.Errorf("--name=%s and argument %s cannot be used at the same time", cfg.Metadata.Name, name)
	}

	if name != "" {
		cfg.Metadata.Name = name
	}

	if cfg.Metadata.Name == "" {
		return fmt.Errorf("--name must be set")
	}

	stackManager := ctl.NewStackManager(cfg)

	stacks, err := stackManager.DescribeStacks(cfg.Metadata.Name)
	if err != nil {
		return err
	}

	if len(stacks) < 2 {
		logger.Warning("only %d stacks found, for a ready-to-use cluster there should be at least 2", len(stacks))
	}

	for _, s := range stacks {
		if !utilsDescribeStackAll && *s.StackStatus == cloudformation.StackStatusDeleteComplete {
			continue
		}
		logger.Info("stack/%s = %#v", *s.StackName, s)
		if utilsDescribeStackEvents {
			events, err := stackManager.DescribeStackEvents(s)
			if err != nil {
				logger.Critical(err.Error())
			}
			for i, e := range events {
				logger.Info("events/%s[%d] = %#v", *s.StackName, i, e)
			}
		}
	}

	return nil
}
