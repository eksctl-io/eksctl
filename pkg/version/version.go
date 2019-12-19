package version

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
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

// Get return version Info struct
func Get() Info {
	return Info{
		Version:      version,
		PreReleaseId: preReleaseId,
		Metadata: BuildMetadata{
			GitCommit: gitCommit,
			BuildDate: getBuildDate(),
		},
	}
}

// String return version info as JSON
func String() string {
	if data, err := json.Marshal(Get()); err == nil {
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
		getBuildDate(),
	)
}

func isReleaseCandidate(preReleaseId string) bool {
	return strings.HasPrefix(preReleaseId, "rc.")
}

func getBuildDate() string {
	return time.Now().UTC().Format("20060102T150405")
}
