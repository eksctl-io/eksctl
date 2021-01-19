package identitymapping_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIdentitymapping(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Identitymapping Suite")
}
