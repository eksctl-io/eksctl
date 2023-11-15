package get

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getPodIdentityAssociationCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	params := &getCmdParams{}

	cmd.SetDescription(
		"podidentityassociation",
		"Get a pod identity association",
		"",
		"podidentityassociations",
	)

	pia := &api.PodIdentityAssociation{}
	configureGetPodIdentityAssociationCmd(cmd, params, pia)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		if err := cmdutils.NewGetPodIdentityAssociationLoader(cmd, pia).Load(); err != nil {
			return err
		}
		return doGetPodIdentityAssociation(cmd, pia.Namespace, pia.ServiceAccountName, params)
	}
}

func doGetPodIdentityAssociation(cmd *cmdutils.Cmd, namespace, serviceAccountName string, params *getCmdParams) error {
	cfg := cmd.ClusterConfig
	ctx := context.Background()

	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	summaries, err := podidentityassociation.NewGetter(cfg.Metadata.Name, ctl.AWSProvider.EKS()).
		GetPodIdentityAssociations(ctx, namespace, serviceAccountName)
	if err != nil {
		return err
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		addPodIdentityAssociationSummaryTableColumns(printer.(*printers.TablePrinter))
	}

	return printer.PrintObjWithKind("podidentityassociations", summaries, cmd.CobraCommand.OutOrStdout())
}

func configureGetPodIdentityAssociationCmd(cmd *cmdutils.Cmd, params *getCmdParams, pia *api.PodIdentityAssociation) {
	cmd.FlagSetGroup.InFlagSet("PodIdentityAssociation", func(fs *pflag.FlagSet) {
		fs.StringVar(&pia.Namespace, "namespace", "", "Namespace the service account belongs to")
		fs.StringVar(&pia.ServiceAccountName, "service-account-name", "", "Name of the service account")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
}

func addPodIdentityAssociationSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("ASSOCIATION ARN", func(s podidentityassociation.Summary) string {
		return s.AssociationARN
	})
	printer.AddColumn("NAMESPACE", func(s podidentityassociation.Summary) string {
		return s.Namespace
	})
	printer.AddColumn("SERVICE ACCOUNT NAME", func(s podidentityassociation.Summary) string {
		return s.ServiceAccountName
	})
	printer.AddColumn("IAM ROLE ARN", func(s podidentityassociation.Summary) string {
		return s.RoleARN
	})
}
