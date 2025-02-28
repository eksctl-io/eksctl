package version

import (
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
)

// ParseEksctlVersion parses the an eksctl version as semver while ignoring
// extra build metadata
func ParseEksctlVersion(raw string) (semver.Version, error) {
	// We don't want any extra info from the version
	semverVersion := strings.Split(raw, ExtraSep)[0]
	v, err := semver.ParseTolerant(semverVersion)
	if err != nil {
		return v, fmt.Errorf("unexpected error parsing eksctl version %q: %w", raw, err)
	}
	return v, nil
}
