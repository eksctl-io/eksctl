package nodegroup

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
)

// ParseReleaseVersion parses an AMI release version string that's in the format `1.18.8-20201007`
func ParseReleaseVersion(releaseVersion string) (AMIReleaseVersion, error) {
	parts := strings.Split(releaseVersion, "-")
	if len(parts) != 2 {
		return AMIReleaseVersion{}, fmt.Errorf("unexpected format for release version: %q", releaseVersion)
	}
	v, err := semver.ParseTolerant(parts[0])
	if err != nil {
		return AMIReleaseVersion{}, fmt.Errorf("invalid SemVer version: %w", err)
	}
	return AMIReleaseVersion{
		Version: v,
		Date:    parts[1],
	}, nil
}

type AMIReleaseVersion struct {
	Version semver.Version
	Date    string
}

// LTE checks if a is less than or equal to b.
func (a AMIReleaseVersion) LTE(b AMIReleaseVersion) bool {
	return a.Compare(b) <= 0
}

// GTE checks if a is greater than or equal to b.
func (a AMIReleaseVersion) GTE(b AMIReleaseVersion) bool {
	return a.Compare(b) >= 0
}

// Compare returns 0 if a==b, -1 if a < b, and +1 if a > b.
func (a AMIReleaseVersion) Compare(b AMIReleaseVersion) int {
	cmp := a.Version.Compare(b.Version)
	if cmp == 0 {
		return strings.Compare(a.Date, b.Date)
	}
	return cmp
}
