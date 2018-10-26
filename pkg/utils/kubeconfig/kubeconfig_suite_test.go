package kubeconfig

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestKubeConfig(t *testing.T) {
	testutils.RegisterAndRun(t, "KubeConfig Suite")
}
