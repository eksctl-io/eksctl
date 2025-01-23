package eks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go/aws"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
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
							ClusterVersion: aws.String(api.Version1_27),
						},
						{
							ClusterVersion: aws.String(api.Version1_28),
						},
						{
							ClusterVersion: aws.String(api.Version1_29),
						},
						{
							ClusterVersion: aws.String(api.Version1_30),
						},
						{
							ClusterVersion: aws.String(api.Version1_31),
							DefaultVersion: true,
						},
					},
				}, nil)

			cvm, err := eks.NewClusterVersionsManager(mockProvider.EKS())
			Expect(err).NotTo(HaveOccurred())

			upgradeVersion, err := cvm.ResolveUpgradeVersion(c.givenVersion, c.eksVersion)

			if c.expectedErrorText != "" {
				if c.expectedErrorText != "cannot upgrade to a lower version" {
				} else {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(c.expectedErrorText))
				}
			} else {
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

		Entry("fails when the given version is lower than the current one", upgradeCase{
			givenVersion:      api.Version1_29,
			eksVersion:        api.Version1_30,
			expectedErrorText: "cannot upgrade to a lower version",
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
