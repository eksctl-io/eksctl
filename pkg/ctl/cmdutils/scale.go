package cmdutils

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewScaleNodeGroupLoader will load config or use flags for 'eksctl scale nodegroup'
func NewScaleNodeGroupLoader(cmd *Cmd, ng *api.NodeGroup) ClusterConfigLoader {
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

		loadedNG := l.ClusterConfig.FindNodegroup(ng.Name)
		if loadedNG == nil {
			return fmt.Errorf("node group %s not found", ng.Name)
		}

		if err := validateNumberOfNodes(loadedNG); err != nil {
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

func validateNameArgument(cmd *Cmd, ng *api.NodeGroup) error {
	if ng.Name != "" && cmd.NameArg != "" {
		return ErrFlagAndArg("--name", ng.Name, cmd.NameArg)
	}

	if cmd.NameArg != "" {
		ng.Name = cmd.NameArg
	}

	if ng.Name == "" {
		return ErrMustBeSet("--name")
	}

	return nil
}

func validateNumberOfNodes(ng *api.NodeGroup) error {
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

//only 1 of desired/min/max has to be set on the cli
func validateNumberOfNodesCLI(ng *api.NodeGroup) error {
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
