package v1alpha3

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestCFNManager(t *testing.T) {
	testutils.RegisterAndRun(t, "eks api Suite")
}
