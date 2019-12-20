// +build !release

package version

var (
	// Version is the version number in semver format X.Y.Z
	Version = "0.12.0"

	// PreReleaseId can be empty for releases, "rc.X" for release candidates and "dev" for snapshots
	PreReleaseId = ""

	// gitCommit is the short commit hash
	gitCommit = ""

	// buildDate is the time of the build with format yyyymmddThhmmss
	buildDate = ""
)
