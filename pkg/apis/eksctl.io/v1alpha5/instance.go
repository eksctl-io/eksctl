package v1alpha5

import (
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
)

// SelectInstanceType determines which instanceType is relevant for selecting an AMI
// If the nodegroup has mixed instances it will prefer a GPU instance type over a general class one
// This is to make sure that the AMI that is selected later is valid for all the types
func SelectInstanceType(np NodePool) string {
	var instanceTypes []string
	switch ng := np.(type) {
	case *NodeGroup:
		if ng.InstancesDistribution != nil {
			instanceTypes = ng.InstancesDistribution.InstanceTypes
		}
	case *ManagedNodeGroup:
		instanceTypes = ng.InstanceTypes
	}

	hasMixedInstances := len(instanceTypes) > 0
	if hasMixedInstances {
		for _, instanceType := range instanceTypes {
			if instanceutils.IsGPUInstanceType(instanceType) {
				return instanceType
			}
		}
		return instanceTypes[0]
	}

	return np.BaseNodeGroup().InstanceType
}
