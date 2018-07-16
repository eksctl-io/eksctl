package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kubicorn/kubicorn/pkg/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mocks"
)

var _ = Describe("Eks", func() {
	var (
		cp     *ClusterProvider
		config ClusterConfig
	)

	BeforeEach(func() {

	})

	Describe("ListAll", func() {
		Context("With a cluster name", func() {
			var (
				clusterName *string
				err         error
			)

			BeforeEach(func() {
				clusterName = aws.String("test-cluster")

				config = ClusterConfig{
					ClusterName: *clusterName,
				}
				cp = &ClusterProvider{
					Cfg: &config,
					Svc: &ProviderServices{
						CFN: &mocks.CloudFormationAPI{},
						EKS: &mocks.EKSAPI{},
						EC2: &mocks.EC2API{},
						STS: &mocks.STSAPI{},
					},
				}

				cp.Svc.EKS.(*mocks.EKSAPI).On("DescribeCluster", mock.MatchedBy(func(input *eks.DescribeClusterInput) bool {
					return *input.Name == *clusterName
				})).Return(&eks.DescribeClusterOutput{
					Cluster: &eks.Cluster{
						Name:   clusterName,
						Status: aws.String(eks.ClusterStatusActive),
					},
				}, nil)
			})

			Context("and normal log level", func() {
				BeforeEach(func() {
					logger.Level = 3
				})

				JustBeforeEach(func() {
					err = cp.ListClusters()
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EKS service once", func() {
					Expect(cp.Svc.EKS.(*mocks.EKSAPI).AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)).To(BeTrue())
				})

				It("should not call AWS CFN ListStackPages", func() {
					Expect(cp.Svc.CFN.(*mocks.CloudFormationAPI).AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 0)).To(BeTrue())
				})
			})

			Context("and debug log level", func() {
				var (
					expectedStatusFilter string
				)
				BeforeEach(func() {
					expectedStatusFilter = "CREATE_COMPLETE"

					logger.Level = 4

					cp.Svc.CFN.(*mocks.CloudFormationAPI).On("ListStacksPages", mock.MatchedBy(func(input *cfn.ListStacksInput) bool {
						return *input.StackStatusFilter[0] == expectedStatusFilter
					}), mock.Anything).Return(nil)
				})

				JustBeforeEach(func() {
					err = cp.ListClusters()
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EKS service once", func() {
					Expect(cp.Svc.EKS.(*mocks.EKSAPI).AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)).To(BeTrue())
				})

				It("should have called AWS CFN ListStackPages", func() {
					Expect(cp.Svc.CFN.(*mocks.CloudFormationAPI).AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 1)).To(BeTrue())
				})
			})
		})
	})
})
