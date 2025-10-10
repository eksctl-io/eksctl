package efa_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEFA(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EFA Suite")
}
