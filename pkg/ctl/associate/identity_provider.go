package associate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func associateIdentityProvider(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("identityprovider", "Associate an identity provider with a cluster", "")

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doAssociateIdentityProvider(cmd)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)
}

func newAssociateIdentityProviderLoader(cmd *cmdutils.Cmd) cmdutils.ClusterConfigLoader {
	l := cmdutils.NewConfigLoaderBuilder()

	l.Validate(func(cmd *cmdutils.Cmd) error {
		if len(cmd.ClusterConfig.IdentityProviders) == 0 {
			return fmt.Errorf("No identity providers provided")
		}
		return nil
	})

	return l.Build(cmd)
}

func doAssociateIdentityProvider(cmd *cmdutils.Cmd) error {
	if err := newAssociateIdentityProviderLoader(cmd).Load(); err != nil {
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

	manager := identityproviders.NewIdentityProviderManager(
		*cfg.Metadata,
		ctl.Provider.EKS(),
	)

	return manager.Associate(identityproviders.AssociateIdentityProvidersOptions{
		Providers: cfg.IdentityProviders,
	})
}
