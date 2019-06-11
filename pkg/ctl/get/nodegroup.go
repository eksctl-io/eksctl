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
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getNodeGroupCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()
	rc.ClusterConfig = cfg

	params := &getCmdParams{}

	rc.SetDescription("nodegroup", "Get nodegroup(s)", "", "ng", "nodegroups")

	rc.SetRunFuncWithNameArg(func() error {
		return doGetNodeGroup(rc, ng, params)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup")
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doGetNodeGroup(rc *cmdutils.ResourceCmd, ng *api.NodeGroup, params *getCmdParams) error {
	cfg := rc.ClusterConfig
	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--cluster")
	}

	if ng.Name != "" && rc.NameArg != "" {
		return cmdutils.ErrNameFlagAndArg(ng.Name, rc.NameArg)
	}

	if rc.NameArg != "" {
		ng.Name = rc.NameArg
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
