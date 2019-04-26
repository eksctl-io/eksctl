package utils

import (
	"strings"
)

// IsGPUInstanceType returns tru of the instance type is GPU
// optimised.
func IsGPUInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p2") || strings.HasPrefix(instanceType, "p3")
}
