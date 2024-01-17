package accessentry_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAccessEntry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Access Entry Suite")
}
