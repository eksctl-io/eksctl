package get

import (
	"context"
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	capabilityactions "github.com/weaveworks/eksctl/pkg/actions/capability"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getCapabilityCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	params := &getCmdParams{}

	cmd.SetDescription(
		"capability",
		"Get capability(ies)",
		"",
		"capabilities",
	)

	var capabilityName string
	cmd.FlagSetGroup.InFlagSet("Capability", func(fs *pflag.FlagSet) {
		fs.StringVar(&capabilityName, "name", "", "name of the capability")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
	})
	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doGetCapability(cmd, capabilityName, params)
	}
}

func doGetCapability(cmd *cmdutils.Cmd, capabilityName string, params *getCmdParams) error {
	if err := cmdutils.NewGetCapabilityLoader(cmd).Load(); err != nil {
		return err
	}

	if params.output != printers.TableType {
		logger.Writer = os.Stderr
	}

	ctx := context.Background()
	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	capabilityGetter := capabilityactions.NewGetter(cmd.ClusterConfig.Metadata.Name, clusterProvider.AWSProvider.EKS())

	summaries, err := capabilityGetter.Get(ctx, capabilityName)
	if err != nil {
		return fmt.Errorf("failed to retrieve capabilities for cluster %s: %w", cmd.ClusterConfig.Metadata.Name, err)
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		addCapabilitySummaryTableColumns(printer.(*printers.TablePrinter))
	}

	return printer.PrintObjWithKind("capabilities", summaries, os.Stdout)
}

func addCapabilitySummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c capabilityactions.Summary) string {
		return c.Name
	})
	printer.AddColumn("TYPE", func(c capabilityactions.Summary) string {
		return c.Type
	})
	printer.AddColumn("STATUS", func(c capabilityactions.Summary) string {
		return c.Status
	})
	printer.AddColumn("VERSION", func(c capabilityactions.Summary) string {
		return c.Version
	})
}
