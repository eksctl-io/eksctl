package v1alpha5_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Resource Account IDs", func() {
	DescribeTable("when 'ISO_EKS_ACCOUNT_ID' environment variable is set", func(region string) {
		Expect(os.Setenv(api.IsoEKSAccountIDEnv, "1234567890")).To(Succeed())

		Expect(api.ValidateEKSAccountID(region)).To(Succeed())
	},
		Entry("using RegionUSIsoEast1 should not produce an error", api.RegionUSIsoEast1),
		Entry("using RegionUSIsobEast1 should not produce an error", api.RegionUSIsobEast1),
	)

	DescribeTable("when 'ISO_EKS_ACCOUNT_ID' environment variable IS present but is NOT set", func(region string) {
		Expect(os.Setenv(api.IsoEKSAccountIDEnv, "")).To(Succeed())

		Expect(api.ValidateEKSAccountID(region)).To(MatchError(ContainSubstring("ISO_EKS_ACCOUNT_ID not set, required for use of region:")))
	},
		Entry("setting RegionUSIsoEast1 should fail", api.RegionUSIsoEast1),
		Entry("setting RegionUSIsobEast1 should fail", api.RegionUSIsobEast1),
	)

	DescribeTable("when 'ISO_EKS_ACCOUNT_ID' environment variable is NOT set", func(region string) {
		Expect(os.Unsetenv(api.IsoEKSAccountIDEnv)).To(Succeed())

		Expect(api.ValidateEKSAccountID(region)).To(MatchError(ContainSubstring("ISO_EKS_ACCOUNT_ID not set, required for use of region:")))
	},
		Entry("setting RegionUSIsoEast1 should fail", api.RegionUSIsoEast1),
		Entry("setting RegionUSIsobEast1 should fail", api.RegionUSIsobEast1),
	)
})
