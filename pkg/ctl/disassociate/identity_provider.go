package disassociate

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

var defaultDisassociateTimeout = 35 * time.Minute

func disassociateIdentityProvider(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("identityprovider", "Disassociate an identity provider from a cluster", "")

	var provider identityproviders.DisassociateIdentityProvider
	var timeout time.Duration

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doDisassociateIdentityProvider(cmd, provider, timeout)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmdutils.AddWaitFlag(fs, &cmd.Wait, "providers to disassociate")
		cmdutils.AddTimeoutFlagWithValue(fs, &timeout, defaultDisassociateTimeout)
		fs.StringVar(&provider.Name, "name", "", "name of the provider to disassociate")
		fs.StringVar(&provider.Type, "type", "", "type of the provider to disassociate")
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

func doDisassociateIdentityProvider(cmd *cmdutils.Cmd, provider identityproviders.DisassociateIdentityProvider, timeout time.Duration) error {
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

	manager := identityproviders.NewIdentityProviderManager(
		*cfg.Metadata,
		ctl.Provider.EKS(),
	)

	options := identityproviders.DisassociateIdentityProvidersOptions{
		Providers: providers,
	}
	if cmd.Wait {
		options.WaitTimeout = &timeout
	}

	return manager.Disassociate(options)
}
