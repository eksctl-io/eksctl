package get

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getIAMIdentityMappingCmd(rc *cmdutils.ResourceCmd) {
	cfg := api.NewClusterConfig()
	rc.ClusterConfig = cfg

	var arn string

	params := &getCmdParams{}

	rc.SetDescription("iamidentitymapping", "Get IAM identity mapping(s)", "")

	rc.SetRunFunc(func() error {
		return doGetIAMIdentityMapping(rc, params, arn)
	})

	rc.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&arn, "arn", "", "ARN of the IAM role or user")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, rc.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddConfigFileFlag(fs, &rc.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(rc.FlagSetGroup, rc.ProviderConfig, false)
}

func doGetIAMIdentityMapping(rc *cmdutils.ResourceCmd, params *getCmdParams, arn string) error {
	if err := cmdutils.NewMetadataLoader(rc).Load(); err != nil {
		return err
	}

	cfg := rc.ClusterConfig

	ctl := eks.New(rc.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
	}

	if err := ctl.GetCredentials(cfg); err != nil {
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
	identities, err := acm.Identities()
	if err != nil {
		return err
	}
	if arn != "" {
		identities = identities.Get(arn)
		// If a filter was given, we error if none was found
		if len(identities) == 0 {
			return fmt.Errorf("no iamidentitymapping with arn %q found", arn)
		}
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}
	if params.output == "table" {
		addIAMIdentityMappingTableColumns(printer.(*printers.TablePrinter))
	}

	if err := printer.PrintObjWithKind("iamidentitymappings", identities, os.Stdout); err != nil {
		return err
	}

	return nil
}

func addIAMIdentityMappingTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("ARN", func(r authconfigmap.MapIdentity) string {
		return r.ARN
	})
	printer.AddColumn("USERNAME", func(r authconfigmap.MapIdentity) string {
		return r.Username
	})
	printer.AddColumn("GROUPS", func(r authconfigmap.MapIdentity) string {
		return strings.Join(r.Groups, ",")
	})
}
