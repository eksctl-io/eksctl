package version

var (
	// Version is the version number in semver format X.Y.Z
	Version = "0.13.0"

	// PreReleaseId can be empty for releases, "rc.X" for release candidates and "dev" for snapshots
	PreReleaseId = "dev"

	// gitCommit is the short commit hash
	gitCommit = ""

	// buildDate is the time of the build with format yyyymmddThhmmss
	buildDate = ""
)
