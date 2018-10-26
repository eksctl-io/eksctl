package az_test

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestAZ(t *testing.T) {
	testutils.RegisterAndRun(t, "AZ Suite")
}
