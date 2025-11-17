package utils

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateClusterDeletionProtectionCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	var enabled string

	cmd.SetDescription("deletion-protection", "Update cluster deletion protection", "")

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		fs.StringVar(&enabled, "enabled", "", "Enable or disable deletion protection (true|false)")
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)

		if enabled == "" {
			return fmt.Errorf("--enabled flag is required (true|false)")
		}

		val, err := strconv.ParseBool(enabled)
		if err != nil {
			return fmt.Errorf("--enabled must be 'true' or 'false', got: %s", enabled)
		}

		cfg.DeletionProtection = &val

		return doUpdateClusterDeletionProtection(cmd)
	}
}

func doUpdateClusterDeletionProtection(cmd *cmdutils.Cmd) error {
	ctx := context.Background()
	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	cfg := cmd.ClusterConfig
	if cfg.Metadata.Name == "" {
		return fmt.Errorf("cluster name is required")
	}

	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	if cmd.Plan {
		logger.Critical("--dry-run is not supported for this command")
		return nil
	}

	action := "disabling"
	if cfg.DeletionProtection != nil && *cfg.DeletionProtection {
		action = "enabling"
	}

	logger.Info("%s deletion protection for cluster %q", action, cfg.Metadata.Name)
	return ctl.UpdateClusterConfigForDeletionProtection(ctx, cfg)
}
