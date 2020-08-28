package flux

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestGitopsFlux(t *testing.T) {
	testutils.RegisterAndRun(t)
}
