package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateZonalShiftConfig(cmd *cmdutils.Cmd, handler func(*cmdutils.Cmd) error) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-zonal-shift-config", "update zonal shift config", "update zonal shift config on a cluster")

	var enableZonalShift bool
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewZonalShiftConfigLoader(cmd).Load(); err != nil {
			return err
		}
		if cmd.ClusterConfigFile == "" {
			cfg.ZonalShiftConfig = &api.ZonalShiftConfig{
				Enabled: &enableZonalShift,
			}
		}
		return handler(cmd)
	}

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		fs.BoolVar(&enableZonalShift, "enable-zonal-shift", true, "Enable zonal shift on a cluster")
	})

}

func updateZonalShiftConfigCmd(cmd *cmdutils.Cmd) {
	updateZonalShiftConfig(cmd, doUpdateZonalShiftConfig)
}

func doUpdateZonalShiftConfig(cmd *cmdutils.Cmd) error {
	cfg := cmd.ClusterConfig
	ctx := context.Background()
	if cfg.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}
	ctl, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	makeZonalShiftStatus := func(enabled *bool) string {
		if api.IsEnabled(enabled) {
			return "enabled"
		}
		return "disabled"
	}
	if zsc := ctl.Status.ClusterInfo.Cluster.ZonalShiftConfig; zsc != nil && *zsc.Enabled == api.IsEnabled(cfg.ZonalShiftConfig.Enabled) {
		logger.Info("zonal shift is already %s", makeZonalShiftStatus(zsc.Enabled))
		return nil
	}
	if err := ctl.UpdateClusterConfig(ctx, &eks.UpdateClusterConfigInput{
		Name: aws.String(cfg.Metadata.Name),
		ZonalShiftConfig: &ekstypes.ZonalShiftConfigRequest{
			Enabled: cfg.ZonalShiftConfig.Enabled,
		},
	}); err != nil {
		return fmt.Errorf("updating zonal shift config: %w", err)
	}
	logger.Info("zonal shift %s successfully", makeZonalShiftStatus(cfg.ZonalShiftConfig.Enabled))
	return nil
}
