package apply_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestApply(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Apply Suite")
}
