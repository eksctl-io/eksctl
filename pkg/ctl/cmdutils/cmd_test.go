package cmdutils

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("resolveClusterVersion", func() {
	type resolveClusterVersionEntry struct {
		eksFlagsDefault bool
		version         string
		expectedVersion string
		expectedErr     string
	}

	DescribeTable("resolves the cluster version", func(e resolveClusterVersionEntry) {
		mockProvider := mockprovider.NewMockProvider()
		mockProvider.MockEKS().On("DescribeClusterVersions", mock.Anything, &awseks.DescribeClusterVersionsInput{}).
			Return(&awseks.DescribeClusterVersionsOutput{
				ClusterVersions: []ekstypes.ClusterVersionInformation{
					{
						ClusterVersion: aws.String(api.Version1_34),
						DefaultVersion: e.eksFlagsDefault,
					},
					{
						ClusterVersion: aws.String(api.Version1_35),
					},
				},
			}, nil)
		cvm, err := eks.NewClusterVersionsManager(mockProvider.EKS())
		Expect(err).NotTo(HaveOccurred())

		meta := &api.ClusterMeta{Version: e.version}
		err = resolveClusterVersion(cvm, meta)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(meta.Version).To(Equal(e.expectedVersion))
	},
		Entry("uses the EKS default version when the version is not set", resolveClusterVersionEntry{
			eksFlagsDefault: true,
			expectedVersion: api.Version1_34,
		}),
		Entry("leaves the version unset when EKS does not flag a default", resolveClusterVersionEntry{
			expectedVersion: "",
		}),
		Entry("uses the specified version when EKS does not flag a default", resolveClusterVersionEntry{
			version:         api.Version1_35,
			expectedVersion: api.Version1_35,
		}),
		Entry("fails for auto when EKS does not flag a default", resolveClusterVersionEntry{
			version:     "auto",
			expectedErr: "EKS did not report a default Kubernetes version",
		}),
	)
})
