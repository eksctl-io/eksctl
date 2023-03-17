package get

import (
	"context"
	"os"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getIAMServiceAccountCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var name, namespace string
	cfg.IAM.WithOIDC = api.Enabled()

	params := &getCmdParams{}

	cmd.SetDescription("iamserviceaccount", "Get iamserviceaccount(s)", "", "iamserviceaccounts")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doGetIAMServiceAccount(cmd, IAMServiceAccountOptions{
			GetOptions:   irsa.GetOptions{Name: name, Namespace: namespace},
			getCmdParams: params,
		})
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		fs.StringVar(&namespace, "namespace", "", "namespace to look for iamserviceaccount")
		fs.StringVar(&name, "name", "", "name of iamserviceaccount to get")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

// IAMServiceAccountOptions holds the configuration for the get
// iamserviceaccounts action
type IAMServiceAccountOptions struct {
	irsa.GetOptions
	*getCmdParams
}

func doGetIAMServiceAccount(cmd *cmdutils.Cmd, options IAMServiceAccountOptions) error {
	if err := cmdutils.NewGetIAMServiceAccountLoader(cmd, &options.GetOptions).Load(); err != nil {
		return err
	}

	if options.output != printers.TableType {
		logger.Writer = os.Stderr
	}

	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)
	irsaManager := irsa.New(cfg.Metadata.Name, stackManager, nil, nil)
	serviceAccounts, err := irsaManager.Get(ctx, options.GetOptions)

	if err != nil {
		return err
	}

	printer, err := printers.NewPrinter(options.output)
	if err != nil {
		return err
	}

	if options.output == printers.TableType {
		addIAMServiceAccountSummaryTableColumns(printer.(*printers.TablePrinter))
	}

	return printer.PrintObjWithKind("iamserviceaccounts", serviceAccounts, os.Stdout)
}

func addIAMServiceAccountSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAMESPACE", func(sa *api.ClusterIAMServiceAccount) string {
		return sa.Namespace
	})
	printer.AddColumn("NAME", func(sa *api.ClusterIAMServiceAccount) string {
		return sa.Name
	})
	printer.AddColumn("ROLE ARN", func(sa *api.ClusterIAMServiceAccount) string {
		return *sa.Status.RoleARN
	})
}
