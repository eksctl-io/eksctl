package version

import (
	"github.com/blang/semver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseEksctlVersion", func() {
	It("handles versions with metadata", func() {
		gitVersion := "0.27.0-dev+001eeced.2020-08-27T03:03:31Z"

		v, err := ParseEksctlVersion(gitVersion)

		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal(
			semver.Version{
				Major: 0,
				Minor: 27,
				Patch: 0,
			},
		))
	})
	It("handles versions without metadata", func() {
		gitVersion := "0.27.0"

		v, err := ParseEksctlVersion(gitVersion)

		Expect(err).NotTo(HaveOccurred())
		Expect(v).To(Equal(
			semver.Version{
				Major: 0,
				Minor: 27,
				Patch: 0,
			},
		))
	})
})
