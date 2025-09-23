package version_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/utils/version"
)

var _ = Describe("Version", func() {
	Describe("IsMinVersion", func() {
		DescribeTable("version comparison",
			func(minimumVersion, versionString string, expected bool, expectError bool) {
				result, err := version.IsMinVersion(minimumVersion, versionString)
				if expectError {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal(expected))
				}
			},
			Entry("equal versions", "1.33", "1.33", true, false),
			Entry("higher version", "1.33", "1.34", true, false),
			Entry("lower version", "1.33", "1.32", false, false),
			Entry("patch version higher", "1.33.0", "1.33.1", true, false),
			Entry("patch version lower", "1.33.1", "1.33.0", false, false),
			Entry("invalid minimum version", "invalid", "1.33", false, true),
			Entry("invalid target version", "1.33", "invalid", false, true),
			Entry("version with v prefix", "1.33", "v1.33", true, false),
			Entry("version with v prefix lower", "1.33", "v1.32", false, false),
		)
	})

	Describe("CompareVersions", func() {
		DescribeTable("version comparison",
			func(a, b string, expected int, expectError bool) {
				result, err := version.CompareVersions(a, b)
				if expectError {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal(expected))
				}
			},
			Entry("equal versions", "1.33", "1.33", 0, false),
			Entry("first version higher", "1.34", "1.33", 1, false),
			Entry("first version lower", "1.32", "1.33", -1, false),
			Entry("invalid first version", "invalid", "1.33", 0, true),
			Entry("invalid second version", "1.33", "invalid", 0, true),
		)
	})

})
