package info

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/weaveworks/eksctl/pkg/utils/kubectl"
	"github.com/weaveworks/eksctl/pkg/version"
)

//go:generate go run ./release_generate.go

// Info holds versions info
type Info struct {
	EksctlVersion  string
	KubectlVersion string
	OS             string
}

// GetInfo returns versions info
func GetInfo() Info {
	return Info{
		EksctlVersion:  getEksctlVersion(),
		KubectlVersion: getKubectlVersion(),
		OS:             runtime.GOOS,
	}
}

// getEksctlVersion returns the eksctl version
func getEksctlVersion() string {
	return version.GetVersion()
}

func getKubectlVersion() string {
	clientVersion, err := kubectl.NewVersionManager().ClientVersion()
	if err != nil {
		return err.Error()
	}
	return clientVersion
}

// String return info as JSON
func String() string {
	data, err := json.Marshal(GetInfo())
	if err != nil {
		return fmt.Sprintf("failed to marshal info into json: %q", err)
	}

	return string(data)
}
