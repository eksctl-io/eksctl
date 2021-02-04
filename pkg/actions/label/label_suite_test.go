package label_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLabel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Label Suite")
}
