package printers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPrinters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Printers Suite")
}
