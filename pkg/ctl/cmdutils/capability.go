package cmdutils

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var capabilityFlagsIncompatibleWithConfigFile = []string{
	"name",
	"type",
	"role-arn",
	"delete-propagation-policy",
	"tags",
}

func NewCreateCapabilityLoader(cmd *Cmd, capability *api.Capability) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(capabilityFlagsIncompatibleWithConfigFile...)
	l.validateWithConfigFile = func() error {
		if len(cmd.ClusterConfig.Capabilities) == 0 {
			return fmt.Errorf("no capabilities specified")
		}
		for _, c := range cmd.ClusterConfig.Capabilities {
			if err := c.Validate(); err != nil {
				return err
			}
		}
		return nil
	}
	l.validateWithoutConfigFile = func() error {
		if err := validateCluster(cmd); err != nil {
			return err
		}
		if capability.Name == "" {
			return fmt.Errorf("must specify capability name")
		}
		if capability.Type == "" {
			return fmt.Errorf("must specify capability type")
		}
		return capability.Validate()
	}
	return l
}

func NewDeleteCapabilityLoader(cmd *Cmd, capability *api.Capability) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(capabilityFlagsIncompatibleWithConfigFile...)
	l.validateWithConfigFile = func() error {
		if len(cmd.ClusterConfig.Capabilities) == 0 {
			return fmt.Errorf("no capabilities specified")
		}
		for _, c := range cmd.ClusterConfig.Capabilities {
			if c.Name == "" {
				return fmt.Errorf("must specify capability name")
			}
		}
		return nil
	}
	l.validateWithoutConfigFile = func() error {
		if err := validateCluster(cmd); err != nil {
			return err
		}
		if capability.Name == "" {
			return fmt.Errorf("must specify capability name")
		}
		return nil
	}
	return l
}

func NewGetCapabilityLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.validateWithConfigFile = func() error {
		return nil
	}
	l.validateWithoutConfigFile = func() error {
		return validateCluster(cmd)
	}
	return l
}