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
