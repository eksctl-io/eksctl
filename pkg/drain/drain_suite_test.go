package drain_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDrain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Drain Suite")
}
