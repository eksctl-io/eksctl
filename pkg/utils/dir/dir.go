package dir

import (
	"io"
	"os"

	"github.com/weaveworks/eksctl/pkg/utils/file"
)

// IsEmpty checks if the provided directory is empty or not.
func IsEmpty(path string) (bool, error) {
	path = file.ExpandPath(path)
	d, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer d.Close()
	// Check ONLY the next file's metadata, as no need to check more to know whether the directory is empty or not.
	_, err = d.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
