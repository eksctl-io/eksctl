package get

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getIAMIdentityMappingCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var role string

	params := &getCmdParams{}

	cmd.SetDescription("iamidentitymapping", "Get IAM identity mapping(s)", "")

	cmd.SetRunFunc(func() error {
		return doGetIAMIdentityMapping(cmd, params, role)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&role, "role", "", "ARN of the IAM role")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &params.chunkSize, &params.output)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doGetIAMIdentityMapping(cmd *cmdutils.Cmd, params *getCmdParams, role string) error {
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
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

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--name")
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
	roles, err := acm.Roles()
	if err != nil {
		return err
	}
	if role != "" {
		roles = roles.Get(role)
		// If a filter was given, we error if none was found
		if len(roles) == 0 {
			return fmt.Errorf("no iamidentitymapping with role %q found", role)
		}
	}

	printer, err := printers.NewPrinter(params.output)
	if err != nil {
		return err
	}
	if params.output == "table" {
		addIAMIdentityMappingTableColumns(printer.(*printers.TablePrinter))
	}

	if err := printer.PrintObjWithKind("iamidentitymappings", roles, os.Stdout); err != nil {
		return err
	}

	return nil
}

func addIAMIdentityMappingTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("ROLE", func(r authconfigmap.MapRole) string {
		return r.RoleARN
	})
	printer.AddColumn("USERNAME", func(r authconfigmap.MapRole) string {
		return r.Username
	})
	printer.AddColumn("GROUPS", func(r authconfigmap.MapRole) string {
		return strings.Join(r.Groups, ",")
	})
}
