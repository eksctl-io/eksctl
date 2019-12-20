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
	PreReleaseId string
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
		Version:      version,
		PreReleaseId: preReleaseId,
		Metadata: BuildMetadata{
			GitCommit: gitCommit,
			BuildDate: buildDate,
		},
	}
}

// String return version info as JSON
func String() string {
	if data, err := json.Marshal(GetVersionInfo()); err == nil {
		return string(data)
	}
	return ""
}

// GetVersion return the exact version of this build
func GetVersion() string {
	if preReleaseId == "" {
		return version
	}

	if isReleaseCandidate(preReleaseId) {
		return fmt.Sprintf("%s-%s", version, preReleaseId)
	}

	//  Include build metadata
	return fmt.Sprintf("%s-%s+%s.%s",
		version,
		preReleaseId,
		gitCommit,
		buildDate,
	)
}

func isReleaseCandidate(preReleaseId string) bool {
	return strings.HasPrefix(preReleaseId, "rc.")
}
