package utils

import (
	"context"
	"fmt"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func migrateAccessEntryCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("migrate-to-access-entry", "Migrates aws-auth to API authentication mode for the cluster", "")

	var options accessentry.AccessEntryMigrationOptions
	cmd.FlagSetGroup.InFlagSet("Migrate to Access Entry", func(fs *pflag.FlagSet) {
		fs.BoolVar(&options.RemoveOIDCProviderTrustRelationship, "remove-aws-auth", false, "Remove aws-auth from cluster")
		fs.StringVar(&options.TargetAuthMode, "authentication-mode", "API_AND_CONFIG_MAP", "Target Authentication mode of migration")
		fs.BoolVar(&options.Approve, "approve", false, "Apply the changes")

		// fs.BoolVar(&options.SkipAgentInstallation, "skip-agent-installation", false, "Skip installing pod-identity-agent addon")
		// cmdutils.AddIAMServiceAccountFilterFlags(fs, &cmd.Include, &cmd.Exclude)
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

func doMigrateToAccessEntry(cmd *cmdutils.Cmd, options accessentry.AccessEntryMigrationOptions) error {
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

	if curAuthMode != tgAuthMode {
		logger.Info("target authentication mode %v is different than the current authentication mode %v, Updating the Cluster authentication mode", tgAuthMode, curAuthMode)
		// Add UpdateAuthentication Mode Method call here
	}

	// Get Access Entries from Cluster Provider
	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	asgetter := accessentry.NewGetter(cfg.Metadata.Name, clusterProvider.AWSProvider.EKS())
	accessEntries, err := asgetter.Get(ctx, api.ARN{})
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", accessEntries)

	// Get CONFIGMAP Entries
	clientSet, err := ctl.NewStdClientSet(cfg)
	if err != nil {
		return err
	}
	acm, err := authconfigmap.NewFromClientSet(clientSet)
	if err != nil {
		return err
	}
	cmEntries, err := acm.GetIdentities()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", cmEntries)

	// Check if any of the cmEntries are in accessEntries, and add the remaining to NeedsUpdateList

	// Perform doCreateAccessEntry() on NeedsUpdateList

	return nil
}
