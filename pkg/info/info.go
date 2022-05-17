package info

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/weaveworks/eksctl/pkg/version"
)

//go:generate go run ./release_generate.go

// Info holds versions info
type Info struct {
	EksctlVersion  string
	KubectlVersion string
	OS             string
}

// clientVersion holds git version info of kubectl client
type clientVersion struct {
	GitVersion string `json:"gitVersion"`
}

// kubectlInfo holds version info of kubectl client
type kubectlInfo struct {
	ClientVersion clientVersion `json:"clientVersion"`
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

// getKubectlVersion returns the kubectl version
func getKubectlVersion() string {
	cmd := exec.Command("kubectl", "version", "--client", "--output", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error : %v", err)
	}

	var info kubectlInfo

	if err := json.Unmarshal(out, &info); err != nil {
		return fmt.Sprintf("error parsing `kubectl version` output: %v", err)
	}

	if info.ClientVersion.GitVersion == "" {
		return "unknown version"
	}

	return info.ClientVersion.GitVersion
}

// String return info as JSON
func String() string {
	data, err := json.Marshal(GetInfo())
	if err != nil {
		return fmt.Sprintf("failed to marshal info into json: %q", err)
	}

	return string(data)
}
