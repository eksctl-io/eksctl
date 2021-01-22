package nodegroup_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNodegroup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nodegroup Suite")
}
