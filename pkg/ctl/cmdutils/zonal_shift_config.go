package cmdutils

import (
	"errors"
	"fmt"
)

// NewZonalShiftConfigLoader creates a new loader for zonal shift config.
func NewZonalShiftConfigLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(
		"enable-zonal-shift",
		"cluster",
	)

	l.validateWithConfigFile = func() error {
		if cmd.NameArg != "" {
			return fmt.Errorf("config file and enable-zonal-shift %s", IncompatibleFlags)
		}
		if l.ClusterConfig.ZonalShiftConfig == nil || l.ClusterConfig.ZonalShiftConfig.Enabled == nil {
			return errors.New("field zonalShiftConfig.enabled is required")
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		if !cmd.CobraCommand.Flag("enable-zonal-shift").Changed {
			return errors.New("--enable-zonal-shift is required when a config file is not specified")
		}
		return nil
	}
	return l
}
