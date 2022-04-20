package flux_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGitops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gitops Flux Suite")
}
