package iam_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Iam Suite")
}
