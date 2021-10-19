package instance

import (
	"strings"
)

// IsARMInstanceType returns true if the instance type is ARM architecture
func IsARMInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "a1") ||
		strings.HasPrefix(instanceType, "t4g") ||
		strings.HasPrefix(instanceType, "m6g") ||
		strings.HasPrefix(instanceType, "m6gd") ||
		strings.HasPrefix(instanceType, "c6g") ||
		strings.HasPrefix(instanceType, "c6gd") ||
		strings.HasPrefix(instanceType, "r6g") ||
		strings.HasPrefix(instanceType, "r6gd")
}

// IsGPUInstanceType returns true if the instance type is GPU optimised
func IsGPUInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p2") ||
		strings.HasPrefix(instanceType, "p3") ||
		strings.HasPrefix(instanceType, "p4") ||
		strings.HasPrefix(instanceType, "g3") ||
		strings.HasPrefix(instanceType, "g4") ||
		strings.HasPrefix(instanceType, "inf1")
}

// IsInferentiaInstanceType returns true if the instance type requires AWS Neuron
func IsInferentiaInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "inf1")
}
