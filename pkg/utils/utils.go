package utils

import (
	"os"
	"strings"

	kopsutils "k8s.io/kops/upup/pkg/fi/utils"
)

// IsGPUInstanceType returns tru of the instance type is GPU
// optimised.
func IsGPUInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p2") || strings.HasPrefix(instanceType, "p3")
}

// FileExists checks to see if a file exists.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// ExpandPath expands path with ~ notation
func ExpandPath(p string) string { return kopsutils.ExpandPath(p) }
