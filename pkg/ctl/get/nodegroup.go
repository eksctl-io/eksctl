package get

import (
	"os"
	"strconv"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getNodegroupCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	ng := cfg.NewNodeGroup()

	cmd := &cobra.Command{
		Use:     "nodegroup",
		Short:   "Get nodegroup(s)",
		Aliases: []string{"ng", "nodegroups"},
		Run: func(_ *cobra.Command, args []string) {
			if err := doGetNodegroups(p, cfg, ng, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&ng.Name, "name", "n", "", "Name of the nodegroup")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddCommonFlagsForGetCmd(fs, &chunkSize, &output)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)

	return cmd
}

func doGetNodegroups(p *api.ProviderConfig, cfg *api.ClusterConfig, ng *api.NodeGroup, nameArg string) error {
	ctl := eks.New(p, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--cluster")
	}

	if ng.Name != "" && nameArg != "" {
		return cmdutils.ErrNameFlagAndArg(ng.Name, nameArg)
	}

	if nameArg != "" {
		ng.Name = nameArg
	}

	manager := ctl.NewStackManager(cfg)
	summaries, err := manager.GetNodeGroupSummaries(ng.Name)
	if err != nil {
		return errors.Wrap(err, "getting nodegroup stack summaries")
	}

	printer, err := printers.NewPrinter(output)
	if err != nil {
		return err
	}

	if output == "table" {
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
