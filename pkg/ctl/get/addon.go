package get

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getAddonCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	params := &getCmdParams{}

	cmd.SetDescription(
		"addon",
		"Get an Addon",
		"",
		"addons",
	)

	var a api.Addon
	cmd.FlagSetGroup.InFlagSet("Addon", func(fs *pflag.FlagSet) {
		fs.StringVar(&a.Name, "name", "", "Addon name")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return getAddon(cmd, &a, params)
	}
}

func getAddon(cmd *cmdutils.Cmd, a *api.Addon, params *getCmdParams) error {
	if err := cmdutils.NewGetAddonsLoader(cmd).Load(); err != nil {
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

	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)

	output, err := clusterProvider.AWSProvider.EKS().DescribeCluster(ctx, &awseks.DescribeClusterInput{
		Name: &cmd.ClusterConfig.Metadata.Name,
	})

	if err != nil {
		return fmt.Errorf("failed to fetch cluster %q version: %v", cmd.ClusterConfig.Metadata.Name, err)
	}

	logger.Info("Kubernetes version %q in use by cluster %q", *output.Cluster.Version, cmd.ClusterConfig.Metadata.Name)
	cmd.ClusterConfig.Metadata.Version = *output.Cluster.Version

	addonManager, err := addon.New(cmd.ClusterConfig, clusterProvider.AWSProvider.EKS(), stackManager, *cmd.ClusterConfig.IAM.WithOIDC, nil, nil)

	if err != nil {
		return err
	}

	var summaries []addon.Summary
	if a.Name == "" {
		summaries, err = addonManager.GetAll(ctx)
		if err != nil {
			return err
		}
	} else {
		summary, err := addonManager.Get(ctx, a)
		if err != nil {
			return err
		}
		summaries = []addon.Summary{summary}
	}

	if len(summaries) > 0 {
		logger.Info("to see issues for an addon run `eksctl get addon --name <addon-name> --cluster <cluster-name>`")
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if tablePrinter, ok := printer.(*printers.TablePrinter); ok {
		if slices.ContainsFunc(summaries, func(summary addon.Summary) bool {
			return len(summary.PodIdentityAssociations) > 0
		}) {
			logger.Info("to view pod identity associations for an addon, rerun the command with --output=json or --output=yaml")
		}
		addAddonSummaryTableColumns(tablePrinter)
	}

	if err := printer.PrintObjWithKind("addons", summaries, cmd.CobraCommand.OutOrStdout()); err != nil {
		return err
	}

	// if getting a particular addon, print the issue.
	if a.Name != "" {
		for _, issue := range summaries[0].Issues {
			fmt.Printf("Issue: %+v\n", issue)
		}
	}

	return nil
}

func addAddonSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(s addon.Summary) string {
		return s.Name
	})
	printer.AddColumn("VERSION", func(s addon.Summary) string {
		return s.Version
	})
	printer.AddColumn("STATUS", func(s addon.Summary) string {
		return s.Status
	})
	printer.AddColumn("ISSUES", func(s addon.Summary) int {
		return len(s.Issues)
	})
	printer.AddColumn("IAMROLE", func(s addon.Summary) string {
		return s.IAMRole
	})
	printer.AddColumn("UPDATE AVAILABLE", func(s addon.Summary) string {
		return s.NewerVersion
	})
	printer.AddColumn("CONFIGURATION VALUES", func(s addon.Summary) string {
		return s.ConfigurationValues
	})
	printer.AddColumn("POD IDENTITY ASSOCIATION ROLES", func(s addon.Summary) string {
		var roleARNs []string
		for _, pia := range s.PodIdentityAssociations {
			roleARNs = append(roleARNs, pia.RoleARN)
		}
		return strings.Join(roleARNs, ",")
	})
}
