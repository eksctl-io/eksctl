package info

import (
	"encoding/json"
	"strings"

	"github.com/Masterminds/semver/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Info", func() {
	Context("GetInfo", func() {
		It("returns the eksctl and kubectl versions, and the host os", func() {
			result := GetInfo()

			Expect(result.EksctlVersion).NotTo(Equal(""))
			Expect(result.KubectlVersion).NotTo(Equal(""))
			Expect(result.OS).NotTo(Equal(""))

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

			Expect(infos.EksctlVersion).NotTo(Equal(""))
			Expect(infos.KubectlVersion).NotTo(Equal(""))
			Expect(infos.OS).NotTo(Equal(""))
		})
	})

})
