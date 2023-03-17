package associate

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

var defaultAssociateTimeout = 35 * time.Minute

func associateIdentityProvider(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("identityprovider", "Associate an identity provider with a cluster", "")

	var timeout time.Duration

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doAssociateIdentityProvider(cmd, timeout)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)

		cmdutils.AddWaitFlag(fs, &cmd.Wait, "providers to associate")
		cmdutils.AddTimeoutFlagWithValue(fs, &timeout, defaultAssociateTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
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

func doAssociateIdentityProvider(cmd *cmdutils.Cmd, timeout time.Duration) error {
	if err := newAssociateIdentityProviderLoader(cmd).Load(); err != nil {
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

	manager := identityproviders.NewManager(
		*cfg.Metadata,
		ctl.AWSProvider.EKS(),
	)

	options := identityproviders.AssociateIdentityProvidersOptions{
		Providers: cfg.IdentityProviders,
	}
	if cmd.Wait {
		options.WaitTimeout = timeout
	}

	return manager.Associate(ctx, options)
}
