package addon_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAddon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Addon Suite")
}
