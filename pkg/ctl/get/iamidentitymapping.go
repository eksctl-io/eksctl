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
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()
	var roleFlag string
	cmd := &cobra.Command{
		Use:   "iamidentitymapping",
		Short: "Get IAM identity mapping(s)",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doGetIAMIdentityMapping(p, cfg, roleFlag); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}
	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&roleFlag, "role", "", "ARN of the IAM role")
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddCommonFlagsForGetCmd(fs, &chunkSize, &output)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, false)

	group.AddTo(cmd)

	return cmd
}

func doGetIAMIdentityMapping(p *api.ProviderConfig, cfg *api.ClusterConfig, roleFlag string) error {
	ctl := eks.New(p, cfg)

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
	if roleFlag != "" {
		roles = roles.Get(roleFlag)
		// If a filter was given, we error if none was found
		if len(roles) == 0 {
			return fmt.Errorf("no iamidentitymapping with role %q found", roleFlag)
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
