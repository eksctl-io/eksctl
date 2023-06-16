package get

import (
	"context"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/weaveworks/eksctl/pkg/actions/label"
	"github.com/weaveworks/eksctl/pkg/printers"

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
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		fs.StringVarP(&nodeGroupName, "nodegroup", "n", "", "Nodegroup name")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

}

func getLabels(cmd *cmdutils.Cmd, nodeGroupName string) error {
	if err := cmdutils.NewGetLabelsLoader(cmd, nodeGroupName).Load(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig

	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	service := managed.NewService(ctl.AWSProvider.EKS(), ctl.AWSProvider.EC2(), manager.NewStackCollection(ctl.AWSProvider, cfg), cfg.Metadata.Name)
	manager := label.New(cfg.Metadata.Name, service, ctl.AWSProvider.EKS())
	labels, err := manager.Get(ctx, nodeGroupName)
	if err != nil {
		return err
	}

	printer := printers.NewTablePrinter()
	addColumns(printer.(*printers.TablePrinter))
	return printer.PrintObjWithKind("labels", labels, cmd.CobraCommand.OutOrStdout())
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
