package testutils

import (
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// RegisterAndRun sets up and runs Ginkgo tests
func RegisterAndRun(t *testing.T) {
	_, suitePath, _, _ := runtime.Caller(1)
	RegisterFailHandler(Fail)
	RunSpecs(t, suitePath)
}
