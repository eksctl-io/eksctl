package get

import (
	"context"
	"fmt"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/accessentry"
	accessentryactions "github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getAccessEntryCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	params := &getCmdParams{}

	cmd.SetDescription(
		"accessentry",
		"Get access entry(ies)",
		"",
		"accessentries",
	)

	var principalARN api.ARN
	cmd.FlagSetGroup.InFlagSet("AccessEntry", func(fs *pflag.FlagSet) {
		fs.VarP(&principalARN, "principal-arn", "", "principal ARN to which the access entry is associated")
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
		return doGetAccessEntry(cmd, principalARN, params)
	}
}

func doGetAccessEntry(cmd *cmdutils.Cmd, principalARN api.ARN, params *getCmdParams) error {
	if err := cmdutils.NewGetAccessEntryLoader(cmd).Load(); err != nil {
		return err
	}

	if params.output != printers.TableType {
		//log warnings and errors to stdout
		logger.Writer = os.Stderr
	}

	ctx := context.Background()
	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	accessEntry := &accessentry.Service{
		ClusterStateGetter: clusterProvider,
	}
	if !accessEntry.IsEnabled() {
		return accessentry.ErrDisabledAccessEntryAPI
	}

	accessEntryGetter := accessentryactions.NewGetter(cmd.ClusterConfig.Metadata.Name, clusterProvider.AWSProvider.EKS())

	summaries, err := accessEntryGetter.Get(ctx, principalARN)
	if err != nil {
		return fmt.Errorf("failed to retrieve access entries for cluster %s: %w", cmd.ClusterConfig.Metadata.Name, err)
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		addAccessEntrySummaryTableColumns(printer.(*printers.TablePrinter))
		logger.Info("to get a detailed view of Kubernetes groups or policies associated with each access entry, use --output yaml or json")
	}

	return printer.PrintObjWithKind("accessentries", summaries, os.Stdout)
}

func addAccessEntrySummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("PRINCIPAL ARN", func(s accessentryactions.Summary) string {
		return s.PrincipalARN
	})
	printer.AddColumn("KUBERNETES GROUPS", func(s accessentryactions.Summary) int {
		return len(s.KubernetesGroups)
	})
	printer.AddColumn("ACCESS POLICIES", func(s accessentryactions.Summary) int {
		return len(s.AccessPolicies)
	})
}
