package fargate_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFargate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fargate Suite")
}
