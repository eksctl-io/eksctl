package get

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/amazon-ec2-instance-selector/v2/pkg/selector"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getNodeGroupCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	ng := api.NewNodeGroup()
	cmd.ClusterConfig = cfg

	params := &getCmdParams{}

	cmd.SetDescription("nodegroup", "Get nodegroup(s)", "", "ng", "nodegroups")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doGetNodeGroup(cmd, ng, params)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func doGetNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup, params *getCmdParams) error {
	if err := cmdutils.NewGetNodegroupLoader(cmd, ng).Load(); err != nil {
		return err
	}

	if params.output != printers.TableType {
		//log warnings and errors to stderr
		logger.Writer = os.Stderr
	}

	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	var summaries []*nodegroup.Summary
	manager := nodegroup.New(cfg, ctl, clientSet, selector.New(ctl.AWSProvider.Session()))
	if ng.Name == "" {
		summaries, err = manager.GetAll(ctx)
		if err != nil {
			return err
		}
	} else {
		summary, err := manager.Get(ctx, ng.Name)
		if err != nil {
			return err
		}
		summaries = append(summaries, summary)
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		// Empty summary implies no nodegroups
		// We only error if the output is table, since if the output
		// is yaml or json we should return an empty object.
		if len(summaries) == 0 {
			if ng.Name == "" {
				return errors.Errorf("No nodegroups found")
			}
			return errors.Errorf("nodegroup with name %v not found", ng.Name)
		}
		addSummaryTableColumns(printer.(*printers.TablePrinter))
	}

	return printer.PrintObjWithKind("nodegroups", summaries, os.Stdout)
}

func addSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("CLUSTER", func(s *nodegroup.Summary) string {
		return s.Cluster
	})
	printer.AddColumn("NODEGROUP", func(s *nodegroup.Summary) string {
		return s.Name
	})
	printer.AddColumn("STATUS", func(s *nodegroup.Summary) string {
		return s.Status
	})
	printer.AddColumn("CREATED", func(s *nodegroup.Summary) string {
		return s.CreationTime.Format(time.RFC3339)
	})
	printer.AddColumn("MIN SIZE", func(s *nodegroup.Summary) string {
		return strconv.Itoa(s.MinSize)
	})
	printer.AddColumn("MAX SIZE", func(s *nodegroup.Summary) string {
		return strconv.Itoa(s.MaxSize)
	})
	printer.AddColumn("DESIRED CAPACITY", func(s *nodegroup.Summary) string {
		return strconv.Itoa(s.DesiredCapacity)
	})
	printer.AddColumn("INSTANCE TYPE", func(s *nodegroup.Summary) string {
		return s.InstanceType
	})
	printer.AddColumn("IMAGE ID", func(s *nodegroup.Summary) string {
		return s.ImageID
	})
	printer.AddColumn("ASG NAME", func(s *nodegroup.Summary) string {
		return s.AutoScalingGroupName
	})
	printer.AddColumn("TYPE", func(s *nodegroup.Summary) api.NodeGroupType {
		return s.NodeGroupType
	})
}
