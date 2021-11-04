package v1alpha5

import "fmt"

// HasInstanceType returns whether some node in the group fulfils the type check
func HasInstanceType(nodeGroup *NodeGroup, hasType func(string) bool) bool {
	if hasType(nodeGroup.InstanceType) {
		return true
	}
	if nodeGroup.InstancesDistribution != nil {
		for _, instanceType := range nodeGroup.InstancesDistribution.InstanceTypes {
			if hasType(instanceType) {
				return true
			}
		}
	}
	return false
}

// HasInstanceTypeManaged returns whether some node in the managed group fulfils the type check
func HasInstanceTypeManaged(nodeGroup *ManagedNodeGroup, hasType func(string) bool) bool {
	if hasType(nodeGroup.InstanceType) {
		return true
	}
	for _, instanceType := range nodeGroup.InstanceTypes {
		if hasType(instanceType) {
			return true
		}
	}
	return false
}

// ClusterHasInstanceType checks all nodegroups and managed nodegroups for a specific instance type
func ClusterHasInstanceType(cfg *ClusterConfig, hasType func(string) bool) bool {
	for _, ng := range cfg.NodeGroups {
		if HasInstanceType(ng, hasType) {
			return true
		}
	}

	for _, mng := range cfg.ManagedNodeGroups {
		if hasType(mng.InstanceType) {
			return true
		}
	}
	return false
}

// FindNodegroup checks if the clusterConfig contains a nodegroup with the given name
func (c *ClusterConfig) FindNodegroup(name string) (*NodeGroupBase, error) {
	var foundNg []*NodeGroupBase
	for _, ng := range c.NodeGroups {
		if name == ng.NameString() {
			foundNg = append(foundNg, ng.NodeGroupBase)
		}
	}

	for _, ng := range c.ManagedNodeGroups {
		if name == ng.NameString() {
			foundNg = append(foundNg, ng.NodeGroupBase)
		}
	}

	if len(foundNg) == 0 {
		return nil, fmt.Errorf("nodegroup %s not found in config file", name)
	} else if len(foundNg) > 1 {
		return nil, fmt.Errorf("found more than 1 nodegroup with name %s", name)
	}

	return foundNg[0], nil
}

// GetAllNodeGroupNames collects and returns names for both managed and unmanaged nodegroups
func (c *ClusterConfig) GetAllNodeGroupNames() []string {
	var ngNames []string
	for _, ng := range c.NodeGroups {
		ngNames = append(ngNames, ng.NameString())
	}
	for _, ng := range c.ManagedNodeGroups {
		ngNames = append(ngNames, ng.NameString())
	}
	return ngNames
}

// AllNodeGroups combines managed and self-managed nodegroups and returns a slice of *api.NodeGroupBase containing
// both types of nodegroups
func (c *ClusterConfig) AllNodeGroups() []*NodeGroupBase {
	var baseNodeGroups []*NodeGroupBase
	for _, ng := range c.NodeGroups {
		baseNodeGroups = append(baseNodeGroups, ng.NodeGroupBase)
	}
	for _, ng := range c.ManagedNodeGroups {
		baseNodeGroups = append(baseNodeGroups, ng.NodeGroupBase)
	}
	return baseNodeGroups
}
