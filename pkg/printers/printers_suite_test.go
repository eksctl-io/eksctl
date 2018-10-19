package printers_test

import (
	"testing"

	"github.com/weaveworks/eksctl/pkg/testutils"
)

func TestPrinters(t *testing.T) {
	testutils.RegisterAndRun(t, "Printers Suite")
}
