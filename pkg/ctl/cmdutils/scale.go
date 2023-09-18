package cmdutils

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewScaleNodeGroupLoader will load config or use flags for 'eksctl scale nodegroup'
func NewScaleNodeGroupLoader(cmd *Cmd, ng *api.NodeGroupBase) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"nodes",
		"nodes-min",
		"nodes-max",
	)
	l.flagsIncompatibleWithConfigFile.Delete("name")

	l.validateWithConfigFile = func() error {
		if err := validateNameArgument(cmd, ng); err != nil {
			return err
		}

		loadedNG, err := l.ClusterConfig.FindNodegroup(ng.Name)
		if err != nil {
			return err
		}

		if err := ValidateNumberOfNodes(loadedNG); err != nil {
			return err
		}
		*ng = *loadedNG
		l.Plan = false
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.Metadata.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}

		if err := validateNameArgument(cmd, ng); err != nil {
			return err
		}

		if err := validateNumberOfNodesCLI(ng); err != nil {
			return err
		}
		l.Plan = false
		return nil
	}

	return l
}

// NewScaleAllNodeGroupLoader will load config or use flags for 'eksctl scale nodegroup'
func NewScaleAllNodeGroupLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.flagsIncompatibleWithConfigFile.Insert(
		"nodes",
		"nodes-min",
		"nodes-max",
	)

	l.validateWithConfigFile = func() error {
		if len(l.ClusterConfig.AllNodeGroups()) == 0 {
			return fmt.Errorf("no nodegroups found in config file")
		}
		l.Plan = false
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		return fmt.Errorf("a config file is required when --name is not set or when scaling multiple nodegroups")
	}

	return l
}

func validateNameArgument(cmd *Cmd, ng *api.NodeGroupBase) error {
	if ng.Name != "" && cmd.NameArg != "" {
		return ErrFlagAndArg("--name", ng.Name, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		ng.Name = cmd.NameArg
	}

	return nil
}

// ValidateNumberOfNodes validates the scaling config of a nodegroup.
func ValidateNumberOfNodes(ng *api.NodeGroupBase) error {
	if ng.ScalingConfig == nil {
		ng.ScalingConfig = &api.ScalingConfig{}
	}

	if ng.DesiredCapacity == nil || *ng.DesiredCapacity < 0 {
		return fmt.Errorf("number of nodes must be 0 or greater")
	}

	if ng.MaxSize != nil && *ng.MaxSize < 0 {
		return fmt.Errorf("maximum number of nodes must be 0 or greater")
	}

	if ng.MaxSize != nil && ng.MinSize != nil && (*ng.MinSize > *ng.DesiredCapacity || *ng.MaxSize < *ng.DesiredCapacity) {
		return fmt.Errorf("number of nodes must be within range of min nodes and max nodes")
	}

	if ng.MaxSize != nil && *ng.MaxSize < *ng.DesiredCapacity {
		return fmt.Errorf("maximum number of nodes must be greater than or equal to number of nodes")
	}

	if ng.MinSize != nil && *ng.MinSize < 0 {
		return fmt.Errorf("minimum number of nodes must be 0 or greater")
	}

	if ng.MinSize != nil && *ng.MinSize > *ng.DesiredCapacity {
		return fmt.Errorf("minimum number of nodes must be less than or equal to number of nodes")
	}

	return nil
}

// only 1 of desired/min/max has to be set on the cli
func validateNumberOfNodesCLI(ng *api.NodeGroupBase) error {
	if ng.ScalingConfig == nil {
		ng.ScalingConfig = &api.ScalingConfig{}
	}

	if ng.DesiredCapacity == nil && ng.MinSize == nil && ng.MaxSize == nil {
		return fmt.Errorf("at least one of minimum, maximum and desired nodes must be set")
	}

	if ng.DesiredCapacity != nil && *ng.DesiredCapacity < 0 {
		return fmt.Errorf("number of nodes must be 0 or greater")
	}

	if ng.MinSize != nil && *ng.MinSize < 0 {
		return fmt.Errorf("minimum of nodes must be 0 or greater")
	}
	if ng.MaxSize != nil && *ng.MaxSize < 0 {
		return fmt.Errorf("maximum of nodes must be 0 or greater")
	}

	if ng.MaxSize != nil && ng.MinSize != nil && *ng.MaxSize < *ng.MinSize {
		return fmt.Errorf("maximum number of nodes must be greater than minimum number of nodes")
	}

	if ng.MaxSize != nil && ng.DesiredCapacity != nil && *ng.MaxSize < *ng.DesiredCapacity {
		return fmt.Errorf("maximum number of nodes must be greater than or equal to number of nodes")
	}

	if ng.MinSize != nil && ng.DesiredCapacity != nil && *ng.MinSize > *ng.DesiredCapacity {
		return fmt.Errorf("minimum number of nodes must be fewer than or equal to number of nodes")
	}
	return nil
}
