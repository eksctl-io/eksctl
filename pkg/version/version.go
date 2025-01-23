package version

import (
	"encoding/json"
	"fmt"
	"strings"
)

//go:generate go run ./release_generate.go

// Info holds version information
type Info struct {
	Version      string
	PreReleaseID string
	Metadata     BuildMetadata
}

// BuildMetadata contains the semver build metadata:
// short commit hash and date in format YYYYMMDDTHHmmSS
type BuildMetadata struct {
	BuildDate string
	GitCommit string
}

// GetVersionInfo returns version Info struct
func GetVersionInfo() Info {
	return Info{
		Version:      Version,
		PreReleaseID: PreReleaseID,
		Metadata: BuildMetadata{
			GitCommit: gitCommit,
			BuildDate: buildDate,
		},
	}
}

// ExtraSep separates semver version from any extra version info
const ExtraSep = "-"

// String return version info as JSON
func String() string {
	if data, err := json.Marshal(GetVersionInfo()); err == nil {
		return string(data)
	}
	return ""
}

// GetVersion return the exact version of this build
func GetVersion() string {
	if PreReleaseID == "" {
		return Version
	}

	versionWithPR := fmt.Sprintf("%s%s%s", Version, ExtraSep, PreReleaseID)

	if isReleaseCandidate(PreReleaseID) || (gitCommit == "" || buildDate == "") {
		return versionWithPR
	}

	//  Include build metadata
	return fmt.Sprintf("%s+%s.%s",
		versionWithPR,
		gitCommit,
		buildDate,
	)
}

func isReleaseCandidate(preReleaseID string) bool {
	return strings.HasPrefix(preReleaseID, "rc.")
}
