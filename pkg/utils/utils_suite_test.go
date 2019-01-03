package utils_test

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestUtils(t *testing.T) {
	testutils.RegisterAndRun(t, "utils Suite")
}
