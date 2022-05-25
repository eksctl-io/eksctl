package disassociate

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var defaultDisassociateTimeout = 35 * time.Minute

type cliProvidedIDP struct {
	Name string
	Type string
}

func disassociateIdentityProvider(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("identityprovider", "Disassociate an identity provider from a cluster", "")

	var cliProvidedIDP cliProvidedIDP
	var timeout time.Duration

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doDisassociateIdentityProvider(cmd, cliProvidedIDP, timeout)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmdutils.AddWaitFlag(fs, &cmd.Wait, "providers to disassociate")
		cmdutils.AddTimeoutFlagWithValue(fs, &timeout, defaultDisassociateTimeout)
		fs.StringVar(&cliProvidedIDP.Name, "name", "", "name of the provider to disassociate")
		fs.StringVar(&cliProvidedIDP.Type, "type", "", "type of the provider to disassociate")
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func newDisassociateIdentityProviderLoader(cmd *cmdutils.Cmd, cliProvidedIDP cliProvidedIDP) cmdutils.ClusterConfigLoader {
	l := cmdutils.NewConfigLoaderBuilder()

	l.ValidateWithConfigFile(func(cmd *cmdutils.Cmd) error {
		if len(cmd.ClusterConfig.IdentityProviders) == 0 {
			return fmt.Errorf("No identity providers provided")
		}
		return nil
	})
	l.ValidateWithoutConfigFile(func(cmd *cmdutils.Cmd) error {
		if len(cmd.ClusterConfig.IdentityProviders) == 0 && cliProvidedIDP.Name == "" {
			return fmt.Errorf("No identity providers provided")
		}
		return nil
	})

	return l.Build(cmd)
}

func doDisassociateIdentityProvider(cmd *cmdutils.Cmd, cliProvidedIDP cliProvidedIDP, timeout time.Duration) error {
	if err := newDisassociateIdentityProviderLoader(cmd, cliProvidedIDP).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig

	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	providers := cliToProviders(cfg, cliProvidedIDP)

	manager := identityproviders.NewManager(
		*cfg.Metadata,
		ctl.AWSProvider.EKS(),
	)

	options := identityproviders.DisassociateIdentityProvidersOptions{
		Providers: providers,
	}
	if cmd.Wait {
		options.WaitTimeout = timeout
	}

	return manager.Disassociate(ctx, options)
}

func cliToProviders(cfg *api.ClusterConfig, cliProvidedIDP cliProvidedIDP) []identityproviders.DisassociateIdentityProvider {
	if cliProvidedIDP.Name == "" {
		var providers []identityproviders.DisassociateIdentityProvider
		for _, generalIDP := range cfg.IdentityProviders {
			var provider identityproviders.DisassociateIdentityProvider
			switch idP := (generalIDP.Inner).(type) {
			case *api.OIDCIdentityProvider:
				provider = identityproviders.DisassociateIdentityProvider{
					Name: idP.Name,
					Type: idP.Type(),
				}
			default:
				panic("unsupported identity provider")
			}
			providers = append(
				providers,
				provider,
			)
		}
		return providers
	}

	return []identityproviders.DisassociateIdentityProvider{
		{
			Name: cliProvidedIDP.Name,
			Type: api.IdentityProviderType(cliProvidedIDP.Type),
		},
	}
}
