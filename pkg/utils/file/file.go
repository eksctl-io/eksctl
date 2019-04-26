package file

import (
	kopsutils "k8s.io/kops/upup/pkg/fi/utils"
	"os"
)

// Exists checks to see if a file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// ExpandPath expands path with ~ notation
func ExpandPath(p string) string { return kopsutils.ExpandPath(p) }
