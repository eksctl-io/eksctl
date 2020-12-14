package disassociate

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func disassociateIdentityProvider(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("identityprovider", "Disassociate an identity provider from a cluster", "")

	var provider identityproviders.DisassociateIdentityProvider

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doDisassociateIdentityProvider(cmd, provider)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmdutils.AddWaitFlag(fs, &cmd.Wait, "deletion of providers")
		fs.StringVar(&provider.Name, "name", "", "name of the provider to delete")
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func newDisassociateIdentityProviderLoader(cmd *cmdutils.Cmd, provider identityproviders.DisassociateIdentityProvider) cmdutils.ClusterConfigLoader {
	l := cmdutils.NewConfigLoaderBuilder()

	l.ValidateWithConfigFile(func(cmd *cmdutils.Cmd) error {
		if len(cmd.ClusterConfig.IdentityProviders) == 0 {
			return fmt.Errorf("No identity providers provided")
		}
		return nil
	})
	l.ValidateWithoutConfigFile(func(cmd *cmdutils.Cmd) error {
		if len(cmd.ClusterConfig.IdentityProviders) == 0 && provider.Name == "" {
			return fmt.Errorf("No identity providers provided")
		}
		return nil
	})

	return l.Build(cmd)
}

func doDisassociateIdentityProvider(cmd *cmdutils.Cmd, provider identityproviders.DisassociateIdentityProvider) error {
	if err := newDisassociateIdentityProviderLoader(cmd, provider).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	meta := cmd.ClusterConfig.Metadata

	ctl, err := cmd.NewCtl()
	if err != nil {
		return err
	}

	cmdutils.LogRegionAndVersionInfo(meta)

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	providers := []identityproviders.DisassociateIdentityProvider{}
	if provider.Name == "" {
		for _, generalIdP := range cfg.IdentityProviders {
			var provider identityproviders.DisassociateIdentityProvider
			switch idP := generalIdP.Inner().(type) {
			case *api.OIDCIdentityProvider:
				provider = identityproviders.DisassociateIdentityProvider{
					Name: idP.Name,
					Type: string(api.OIDCIdentityProviderType),
				}
			default:
				return errors.New("can't disassociate provider")
			}
			providers = append(
				providers,
				provider,
			)
		}
	} else {
		providers = []identityproviders.DisassociateIdentityProvider{
			provider,
		}
	}
	options := identityproviders.DisassociateIdentityProvidersOptions{
		Providers: providers,
	}
	if cmd.Wait {
		timeout := ctl.Provider.WaitTimeout()
		options.WaitTimeout = &timeout
	}

	manager := identityproviders.NewIdentityProviderManager(
		*cfg.Metadata,
		ctl.Provider.EKS(),
	)

	return manager.Disassociate(options)
}
