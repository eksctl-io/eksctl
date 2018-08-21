package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/kubicorn/kubicorn/pkg/namer"
	kopsutils "k8s.io/kops/upup/pkg/fi/utils"
)

// ClusterName generates a neme string when a and b are empty strings.
// If either a or b are non-empty, it returns whichever is non-empty.
// If neither a nor b are empty, it returns empty name, to indicate
// ambigous usage.
func ClusterName(a, b string) string {
	if a != "" && b != "" {
		return ""
	}
	if a != "" {
		return a
	}
	if b != "" {
		return b
	}
	return fmt.Sprintf("%s-%d", namer.RandomName(), time.Now().Unix())
}

// FileExists checks to see if a file exists.
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// ExpandPath expands path with ~ notation
func ExpandPath(p string) string { return kopsutils.ExpandPath(p) }
