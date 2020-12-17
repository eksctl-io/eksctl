package assetutil

import (
	"fmt"
)

// Error describes an asset error
type Error struct {
	error
}

func (ae *Error) Error() string {
	return fmt.Sprintf("unexpected error generating assets: %v", ae.error.Error())
}

// Func describes an asset loading function
type Func func() ([]byte, error)

// MustLoad loads an asset or panics if the asset couldn't be loaded
func MustLoad(assetFunc Func) []byte {
	bytes, err := assetFunc()
	if err != nil {
		panic(&Error{err})
	}
	return bytes
}
