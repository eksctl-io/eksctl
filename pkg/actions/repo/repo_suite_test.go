package repo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGitops(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gitops Repo Suite")
}
