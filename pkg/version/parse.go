package version

import (
	"strings"

	"github.com/blang/semver"
	"github.com/pkg/errors"
)

// ParseEksctlVersion parses the an eksctl version as semver while ignoring
// extra build metadata
func ParseEksctlVersion(raw string) (semver.Version, error) {
	// We don't want any extra info from the version
	semverVersion := strings.Split(raw, ExtraSep)[0]
	v, err := semver.ParseTolerant(semverVersion)
	if err != nil {
		return v, errors.Wrapf(err, "unexpected error parsing eksctl version %q", raw)
	}
	return v, nil
}
