package cloudconfig_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCFNBuilder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cloud-config Suite")
}
