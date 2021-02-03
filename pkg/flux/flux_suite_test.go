package flux_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flux Suite")
}
