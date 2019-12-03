package retry_test

import (
	"github.com/weaveworks/eksctl/pkg/testutils"
	"testing"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}