package version

import (
	"fmt"

	"github.com/blang/semver/v4"
)

// IsMinVersion compares a given version number with a minimum one and returns true if
// version >= minimumVersion
func IsMinVersion(minimumVersion, version string) (bool, error) {
	minVersion, err := semver.ParseTolerant(minimumVersion)
	if err != nil {
		return false, fmt.Errorf("unable to parse minimum version required %s", minVersion)
	}
	targetVersion, err := semver.ParseTolerant(version)
	if err != nil {
		return false, fmt.Errorf("unable to parse target version %s", version)
	}
	return targetVersion.GE(minVersion), nil
}

// CompareVersions compares two version strings with the usual conventions:
// returns 0 if a == b
// returns 1 if a > b
// returns -1 if a < b
func CompareVersions(a, b string) (int, error) {
	aVersion, err := semver.ParseTolerant(a)
	if err != nil {
		return 0, fmt.Errorf("unable to parse first version %q: %w", a, err)
	}
	bVersion, err := semver.ParseTolerant(b)
	if err != nil {
		return 0, fmt.Errorf("unable to parse second version %q: %w", b, err)
	}
	return aVersion.Compare(bVersion), nil
}
