package kubeconfig

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}
