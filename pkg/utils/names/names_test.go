package names_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/names"
	"testing"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("name", func() {
	Describe("ForNodeGroup", func() {
		It("returns the first non-empty provided name if any", func() {
			first := names.ForNodeGroup("first-name", "")
			Expect(first).To(Equal("first-name"))
			second := names.ForNodeGroup("", "second-name")
			Expect(second).To(Equal("second-name"))
		})

		It("returns an empty string if both provided names are non-empty, so the client can test this and error-out", func() {
			name := names.ForNodeGroup("first-name", "second-name")
			Expect(name).To(Equal(""))
		})

		It("generates a random name otherwise", func() {
			name := names.ForNodeGroup("", "")
			Expect(name).To(MatchRegexp("ng-[abcdef0123456789]{8}"))
		})
	})

	Describe("ForFargateProfile", func() {
		It("returns the provided name if non-empty", func() {
			name := names.ForFargateProfile("my-favourite-name")
			Expect(name).To(Equal("my-favourite-name"))
		})
		It("generates a random name otherwise", func() {
			name := names.ForFargateProfile("")
			Expect(name).To(MatchRegexp("fp-[abcdef0123456789]{8}"))
		})
	})
})
