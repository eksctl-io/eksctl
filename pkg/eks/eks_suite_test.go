package eks_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Eks Suite")
}
