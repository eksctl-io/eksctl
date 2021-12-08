package karpenter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKarpenter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Karpenter Suite")
}
