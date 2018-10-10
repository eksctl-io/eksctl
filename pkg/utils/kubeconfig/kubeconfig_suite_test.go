package kubeconfig

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKubeConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KubeConfig Suite")
}
