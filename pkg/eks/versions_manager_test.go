package eks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils"
)

var _ = Describe("upgrade cluster", func() {
	type upgradeCase struct {
		givenVersion           string
		eksVersion             string
		expectedUpgradeVersion string
		expectedErrorText      string
	}

	DescribeTable("checks the specified version",
		func(c upgradeCase) {
			mockProvider := mockprovider.NewMockProvider()
			mockProvider.MockEKS().On("DescribeClusterVersions", mock.Anything, &awseks.DescribeClusterVersionsInput{}).
				Return(&awseks.DescribeClusterVersionsOutput{
					ClusterVersions: []ekstypes.ClusterVersionInformation{
						{
							ClusterVersion: utils.StringPtr(api.Version1_27),
						},
						{
							ClusterVersion: utils.StringPtr(api.Version1_28),
						},
						{
							ClusterVersion: utils.StringPtr(api.Version1_29),
						},
						{
							ClusterVersion: utils.StringPtr(api.Version1_30),
						},
						{
							ClusterVersion: utils.StringPtr(api.Version1_31),
							DefaultVersion: true,
						},
					},
				}, nil)

			cvm, err := eks.NewClusterVersionsManager(mockProvider.EKS())
			Expect(err).NotTo(HaveOccurred())

			upgradeVersion, err := cvm.ResolveUpgradeVersion(c.givenVersion, c.eksVersion)

			if c.expectedErrorText != "" {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(c.expectedErrorText))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(upgradeVersion).To(Equal(c.expectedUpgradeVersion))
			}
		},

		Entry("upgrades by default when the version is not specified", upgradeCase{
			givenVersion:           "",
			eksVersion:             api.Version1_27,
			expectedUpgradeVersion: api.Version1_28,
		}),

		Entry("upgrades by default when the version is auto", upgradeCase{
			givenVersion:           "auto",
			eksVersion:             api.Version1_27,
			expectedUpgradeVersion: api.Version1_28,
		}),

		Entry("does not upgrade or fail when the cluster is already in the last version", upgradeCase{
			givenVersion:           "",
			eksVersion:             api.Version1_31,
			expectedUpgradeVersion: "",
		}),

		Entry("upgrades to the next version when specified", upgradeCase{
			givenVersion:           api.Version1_28,
			eksVersion:             api.Version1_27,
			expectedUpgradeVersion: api.Version1_28,
		}),

		Entry("does not upgrade when the current version is specified", upgradeCase{
			givenVersion:           api.Version1_30,
			eksVersion:             api.Version1_30,
			expectedUpgradeVersion: "",
		}),

		Entry("fails when the upgrade jumps more than one kubernetes version", upgradeCase{
			givenVersion:      api.Version1_31,
			eksVersion:        api.Version1_29,
			expectedErrorText: "upgrading more than one version at a time is not supported",
		}),

		Entry("downgrades to the previous version when specified", upgradeCase{
			givenVersion:           api.Version1_29,
			eksVersion:             api.Version1_30,
			expectedUpgradeVersion: api.Version1_29,
		}),

		Entry("fails when the downgrade jumps more than one kubernetes version", upgradeCase{
			givenVersion:      api.Version1_28,
			eksVersion:        api.Version1_30,
			expectedErrorText: "downgrading more than one version at a time is not supported",
		}),

		Entry("fails when the version is deprecated", upgradeCase{
			givenVersion:      api.Version1_12,
			eksVersion:        api.Version1_12,
			expectedErrorText: "control plane version \"1.12\" has been deprecated",
		}),

		Entry("fails when the version is not supported", upgradeCase{
			givenVersion:      "1.50",
			eksVersion:        api.Version1_31,
			expectedErrorText: "control plane version \"1.50\" is not supported",
		}),
	)
})

var _ = Describe("resolve cluster version", func() {
	type resolveCase struct {
		eksDefaultVersion      string
		givenVersion           string
		expectedDefaultVersion string
		expectedVersion        string
		expectedErrorText      string
	}

	DescribeTable("resolves the specified version",
		func(c resolveCase) {
			var clusterVersions []ekstypes.ClusterVersionInformation
			for _, v := range []string{"1.33", "1.34", "1.35", "1.36", "1.37"} {
				clusterVersions = append(clusterVersions, ekstypes.ClusterVersionInformation{
					ClusterVersion: utils.StringPtr(v),
					DefaultVersion: v == c.eksDefaultVersion,
				})
			}
			mockProvider := mockprovider.NewMockProvider()
			mockProvider.MockEKS().On("DescribeClusterVersions", mock.Anything, &awseks.DescribeClusterVersionsInput{}).
				Return(&awseks.DescribeClusterVersionsOutput{
					ClusterVersions: clusterVersions,
				}, nil)

			cvm, err := eks.NewClusterVersionsManager(mockProvider.EKS())
			Expect(err).NotTo(HaveOccurred())
			Expect(cvm.DefaultVersion()).To(Equal(c.expectedDefaultVersion))

			version, err := cvm.ResolveClusterVersion(c.givenVersion)

			if c.expectedErrorText != "" {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(c.expectedErrorText))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(version).To(Equal(c.expectedVersion))
			}
		},

		Entry("uses the EKS default version when the version is not specified", resolveCase{
			eksDefaultVersion:      "1.36",
			givenVersion:           "",
			expectedDefaultVersion: "1.36",
			expectedVersion:        "1.36",
		}),

		Entry("uses the EKS default version when the version is auto", resolveCase{
			eksDefaultVersion:      "1.36",
			givenVersion:           "auto",
			expectedDefaultVersion: "1.36",
			expectedVersion:        "1.36",
		}),

		Entry("uses the latest version when the version is latest", resolveCase{
			eksDefaultVersion:      "1.36",
			givenVersion:           "latest",
			expectedDefaultVersion: "1.36",
			expectedVersion:        "1.37",
		}),

		Entry("uses the specified version when it is supported", resolveCase{
			eksDefaultVersion:      "1.36",
			givenVersion:           "1.33",
			expectedDefaultVersion: "1.36",
			expectedVersion:        "1.33",
		}),

		Entry("fails when the version is not specified and EKS does not flag a default", resolveCase{
			eksDefaultVersion: "",
			givenVersion:      "",
			expectedErrorText: "EKS did not report a default Kubernetes version",
		}),

		Entry("fails when the version is auto and EKS does not flag a default", resolveCase{
			eksDefaultVersion: "",
			givenVersion:      "auto",
			expectedErrorText: "EKS did not report a default Kubernetes version",
		}),

		Entry("uses the specified version when EKS does not flag a default", resolveCase{
			eksDefaultVersion: "",
			givenVersion:      "1.35",
			expectedVersion:   "1.35",
		}),
	)
})
