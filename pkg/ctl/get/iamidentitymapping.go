package get

import (
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
)

func getIAMIdentityMappingCmd(g *cmdutils.Grouping) *cobra.Command {
	cfg := api.NewClusterConfig()
	cp := cmdutils.NewCommonParams(cfg)

	var role string

	cp.Command = &cobra.Command{
		Use:   "iamidentitymapping",
		Short: "Get IAM identity mapping(s)",
		Run: func(_ *cobra.Command, _ []string) {
			if err := doGetIAMIdentityMapping(cp, role); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cp.Command)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&role, "role", "", "ARN of the IAM role")
		cmdutils.AddNameFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, cp.ProviderConfig)
		cmdutils.AddCommonFlagsForGetCmd(fs, &chunkSize, &output)
	})

	cmdutils.AddCommonFlagsForAWS(group, cp.ProviderConfig, false)

	group.AddTo(cp.Command)
	return cp.Command
}

func doGetIAMIdentityMapping(cp *cmdutils.CommonParams, role string) error {
	cfg := cp.ClusterConfig

	ctl := eks.New(cp.ProviderConfig, cfg)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet("--cluster")
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

	printer, err := printers.NewPrinter(output)
	if err != nil {
		return err
	}
	if output == "table" {
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
