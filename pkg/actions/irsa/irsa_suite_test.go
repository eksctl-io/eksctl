package irsa_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIRSA(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IRSA Suite")
}
