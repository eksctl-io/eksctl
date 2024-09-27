package cmdutils

import (
	"errors"
	"fmt"
	"slices"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewAutonomousModeLoader creates a new loader for Autonomous Mode.
func NewAutonomousModeLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.validateWithConfigFile = func() error {
		meta := cmd.ClusterConfig.Metadata
		if meta.Name == "" {
			return ErrMustBeSet("metadata.name")
		}
		if meta.Region == "" {
			return ErrMustBeSet("metadata.region")
		}
		cc := cmd.ClusterConfig.AutonomousModeConfig
		if cc == nil {
			return ErrMustBeSet("autonomousModeConfig")
		}
		if cc.Enabled == nil {
			return errors.New("autonomousModeConfig.enabled must be set to either true or false when updating Autonomous Mode config")
		}
		api.SetClusterConfigDefaults(cmd.ClusterConfig)
		if err := api.ValidateAutonomousModeConfig(cmd.ClusterConfig); err != nil {
			return err
		}
		const (
			drainParallel   = "drain-parallel"
			drainNodeGroups = "drain-all-nodegroups"
		)

		if slices.ContainsFunc([]string{drainParallel, drainNodeGroups}, func(flagName string) bool {
			return cmd.CobraCommand.Flag(flagName).Changed
		}) {
			if api.IsDisabled(cc.Enabled) {
				return fmt.Errorf("cannot specify --%s or --%s when autonomousModeConfig.enabled is false", drainParallel, drainNodeGroups)
			}
			if !cc.HasNodePools() {
				return fmt.Errorf("cannot specify --%s or --%s when autonomousModeConfig.nodePools is empty", drainParallel, drainNodeGroups)
			}
		}
		return nil
	}
	l.validateWithoutConfigFile = func() error {
		return ErrMustBeSet("--config-file")
	}
	return l
}
