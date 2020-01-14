package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/version"
)

var _ = Describe("release tests", func() {
	BeforeEach(func() {
		version.Version = "0.5.0"
		version.PreReleaseID = "dev"
	})

	It("produces a release without a pre-release id", func() {
		v, p := prepareRelease()

		Expect(v).To(Equal("0.5.0"))
		Expect(p).To(BeEmpty())
	})

	It("produces the correct release for 2 digit minor versions", func() {
		version.Version = "0.25.0"
		v, p := prepareRelease()

		Expect(v).To(Equal("0.25.0"))
		Expect(p).To(BeEmpty())
	})

	It("increases minor version for the next development iteration", func() {
		v, p := nextDevelopmentIteration()

		Expect(v).To(Equal("0.6.0"))
		Expect(p).To(Equal("dev"))
	})

	It("produces the correct default release candidate", func() {
		v, p := prepareReleaseCandidate("")

		Expect(v).To(Equal("0.5.0"))
		Expect(p).To(Equal("rc.0"))
	})

	It("allows overwritting the release candidate", func() {
		v, p := prepareReleaseCandidate("rc.3")

		Expect(v).To(Equal("0.5.0"))
		Expect(p).To(Equal("rc.3"))
	})

})
