package api

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestCFNManager(t *testing.T) {
	testutils.RegisterAndRun(t, "eks api Suite")
}
