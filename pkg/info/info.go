package info

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"unicode"

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

// getKubectlVersion returns the kubectl version
func getKubectlVersion() string {
	cmd := exec.Command("kubectl", "version", "--short", "--client")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("error : %v", err)
	}

	values := strings.FieldsFunc(string(out), func(c rune) bool {
		return unicode.IsSpace(c)
	})
	if len(values) == 0 {
		return "unknown version"
	}
	return values[len(values)-1]
}

// String return info as JSON
func String() string {
	data, err := json.Marshal(GetInfo())
	if err != nil {
		return fmt.Sprintf("failed to marshal info into json: %q", err)
	}

	return string(data)
}
