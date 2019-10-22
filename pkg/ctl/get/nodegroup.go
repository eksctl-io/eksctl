package get

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
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

	cmd.SetRunFuncWithNameArg(func() error {
		return doGetNodeGroup(cmd, ng, params)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
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

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	manager := ctl.NewStackManager(cfg)
	summaries, err := manager.GetNodeGroupSummaries(ng.Name)
	if err != nil {
		return errors.Wrap(err, "getting nodegroup stack summaries")
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == "table" {
		addSummaryTableColumns(printer.(*printers.TablePrinter))
	}

	if err := printer.PrintObjWithKind("nodegroups", summaries, os.Stdout); err != nil {
		return err
	}

	return nil
}

func addSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("CLUSTER", func(s *manager.NodeGroupSummary) string {
		return s.Cluster
	})
	printer.AddColumn("NODEGROUP", func(s *manager.NodeGroupSummary) string {
		return s.Name
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
}
