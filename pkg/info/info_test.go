package info

import (
	"encoding/json"
	"strings"

	"github.com/Masterminds/semver/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info", func() {
	Context("GetInfo", func() {
		It("returns the eksctl and kubectl versions, and the host os", func() {
			result := GetInfo()

			Expect(result.EksctlVersion).ToNot(Equal(""))
			Expect(result.KubectlVersion).ToNot(Equal(""))
			Expect(result.OS).ToNot(Equal(""))

			_, err := semver.NewVersion(strings.TrimSpace(result.EksctlVersion))
			Expect(err).NotTo(HaveOccurred())

			_, err = semver.NewVersion(strings.TrimSpace(strings.TrimPrefix(result.KubectlVersion, "v")))
			Expect(err).NotTo(HaveOccurred())

			oses := []string{"aix", "android", "darwin", "dragonfly", "freebsd", "hurd", "illumos", "ios", "js", "linux", "nacl", "netbsd", "openbsd", "plan9", "solaris", "windows", "zos"}
			Expect(result.OS).To(BeElementOf(oses))
		})
	})

	Context("String", func() {
		It("returns Info in json", func() {
			result := String()

			infos := Info{}
			err := json.Unmarshal([]byte(result), &infos)
			Expect(err).NotTo(HaveOccurred())

			Expect(infos.EksctlVersion).ToNot(Equal(""))
			Expect(infos.KubectlVersion).ToNot(Equal(""))
			Expect(infos.OS).ToNot(Equal(""))
		})
	})

})
