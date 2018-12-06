package version

import (
	"encoding/json"
)

//go:generate go run ./release_generate.go

// Info hold version information
type Info struct {
	BuiltAt   string
	GitCommit string `json:",omitempty"`
	GitTag    string `json:",omitempty"`
}

var info = Info{
	BuiltAt:   builtAt,
	GitCommit: gitCommit,
	GitTag:    gitTag,
}

// Get return version Info struct
func Get() Info { return info }

// String return version info as JSON
func String() string {
	if data, err := json.Marshal(info); err == nil {
		return string(data)
	}
	return ""
}
