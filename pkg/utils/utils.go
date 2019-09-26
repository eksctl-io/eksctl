package utils

import (
	"strings"
)

// IsGPUInstanceType returns tru of the instance type is GPU
// optimised.
func IsGPUInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p2") || strings.HasPrefix(instanceType, "p3") || strings.HasPrefix(instanceType, "g3") || strings.HasPrefix(instanceType, "g4")
}

// HasGPUInstanceType returns true if it finds a gpu instance among the mixed instances
func HasGPUInstanceType(instanceTypes []string) bool {
	for _, instanceType := range instanceTypes {
		if IsGPUInstanceType(instanceType) {
			return true
		}
	}
	return false
}
