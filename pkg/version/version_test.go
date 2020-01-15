package version

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("release tests", func() {
	BeforeEach(func() {
		Version = "0.5.0"
		PreReleaseID = ""
		gitCommit = ""
		buildDate = ""
	})

	It("ignores pre-release and build metadata for releases", func() {
		v := GetVersion()
		info := GetVersionInfo()

		Expect(v).To(Equal("0.5.0"))
		Expect(info).To(Equal(Info{
			Version:      "0.5.0",
			PreReleaseID: "",
			Metadata:     BuildMetadata{},
		}))
	})

	It("ignores build metadata for release candidates", func() {
		Version = "0.25.0"
		PreReleaseID = "rc.2"
		gitCommit = "abc123"
		buildDate = "today"

		v := GetVersion()
		info := GetVersionInfo()

		Expect(v).To(Equal("0.25.0-rc.2"))
		Expect(info).To(Equal(Info{
			Version:      "0.25.0",
			PreReleaseID: "rc.2",
			Metadata: BuildMetadata{
				GitCommit: "abc123",
				BuildDate: "today",
			},
		}))
	})

	It("wrong rc tag is treated like a dev version with metadata", func() {
		PreReleaseID = "rc1.2"
		gitCommit = "abc123"
		buildDate = "today"

		v := GetVersion()
		info := GetVersionInfo()

		Expect(v).To(Equal("0.5.0-rc1.2+abc123.today"))
		Expect(info).To(Equal(Info{
			Version:      "0.5.0",
			PreReleaseID: "rc1.2",
			Metadata: BuildMetadata{
				GitCommit: "abc123",
				BuildDate: "today",
			},
		}))
	})

	It("produces a dev version with build metadata", func() {
		PreReleaseID = "dev"
		gitCommit = "abc1234"
		buildDate = "2020-01-15T14:03:46Z"

		v := GetVersion()
		info := GetVersionInfo()

		Expect(v).To(Equal("0.5.0-dev+abc1234.2020-01-15T14:03:46Z"))
		Expect(info).To(Equal(Info{
			Version:      "0.5.0",
			PreReleaseID: "dev",
			Metadata: BuildMetadata{
				GitCommit: "abc1234",
				BuildDate: "2020-01-15T14:03:46Z",
			},
		}))
	})

	It("skips build metadata when not present", func() {
		PreReleaseID = "dev"

		v := GetVersion()
		info := GetVersionInfo()

		Expect(v).To(Equal("0.5.0-dev"))
		Expect(info).To(Equal(Info{
			Version:      "0.5.0",
			PreReleaseID: "dev",
			Metadata:     BuildMetadata{},
		}))
	})

})
