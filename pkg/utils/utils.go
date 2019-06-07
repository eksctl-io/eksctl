package utils

import (
	"strings"
)

// IsGPUInstanceType returns tru of the instance type is GPU
// optimised.
func IsGPUInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p2") || strings.HasPrefix(instanceType, "p3")
}

// HasGPUInstanceType returns true if it finds a gpu instance among the mixed instances
func HasGPUInstanceType(instanceTypes []string) bool {
	if instanceTypes == nil || len(instanceTypes) == 0 {
		return false
	}
	for _, instanceType := range instanceTypes {
		if IsGPUInstanceType(instanceType) {
			return true
		}
	}
	return false
}
