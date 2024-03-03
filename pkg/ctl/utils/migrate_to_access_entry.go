package utils

import (
	"context"
	"fmt"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	accessentryactions "github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func migrateAccessEntryCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("migrate-to-access-entry", "Migrates aws-auth to API authentication mode for the cluster", "")

	var options accessentryactions.AccessEntryMigrationOptions
	cmd.FlagSetGroup.InFlagSet("Migrate to Access Entry", func(fs *pflag.FlagSet) {
		fs.StringVar(&options.TargetAuthMode, "target-authentication-mode", "API_AND_CONFIG_MAP", "Target Authentication mode of migration")
		fs.BoolVar(&options.RemoveOIDCProviderTrustRelationship, "remove-aws-auth", false, "Remove aws-auth from cluster")
		fs.BoolVar(&options.Approve, "approve", false, "Apply the changes")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &options.Timeout)
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doMigrateToAccessEntry(cmd, options)
	}
}

func doMigrateToAccessEntry(cmd *cmdutils.Cmd, options accessentryactions.AccessEntryMigrationOptions) error {
	cfg := cmd.ClusterConfig
	cmd.ClusterConfig.AccessConfig.AuthenticationMode = ekstypes.AuthenticationMode(options.TargetAuthMode)
	tgAuthMode := cmd.ClusterConfig.AccessConfig.AuthenticationMode

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

	if tgAuthMode != ekstypes.AuthenticationModeApi && tgAuthMode != ekstypes.AuthenticationModeApiAndConfigMap {
		return fmt.Errorf("target authentication mode is invalid")
	}

	curAuthMode := ctl.GetClusterState().AccessConfig.AuthenticationMode

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)
	aeCreator := accessentryactions.Creator{
		ClusterName:  cmd.ClusterConfig.Metadata.Name,
		StackCreator: stackManager,
	}

	return accessentryactions.NewMigrator(
		cfg.Metadata.Name,
		ctl.AWSProvider.EKS(),
		ctl.AWSProvider.IAM(),
		clientSet,
		aeCreator,
		curAuthMode,
		tgAuthMode,
	).MigrateToAccessEntry(ctx, options)
}
