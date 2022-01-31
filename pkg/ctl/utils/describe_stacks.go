package utils

import (
	"os"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func describeStacksCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var all, events, trail bool
	var output printers.Type

	cmd.SetDescription("describe-stacks", "Describe CloudFormation stack for a given cluster", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		var printer printers.OutputPrinter
		var err error
		if output != "" {
			if cmd.CobraCommand.Flags().Changed("all") ||
				cmd.CobraCommand.Flags().Changed("events") ||
				cmd.CobraCommand.Flags().Changed("trails") {
				return errors.Errorf("since the output flag is specified, the flags `all`, `events` and `trail` cannot be used")
			}
		}
		switch output {
		case printers.TableType:
			return errors.Errorf("output type %q is not supported", output)
		case "":
		default:
			printer, err = printers.NewPrinter(output)
			if err != nil {
				return err
			}
		}
		return doDescribeStacksCmd(cmd, all, events, trail, printer)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		fs.BoolVar(&all, "all", false, "include deleted stacks")
		fs.BoolVar(&events, "events", false, "include stack events")
		fs.BoolVar(&trail, "trail", false, "lookup CloudTrail events for the cluster")
		fs.StringVarP(&output, "output", "o", "", "specifies the output formats (valid option: json and yaml)")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doDescribeStacksCmd(cmd *cmdutils.Cmd, all, events, trail bool, printer printers.OutputPrinter) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}
	if printer == nil {
		cmdutils.LogRegionAndVersionInfo(cfg.Metadata)
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

	if printer != nil {
		return printer.PrintObj(stacks, os.Stdout)
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
