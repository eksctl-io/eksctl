package gitops

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestGitops(t *testing.T) {
	RegisterFailHandler(Fail)
	testutils.RegisterAndRun(t)
}
