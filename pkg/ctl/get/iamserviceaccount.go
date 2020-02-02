package get

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getIAMServiceAccountCmd(cmd *cmdutils.Cmd) {
	getIAMServiceAccountWithRunFunc(cmd, func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, params *getCmdParams) error {
		return doGetIAMServiceAccount(cmd, serviceAccount, params)
	})
}

func getIAMServiceAccountWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, params *getCmdParams) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	serviceAccount := &api.ClusterIAMServiceAccount{}

	cfg.IAM.WithOIDC = api.Enabled()
	cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, serviceAccount)

	params := &getCmdParams{}

	cmd.SetDescription("iamserviceaccount", "Get iamserviceaccount(s)", "", "iamserviceaccounts")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return runFunc(cmd, serviceAccount, params)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")

		fs.StringVar(&serviceAccount.Name, "name", "", "name of the iamserviceaccount")
		fs.StringVar(&serviceAccount.Namespace, "namespace", "default", "namespace where to retrieve the iamserviceaccount")

		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doGetIAMServiceAccount(cmd *cmdutils.Cmd, serviceAccount *api.ClusterIAMServiceAccount, params *getCmdParams) error {
	if err := cmdutils.NewGetIAMServiceAccountLoader(cmd, serviceAccount).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	remoteServiceAccounts, err := stackManager.GetIAMServiceAccounts()
	if err != nil {
		return errors.Wrap(err, "getting iamserviceaccounts")
	}
	// we will show user the object based on given config file,
	// and what we have learned about the iamserviceaccounts;
	// that is not ideal, but we don't have a better option yet
	cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{}

	saFilter := cmdutils.NewIAMServiceAccountFilter()

	if cmd.ClusterConfigFile == "" {
		// reset defaulted fields to avoid output being a complete lie
		cfg.VPC = nil
		cfg.CloudWatch = nil
		// only return the iamserciceaccount that user asked for
		var notFoundErr error
		if serviceAccount.Name != "" { // name was given
			notFoundErr = fmt.Errorf("iamserviceaccount %q not found", serviceAccount.NameString())
			saFilter.AppendIncludeNames(serviceAccount.NameString())
		} else if cmd.CobraCommand.Flag("namespace").Changed { // only namespace was given
			notFoundErr = fmt.Errorf("no iamserviceaccounts found in namespace %q", serviceAccount.Namespace)
			err = saFilter.AppendIncludeGlobs(remoteServiceAccounts, serviceAccount.Namespace+"/*")
			if err != nil {
				return fmt.Errorf("unable to append include globs in namespace %q", serviceAccount.Namespace)
			}
		}
		saSubset, _ := saFilter.MatchAll(remoteServiceAccounts)
		if saSubset.Len() == 0 {
			return notFoundErr
		}
	}

	err = saFilter.ForEach(remoteServiceAccounts, func(_ int, remoteServiceAccount *api.ClusterIAMServiceAccount) error {
		cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, remoteServiceAccount)
		return nil
	})
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
		obj = cfg.IAM.ServiceAccounts
	} else {
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
