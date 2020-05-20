package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver"
)

// IsGPUInstanceType returns true if the instance type is GPU optimised
func IsGPUInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "p2") || strings.HasPrefix(instanceType, "p3") || strings.HasPrefix(instanceType, "g3") || strings.HasPrefix(instanceType, "g4") || strings.HasPrefix(instanceType, "inf1")
}

// IsInferentiaInstanceType returns true if the instance type requires AWS Neuron
func IsInferentiaInstanceType(instanceType string) bool {
	return strings.HasPrefix(instanceType, "inf1")
}

var matchFirstCap = regexp.MustCompile("([0-9]+|[A-Z])")

// ToKebabCase turns a CamelCase string into a kebab-case string
func ToKebabCase(str string) string {
	kebab := matchFirstCap.ReplaceAllString(str, "-${1}")
	kebab = strings.TrimPrefix(kebab, "-")
	return strings.ToLower(kebab)
}

// IsMinVersion compares a given version number with a minimum one and returns true if
// version >= minimumVersion
func IsMinVersion(minimumVersion, version string) (bool, error) {
	minVersion, err := semver.ParseTolerant(minimumVersion)
	if err != nil {
		return false, fmt.Errorf("unable to parse minimum version required %s", minVersion)
	}
	targetVersion, err := semver.ParseTolerant(version)
	if err != nil {
		return false, fmt.Errorf("unable to parse version %s", version)
	}
	return targetVersion.GE(minVersion), nil
}
