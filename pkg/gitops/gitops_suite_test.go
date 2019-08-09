package gitops

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"testing"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	testutils.RegisterAndRun(t)
}
