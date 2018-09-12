package ami_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAmi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ami Suite")
}
