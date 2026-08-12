package cmdutils

import (
	"errors"
)

// NewControlPlaneComponentConfigLoader creates a new loader for control plane component config.
func NewControlPlaneComponentConfigLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(
		"cluster",
	)

	l.validateWithConfigFile = func() error {
		if cmd.NameArg != "" {
			return errors.New("name argument is not supported")
		}
		cfg := l.ClusterConfig
		if cfg.KubeAPIServerConfig == nil && cfg.KubeSchedulerConfig == nil && cfg.KubeControllerManagerConfig == nil {
			return errors.New("at least one of kubeAPIServerConfig, kubeSchedulerConfig or kubeControllerManagerConfig is required")
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		return errors.New("--config-file/-f is required to update control plane component config")
	}

	return l
}
