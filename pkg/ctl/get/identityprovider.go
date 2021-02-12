package get

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getIdentityProvider(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg
	params := getCmdParams{}

	cmd.SetDescription("identityprovider", "Describe identity providers for cluster authentication and authorization", "")

	var name = ""

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doGetIdentityProvider(cmd, params, name)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)

		fs.StringVar(&name, "name", "", "name of the provider to delete")
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doGetIdentityProvider(cmd *cmdutils.Cmd, params getCmdParams, name string) error {
	builder := cmdutils.NewConfigLoaderBuilder()
	if err := builder.Build(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		cmdutils.LogRegionAndVersionInfo(cfg.Metadata)
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	manager := identityproviders.NewManager(
		*cfg.Metadata,
		ctl.Provider.EKS(),
	)

	summaries, err := manager.Get(
		identityproviders.GetIdentityProvidersOptions{Name: name},
	)
	if err != nil {
		return err
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}

	if params.output == printers.TableType {
		addIdentityProviderTableColumns(printer.(*printers.TablePrinter), len(summaries))
	}

	if err := printer.PrintObjWithKind("identity provider summary", summaries, os.Stdout); err != nil {
		return err
	}

	return nil
}

func addIdentityProviderTableColumns(printer *printers.TablePrinter, num int) {
	printer.AddColumn("NAME", func(s identityproviders.Summary) string {
		return s.Name
	})
	printer.AddColumn("TYPE", func(s identityproviders.Summary) string {
		return string(s.Type)
	})
	printer.AddColumn("CLIENT_ID", func(s identityproviders.Summary) string {
		return s.ClientID
	})
	printer.AddColumn("ISSUER_URL", func(s identityproviders.Summary) string {
		return s.IssuerURL
	})
	printer.AddColumn("ARN", func(s identityproviders.Summary) string {
		return s.Arn
	})
	printer.AddColumn("STATUS", func(s identityproviders.Summary) string {
		return s.Status
	})
}
