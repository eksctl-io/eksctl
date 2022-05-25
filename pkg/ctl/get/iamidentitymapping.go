package get

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getIAMIdentityMappingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var arn string

	params := &getCmdParams{}

	cmd.SetDescription("iamidentitymapping", "Get IAM identity mapping(s)", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doGetIAMIdentityMapping(cmd, params, arn)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddIAMIdentityMappingARNFlags(fs, cmd, &arn, "get")
		cmdutils.AddClusterFlagWithDeprecated(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func doGetIAMIdentityMapping(cmd *cmdutils.Cmd, params *getCmdParams, arn string) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	if params.output != printers.TableType {
		logger.Writer = os.Stderr
	}

	ctl, err := cmd.NewProviderForExistingCluster(context.Background())
	if err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}
	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}
	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}
	identities, err := acm.GetIdentities()
	if err != nil {
		return err
	}

	if arn != "" {
		var selectedIdentities []iam.Identity

		for _, identity := range identities {
			if identity.ARN() == arn {
				selectedIdentities = append(selectedIdentities, identity)
			}
		}

		identities = selectedIdentities
		// If a filter was given, we error if none was found
		if len(identities) == 0 {
			return fmt.Errorf("no iamidentitymapping with arn %q found", arn)
		}
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}
	if params.output == printers.TableType {
		addIAMIdentityMappingTableColumns(printer.(*printers.TablePrinter))
	}

	return printer.PrintObjWithKind("iamidentitymappings", identities, os.Stdout)
}

func addIAMIdentityMappingTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("ARN", func(r iam.Identity) string {
		return r.ARN()
	})
	printer.AddColumn("USERNAME", func(r iam.Identity) string {
		return r.Username()
	})
	printer.AddColumn("GROUPS", func(r iam.Identity) string {
		return strings.Join(r.Groups(), ",")
	})
	printer.AddColumn("ACCOUNT", func(r iam.Identity) string {
		return r.Account()
	})
}
