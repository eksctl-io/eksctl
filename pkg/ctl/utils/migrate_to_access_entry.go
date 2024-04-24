package utils

import (
	"context"

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

	var options accessentryactions.MigrationOptions
	cmd.FlagSetGroup.InFlagSet("Migrate to Access Entry", func(fs *pflag.FlagSet) {
		fs.StringVar(&options.TargetAuthMode, "target-authentication-mode", "API_AND_CONFIG_MAP", "Target Authentication mode of migration")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &options.Timeout)
		cmdutils.AddApproveFlag(fs, cmd)
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		options.Approve = !cmd.Plan
		return doMigrateToAccessEntry(cmd, options)
	}
}

func doMigrateToAccessEntry(cmd *cmdutils.Cmd, options accessentryactions.MigrationOptions) error {
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

	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)
	aeCreator := &accessentryactions.Creator{
		ClusterName:  cmd.ClusterConfig.Metadata.Name,
		StackCreator: stackManager,
	}
	aeGetter := accessentryactions.NewGetter(cfg.Metadata.Name, ctl.AWSProvider.EKS())

	if err := accessentryactions.NewMigrator(
		cfg.Metadata.Name,
		ctl.AWSProvider.EKS(),
		ctl.AWSProvider.IAM(),
		clientSet,
		aeCreator,
		aeGetter,
		ctl.GetClusterState().AccessConfig.AuthenticationMode,
		ekstypes.AuthenticationMode(options.TargetAuthMode),
	).MigrateToAccessEntry(ctx, options); err != nil {
		return err
	}

	cmdutils.LogPlanModeWarning(cmd.Plan)
	return nil
}
