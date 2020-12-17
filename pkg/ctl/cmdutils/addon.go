package cmdutils

import (
	"fmt"
)

var addonFlagsIncompatibleWithoutConfigFile = []string{}
var addonFlagsIncompatibleWithConfigFile = []string{
	"name",
	"version",
	"service-account-role-arn",
	"attach-policy-arn",
}

func NewCreateOrUpgradeAddonLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(addonFlagsIncompatibleWithConfigFile...)
	l.flagsIncompatibleWithoutConfigFile.Insert(addonFlagsIncompatibleWithoutConfigFile...)
	l.validateWithConfigFile = func() error {
		if len(cmd.ClusterConfig.Addons) == 0 {
			return fmt.Errorf("no addons specified")
		}
		for _, a := range cmd.ClusterConfig.Addons {
			if err := a.Validate(); err != nil {
				return err
			}
		}
		return nil
	}
	l.validateWithoutConfigFile = func() error {
		if err := validateCluster(cmd); err != nil {
			return err
		}
		if len(cmd.ClusterConfig.Addons) == 0 {
			return fmt.Errorf("no addons specified")
		}
		return cmd.ClusterConfig.Addons[0].Validate()
	}
	return l
}

func NewDeleteAddonLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(addonFlagsIncompatibleWithConfigFile...)
	l.flagsIncompatibleWithoutConfigFile.Insert(addonFlagsIncompatibleWithoutConfigFile...)
	l.validateWithConfigFile = func() error {
		if len(cmd.ClusterConfig.Addons) == 0 {
			return fmt.Errorf("no addons specified")
		}
		for _, a := range cmd.ClusterConfig.Addons {
			if a.Name == "" {
				return fmt.Errorf("must specify addon name")
			}
		}
		return nil
	}
	l.validateWithoutConfigFile = func() error {
		if err := validateCluster(cmd); err != nil {
			return err
		}

		if len(cmd.ClusterConfig.Addons) == 0 {
			return fmt.Errorf("no addons specified")
		}

		if cmd.ClusterConfig.Addons[0].Name == "" {
			return fmt.Errorf("must specify addon name")
		}
		return nil
	}
	return l
}
