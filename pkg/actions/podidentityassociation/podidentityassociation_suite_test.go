package podidentityassociation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPodIdentityAssociation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nodegroup Suite")
}
