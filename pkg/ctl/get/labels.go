package get

import (
	"os"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/label"
	"github.com/weaveworks/eksctl/pkg/printers"
	"k8s.io/apimachinery/pkg/labels"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func getLabelsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("labels", "Get labels for managed nodegroup", "")

	var nodeGroupName string
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return getLabels(cmd, nodeGroupName)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&nodeGroupName, "nodegroup", "n", "", "Nodegroup name")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

}

func getLabels(cmd *cmdutils.Cmd, nodeGroupName string) error {
	if err := cmdutils.NewGetLabelsLoader(cmd, nodeGroupName).Load(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}
	cmdutils.LogRegionAndVersionInfo(cmd.ClusterConfig.Metadata)

	service := managed.NewService(ctl.Provider.EKS(), ctl.Provider.SSM(), ctl.Provider.EC2(), manager.NewStackCollection(ctl.Provider, cfg), cfg.Metadata.Name)
	manager := label.New(cfg.Metadata.Name, service, ctl.Provider.EKS())
	labels, err := manager.Get(nodeGroupName)
	if err != nil {
		return err
	}

	printer := printers.NewTablePrinter()
	addColumns(printer.(*printers.TablePrinter))
	return printer.PrintObjWithKind("labels", labels, os.Stdout)
}

func addColumns(printer *printers.TablePrinter) {
	printer.AddColumn("CLUSTER", func(s label.Summary) string {
		return s.Cluster
	})
	printer.AddColumn("NODEGROUP", func(s label.Summary) string {
		return s.NodeGroup
	})
	printer.AddColumn("LABELS", func(s label.Summary) string {
		return labels.FormatLabels(s.Labels)
	})
}
