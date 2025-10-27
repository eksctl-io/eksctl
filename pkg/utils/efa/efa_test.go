package efa_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/efa"
)

var _ = Describe("EFA Utils", func() {
	Describe("IsBuiltInSupported", func() {
		DescribeTable("EFA built-in support version detection",
			func(kubernetesVersion string, expected bool, expectError bool) {
				result, err := efa.IsBuiltInSupported(kubernetesVersion, true)
				if expectError {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).NotTo(HaveOccurred())
					Expect(result).To(Equal(expected))
				}
			},
			Entry("Kubernetes 1.33 (exact match)", "1.33", true, false),
			Entry("Kubernetes 1.34 (higher version)", "1.34", true, false),
			Entry("Kubernetes 1.35 (higher version)", "1.35", true, false),
			Entry("Kubernetes 1.32 (lower version)", "1.32", false, false),
			Entry("Kubernetes 1.31 (lower version)", "1.31", false, false),
			Entry("Kubernetes 1.30 (lower version)", "1.30", false, false),
			Entry("Kubernetes 1.33.0 (with patch version)", "1.33.0", true, false),
			Entry("Kubernetes 1.33.1 (with patch version)", "1.33.1", true, false),
			Entry("Kubernetes 1.32.5 (lower with patch)", "1.32.5", false, false),
			Entry("Kubernetes 1.34.0 (higher with patch)", "1.34.0", true, false),
			Entry("invalid version string", "invalid", false, true),
			Entry("empty version string", "", false, true),
			Entry("version with v prefix", "v1.33", true, false),
			Entry("version with v prefix lower", "v1.32", false, false),
		)

		It("should use the correct EFA built-in support version constant", func() {
			Expect(api.EFABuiltInSupportVersion).To(Equal("1.33"))
		})

		Context("Test error handling", func() {
			It("should provide detailed error context for invalid versions", func() {
				_, err := efa.IsBuiltInSupported("invalid.version", true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to determine EFA built-in support"))
				Expect(err.Error()).To(ContainSubstring("invalid.version"))
				Expect(err.Error()).To(ContainSubstring("minimum required: 1.33"))
			})

			It("should return false if managed node groups is false", func() {
				result, err := efa.IsBuiltInSupported("1.34", false)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(BeFalse())
			})

			It("should provide detailed error context for malformed versions", func() {
				_, err := efa.IsBuiltInSupported("1.33.invalid.extra", true)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to determine EFA built-in support"))
				Expect(err.Error()).To(ContainSubstring("1.33.invalid.extra"))
			})
		})
	})
})
