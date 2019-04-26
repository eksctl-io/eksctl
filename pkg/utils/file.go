package utils

import (
	kopsutils "k8s.io/kops/upup/pkg/fi/utils"
	"os"
)

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

// CheckFileExists returns an error if the file doesn't exist
func CheckFileExists(filePath string) error {
	extendedPath := ExpandPath(filePath)
	if _, err := os.Stat(extendedPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

// ExpandPath expands path with ~ notation
func ExpandPath(p string) string { return kopsutils.ExpandPath(p) }
