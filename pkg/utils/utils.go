package utils

import (
	"strings"

	kopsutils "k8s.io/kops/upup/pkg/fi/utils"
)

// IsGPUInstanceType returns tru of the instance type is GPU
// optimised.
func IsGPUInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p2") || strings.HasPrefix(instanceType, "p3")
}

// ExpandPath expands path with ~ notation
func ExpandPath(p string) string { return kopsutils.ExpandPath(p) }
