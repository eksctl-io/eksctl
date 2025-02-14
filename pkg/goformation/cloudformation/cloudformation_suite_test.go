package cloudformation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCloudformation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cloudformation Suite")
}
