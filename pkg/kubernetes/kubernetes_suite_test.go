package kubernetes_test

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestKubernetes(t *testing.T) {
	testutils.RegisterAndRun(t)
}
