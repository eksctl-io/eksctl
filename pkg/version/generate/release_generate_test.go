//go:build release
// +build release

package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/version"
)

var _ = Describe("release tests", func() {
	BeforeEach(func() {
		version.Version = "0.5.0"
		version.PreReleaseID = "dev"
	})

	It("increases minor version for the next development iteration from a release", func() {
		version.PreReleaseID = ""

		v, p := nextDevelopmentIteration()

		Expect(v).To(Equal("0.6.0"))
		Expect(p).To(Equal("dev"))
	})

	It("increases minor version for the next development iteration from an rc", func() {
		version.PreReleaseID = "rc.1"

		v, p := nextDevelopmentIteration()

		Expect(v).To(Equal("0.6.0"))
		Expect(p).To(Equal("dev"))
	})

	DescribeTable("gets the next pre-release ID", func(releaseVersion, expectedPreReleaseID string) {
		preReleaseID, err := nextPreReleaseID(releaseVersion)
		Expect(err).NotTo(HaveOccurred())
		Expect(preReleaseID).To(Equal(expectedPreReleaseID))
	},
		Entry("no previous RC", "", "rc.0"),
		Entry("rc.0 exists", "0.100.0-rc.0", "rc.1"),
		Entry("rc.5 exists", "v0.100.0-rc.5", "rc.6"),
	)

})
