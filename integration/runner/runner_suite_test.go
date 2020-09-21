package runner

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestRunner(t *testing.T) {
	testutils.RegisterAndRun(t)
}
