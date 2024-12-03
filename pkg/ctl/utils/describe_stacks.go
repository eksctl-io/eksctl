package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func describeStacksCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var all, events, trail bool
	var resourceStatus []string
	var output printers.Type

	cmd.SetDescription("describe-stacks", "Describe CloudFormation stack for a given cluster", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		var printer printers.OutputPrinter
		var err error
		if output != "" {
			if cmd.CobraCommand.Flags().Changed("all") ||
				cmd.CobraCommand.Flags().Changed("events") ||
				cmd.CobraCommand.Flags().Changed("trails") ||
				cmd.CobraCommand.Flags().Changed("resource-status") {
				return errors.New("since the output flag is specified, the flags `all`, `events`, `trail` and `resource-status` cannot be used")
			}
		}
		switch output {
		case printers.TableType:
			return fmt.Errorf("output type %q is not supported", output)
		case "":
		default:
			printer, err = printers.NewPrinter(output)
			if err != nil {
				return err
			}
		}
		return doDescribeStacksCmd(cmd, all, events, trail, resourceStatus, printer)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		fs.BoolVar(&all, "all", false, "include deleted stacks")
		fs.BoolVar(&events, "events", false, "include stack events")
		fs.StringSliceVar(&resourceStatus, "resource-status", nil, "resource statuses to filter events by, e.g. `CREATE_FAILED`, `UPDATE_FAILED`")
		fs.BoolVar(&trail, "trail", false, "lookup CloudTrail events for the cluster")
		fs.StringVarP(&output, "output", "o", "", "specifies the output formats (valid option: json and yaml)")
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func doDescribeStacksCmd(cmd *cmdutils.Cmd, all, events, trail bool, resourceStatus []string, printer printers.OutputPrinter) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	if printer != nil {
		logger.Writer = os.Stderr
	}

	ctx := context.TODO()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name != "" && cmd.NameArg != "" {
		return cmdutils.ErrFlagAndArg(cmdutils.ClusterNameFlag(cmd), cfg.Metadata.Name, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		cfg.Metadata.Name = cmd.NameArg
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if len(resourceStatus) > 0 && !events {
		return fmt.Errorf("--resource-status flag cannot be specified without setting --events=true")
	}

	stackManager := ctl.NewStackManager(cfg)
	stacks, err := stackManager.ListStacks(ctx)
	if err != nil {
		return err
	}

	if len(stacks) < 2 {
		logger.Warning("only %d stacks found, for a ready-to-use cluster there should be at least 2", len(stacks))
	}

	if printer != nil {
		return printer.PrintObj(stacks, cmd.CobraCommand.OutOrStdout())
	}

	for _, s := range stacks {
		if !all && s.StackStatus == types.StackStatusDeleteComplete {
			continue
		}
		logger.Info("stack/%s = %#v", *s.StackName, s)
		if events {
			events, err := stackManager.DescribeStackEvents(ctx, s)
			if err != nil {
				logger.Critical(err.Error())
			}

			if len(resourceStatus) > 0 {
				for i, rs := range resourceStatus {
					resourceStatus[i] = strings.ToLower(rs)
				}

				filteredEvents := []types.StackEvent{}
				for _, e := range events {
					if slices.Contains(resourceStatus, strings.ToLower(string(e.ResourceStatus))) {
						filteredEvents = append(filteredEvents, e)
					}
				}
				events = filteredEvents
			}

			for i, e := range events {
				logger.Info("CloudFormation.events/%s[%d] = %s", *s.StackName, i, StackEventToString(&e))
			}
		}
		if trail {
			events, err := stackManager.LookupCloudTrailEvents(ctx, s)
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

func StackEventToString(event *types.StackEvent) string {
	internalEvent := struct {
		TimeStamp            time.Time
		LogicalID            string
		ResourceStatus       string
		ResourceStatusReason string
	}{
		TimeStamp:      *event.Timestamp,
		LogicalID:      *event.LogicalResourceId,
		ResourceStatus: string(event.ResourceStatus),
	}
	if event.ResourceStatusReason != nil {
		internalEvent.ResourceStatusReason = *event.ResourceStatusReason
	}
	return fmt.Sprintf("%+v", internalEvent)
}
