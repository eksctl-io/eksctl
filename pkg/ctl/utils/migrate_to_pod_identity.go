package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func migrateToPodIdentityCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("migrate-to-pod-identity", "Migrates all IRSA related config for a cluster to an equivalent pod identity associations config", "")

	var options podidentityassociation.PodIdentityMigrationOptions
	cmd.FlagSetGroup.InFlagSet("Authentication mode", func(fs *pflag.FlagSet) {
		fs.BoolVar(&options.RemoveOIDCProviderTrustRelationship, "remove-oidc-provider-trust-relationship", false, "Remove existing IRSAv1 OIDC provided entities")
		fs.BoolVar(&options.Approve, "approve", false, "Apply the changes")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &options.Timeout)
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doMigrateToPodIdentity(cmd, options)
	}
}

func doMigrateToPodIdentity(cmd *cmdutils.Cmd, options podidentityassociation.PodIdentityMigrationOptions) error {
	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	ctx := context.Background()
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if ok, err := ctl.CanOperate(cfg); !ok {
		return err
	}

	oidc, err := ctl.NewOpenIDConnectManager(ctx, cfg)
	if err != nil {
		return err
	}

	providerExists, err := oidc.CheckProviderExists(ctx)
	if err != nil {
		return err
	}

	if !providerExists {
		logger.Warning("no IAM OIDC provider associated with cluster, hence no iamserviceaccounts to be migrated")
		return nil
	}

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	addonManager, err := addon.New(cfg, ctl.AWSProvider.EKS(), nil, false, nil, nil)
	if err != nil {
		return fmt.Errorf("initializing addon creator %w", err)
	}

	stackManager := ctl.NewStackManager(cfg)
	iamRoleCreator := &podidentityassociation.IAMRoleCreator{
		ClusterName:  cfg.Metadata.Name,
		StackCreator: stackManager,
	}
	return podidentityassociation.NewMigrator(
		cfg.Metadata.Name,
		ctl.AWSProvider.EKS(),
		ctl.AWSProvider.IAM(),
		stackManager,
		clientSet,
		&addonCreator{
			addonManager:   addonManager,
			iamRoleCreator: iamRoleCreator,
		},
	).MigrateToPodIdentity(ctx, options)
}

type addonCreator struct {
	addonManager   *addon.Manager
	iamRoleCreator addon.IAMRoleCreator
}

func (a *addonCreator) Create(ctx context.Context, addon *api.Addon, waitTimeout time.Duration) error {
	return a.addonManager.Create(ctx, addon, a.iamRoleCreator, waitTimeout)
}
