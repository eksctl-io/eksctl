package get

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"
	"github.com/weaveworks/eksctl/pkg/printers"
	"k8s.io/apimachinery/pkg/labels"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func getLabelsCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("labels", "Get nodegroup labels", "")

	var nodeGroupName string
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return getLabels(cmd, nodeGroupName)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		fs.StringVarP(&nodeGroupName, "nodegroup", "n", "", "Nodegroup name")

		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)

}

type summary struct {
	Cluster   string
	NodeGroup string
	Labels    map[string]string
}

func getLabels(cmd *cmdutils.Cmd, nodeGroupName string) error {
	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if cmd.NameArg != "" {
		return cmdutils.ErrUnsupportedNameArg()
	}

	ctl := eks.New(cmd.ProviderConfig, cmd.ClusterConfig)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	stackCollection := manager.NewStackCollection(ctl.Provider, cfg)
	managedService := managed.NewService(ctl.Provider, stackCollection, cfg.Metadata.Name)
	ngLabels, err := managedService.GetLabels(nodeGroupName)
	if err != nil {
		return err
	}

	out := []summary{
		{
			Cluster:   cfg.Metadata.Name,
			NodeGroup: nodeGroupName,
			Labels:    ngLabels,
		},
	}

	printer := printers.NewTablePrinter()
	addColumns(printer.(*printers.TablePrinter))
	return printer.PrintObjWithKind("labels", out, os.Stdout)
}

func addColumns(printer *printers.TablePrinter) {
	printer.AddColumn("CLUSTER", func(s summary) string {
		return s.Cluster
	})
	printer.AddColumn("NODEGROUP", func(s summary) string {
		return s.NodeGroup
	})
	printer.AddColumn("LABELS", func(s summary) string {
		return labels.FormatLabels(s.Labels)
	})
}
