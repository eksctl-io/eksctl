package file

import (
	"os"

	kopsutils "k8s.io/kops/upup/pkg/fi/utils"
)

// Exists checks to see if a file exists.
func Exists(path string) bool {
	extendedPath := ExpandPath(path)
	_, err := os.Stat(extendedPath)
	return err == nil
}

// ExpandPath expands path with ~ notation
func ExpandPath(p string) string { return kopsutils.ExpandPath(p) }
