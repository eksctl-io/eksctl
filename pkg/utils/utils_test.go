package utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/utils"
)

var _ = Describe("Utils", func() {
	Describe("ToKebabCase", func() {
		DescribeTable("converts CamelCase to kebab-case",
			func(input, expected string) {
				result := utils.ToKebabCase(input)
				Expect(result).To(Equal(expected))
			},
			Entry("simple case", "CamelCase", "camel-case"),
			Entry("with numbers", "Test123Case", "test-123-case"),
			Entry("single word", "Test", "test"),
			Entry("already lowercase", "test", "test"),
			Entry("with consecutive caps", "XMLHttpRequest", "x-m-l-http-request"),
		)
	})

	Describe("IsMinVersion (wrapper)", func() {
		It("should delegate to version.IsMinVersion", func() {
			result, err := utils.IsMinVersion("1.33", "1.34")
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeTrue())
		})
	})

	Describe("CompareVersions (wrapper)", func() {
		It("should delegate to version.CompareVersions", func() {
			result, err := utils.CompareVersions("1.34", "1.33")
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(1))
		})
	})

	Describe("FnvHash", func() {
		It("should generate consistent hash for same input", func() {
			input := "test-string"
			hash1 := utils.FnvHash(input)
			hash2 := utils.FnvHash(input)
			Expect(hash1).To(Equal(hash2))
		})

		It("should generate different hashes for different inputs", func() {
			hash1 := utils.FnvHash("string1")
			hash2 := utils.FnvHash("string2")
			Expect(hash1).NotTo(Equal(hash2))
		})

		It("should return non-empty hash", func() {
			hash := utils.FnvHash("test")
			Expect(len(hash)).To(BeNumerically(">", 0))
		})
	})
})
