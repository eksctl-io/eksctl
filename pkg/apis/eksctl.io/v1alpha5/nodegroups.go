package v1alpha5

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

// HasNodegroup returns true if this clusterConfig contains a managed or un-managed nodegroup with the given name
func (c *ClusterConfig) FindNodegroup(name string) *NodeGroup {
	for _, ng := range c.NodeGroups {
		if name == ng.NameString() {
			return ng
		}
	}
	return nil
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
