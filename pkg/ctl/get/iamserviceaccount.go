package get

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/iam"

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
		return doGetIAMServiceAccount(cmd, namespace, name, params)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")

		fs.StringVar(&namespace, "namespace", "", "namespace to look for iamserviceaccount")
		fs.StringVar(&name, "name", "", "name of iamserviceaccount to get")

		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doGetIAMServiceAccount(cmd *cmdutils.Cmd, namespace, name string, params *getCmdParams) error {
	if err := cmdutils.NewGetIAMServiceAccountLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if params.output == "table" {
		cmdutils.LogRegionAndVersionInfo(cmd.ClusterConfig.Metadata)
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)
	iamServiceAccountManager := iam.New(cfg.Metadata.Name, ctl, stackManager, nil, nil)
	serviceAccounts, err := iamServiceAccountManager.Get(namespace, name)

	if err != nil {
		return err
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	var obj interface{}
	if params.output == "table" {
		addIAMServiceAccountSummaryTableColumns(printer.(*printers.TablePrinter))
		obj = serviceAccounts
	} else {
		cfg.IAM.ServiceAccounts = serviceAccounts
		obj = cfg
	}
	return printer.PrintObjWithKind("iamserviceaccounts", obj, os.Stdout)
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
