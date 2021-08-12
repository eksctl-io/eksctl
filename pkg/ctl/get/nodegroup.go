package get

import (
	"os"
	"strconv"
	"time"

	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
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
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doGetNodeGroup(cmd *cmdutils.Cmd, ng *api.NodeGroup, params *getCmdParams) error {
	cfg := cmd.ClusterConfig

	// TODO: move this into a loader when --config-file gets added to this command
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if ng.Name != "" && cmd.NameArg != "" {
		return cmdutils.ErrFlagAndArg("--name", ng.Name, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		ng.Name = cmd.NameArg
	}

	// prevent creation of invalid config object with unnamed nodegroup
	if ng.Name != "" {
		cfg.NodeGroups = append(cfg.NodeGroups, ng)
	}

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		cmdutils.LogRegionAndVersionInfo(cmd.ClusterConfig.Metadata)
	} else {
		//log warnings and errors to stderr
		logger.Writer = os.Stderr
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	var summaries []*manager.NodeGroupSummary
	if ng.Name == "" {
		summaries, err = nodegroup.New(cfg, ctl, clientSet).GetAll()
		if err != nil {
			return err
		}
	} else {
		summary, err := nodegroup.New(cfg, ctl, clientSet).Get(ng.Name)
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
	printer.AddColumn("CLUSTER", func(s *manager.NodeGroupSummary) string {
		return s.Cluster
	})
	printer.AddColumn("NODEGROUP", func(s *manager.NodeGroupSummary) string {
		return s.Name
	})
	printer.AddColumn("STATUS", func(s *manager.NodeGroupSummary) string {
		return s.Status
	})
	printer.AddColumn("CREATED", func(s *manager.NodeGroupSummary) string {
		return s.CreationTime.Format(time.RFC3339)
	})
	printer.AddColumn("MIN SIZE", func(s *manager.NodeGroupSummary) string {
		return strconv.Itoa(s.MinSize)
	})
	printer.AddColumn("MAX SIZE", func(s *manager.NodeGroupSummary) string {
		return strconv.Itoa(s.MaxSize)
	})
	printer.AddColumn("DESIRED CAPACITY", func(s *manager.NodeGroupSummary) string {
		return strconv.Itoa(s.DesiredCapacity)
	})
	printer.AddColumn("INSTANCE TYPE", func(s *manager.NodeGroupSummary) string {
		return s.InstanceType
	})
	printer.AddColumn("IMAGE ID", func(s *manager.NodeGroupSummary) string {
		return s.ImageID
	})
	printer.AddColumn("ASG NAME", func(s *manager.NodeGroupSummary) string {
		return s.AutoScalingGroupName
	})
}
