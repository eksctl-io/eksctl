package anywhere_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAnywhere(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Anywhere Suite")
}
