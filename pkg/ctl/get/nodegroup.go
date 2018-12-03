package get

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getNodegroupCmd() *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:     "nodegroup",
		Short:   "Get nodegroups(s)",
		Aliases: []string{"nodegroups"},
		Run: func(_ *cobra.Command, args []string) {
			if err := doGetNodegroups(p, cfg, cmdutils.GetNameArg(args)); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	fs := cmd.Flags()

	fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")

	fs.StringVarP(&p.Region, "region", "r", "", "AWS region")
	fs.StringVarP(&p.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

	fs.StringVarP(&output, "output", "o", "table", "Specifies the output format. Choose from table,json,yaml. Defaults to table.")
	return cmd
}

func doGetNodegroups(p *api.ProviderConfig, cfg *api.ClusterConfig, name string) error {
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

	manager := ctl.NewStackManager(cfg)
	summaries, err := manager.GetNodeGroupSummaries()
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

	if err := printer.PrintObj("nodegroups", summaries, os.Stdout); err != nil {
		return err
	}

	return nil
}

func addSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("STACKNAME", func(s *manager.NodeGroupSummary) string {
		return s.StackName
	})
	printer.AddColumn("SEQ", func(s *manager.NodeGroupSummary) string {
		return strconv.Itoa(s.Seq)
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
