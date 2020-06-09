package file

import (
	"os"
	"strings"

	"k8s.io/client-go/util/homedir"
)

// Exists checks to see if a file exists.
func Exists(path string) bool {
	extendedPath := ExpandPath(path)
	_, err := os.Stat(extendedPath)
	return err == nil
}

// ExpandPath expands path with ~ notation
func ExpandPath(p string) string {
	if strings.HasPrefix(p, "~/") {
		p = homedir.HomeDir() + p[1:]
	}

	return p
}
