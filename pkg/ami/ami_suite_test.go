package ami_test

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestAmi(t *testing.T) {
	testutils.RegisterAndRun(t, "Ami Suite")
}
