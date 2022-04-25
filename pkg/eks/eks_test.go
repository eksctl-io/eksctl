package eks_test

import (
	"context"

	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("EKS API wrapper", func() {
	var (
		c *ClusterProvider
		p *mockprovider.MockProvider
	)

	Describe("GetClusters", func() {
		var (
			clusterName string
			err         error
			cluster     *awseks.Cluster
		)

		When("the cluster is ready", func() {
			BeforeEach(func() {
				clusterName = "test-cluster"

				p = mockprovider.NewMockProvider()

				c = &ClusterProvider{
					Provider: p,
				}

				p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
					return *input.Name == clusterName
				})).Return(&awseks.DescribeClusterOutput{
					Cluster: testutils.NewFakeCluster(clusterName, awseks.ClusterStatusActive),
				}, nil)
			})

			Context("and normal log level", func() {
				BeforeEach(func() {
					logger.Level = 3
				})

				JustBeforeEach(func() {
					cluster, err = c.GetCluster(context.Background(), clusterName)
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return the cluster", func() {
					Expect(cluster).To(Equal(testutils.NewFakeCluster(clusterName, awseks.ClusterStatusActive)))
				})

				It("should have called AWS EKS service once", func() {
					p.MockEKS().AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)
				})

				It("should not call AWS CFN ListStacks", func() {
					p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacks", 0)
				})
			})

			Context("and debug log level", func() {

				BeforeEach(func() {
					expectedStatusFilter := []string{
						"CREATE_IN_PROGRESS",
						"CREATE_FAILED",
						"CREATE_COMPLETE",
						"ROLLBACK_IN_PROGRESS",
						"ROLLBACK_FAILED",
						"ROLLBACK_COMPLETE",
						"DELETE_IN_PROGRESS",
						"DELETE_FAILED",
						"UPDATE_IN_PROGRESS",
						"UPDATE_COMPLETE_CLEANUP_IN_PROGRESS",
						"UPDATE_COMPLETE",
						"UPDATE_ROLLBACK_IN_PROGRESS",
						"UPDATE_ROLLBACK_FAILED",
						"UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS",
						"UPDATE_ROLLBACK_COMPLETE",
						"REVIEW_IN_PROGRESS",
					}

					logger.Level = 4

					p.MockCloudFormation().On("ListStacks", mock.Anything, mock.MatchedBy(func(input *cfn.ListStacksInput) bool {
						matches := 0
						for i := range input.StackStatusFilter {
							if input.StackStatusFilter[i] == types.StackStatus(expectedStatusFilter[i]) {
								matches++
							}
						}
						return matches == len(expectedStatusFilter)
					})).Return(&cfn.ListStacksOutput{}, nil)
				})

				JustBeforeEach(func() {
					cluster, err = c.GetCluster(context.Background(), clusterName)
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return the cluster", func() {
					Expect(cluster).To(Equal(testutils.NewFakeCluster(clusterName, awseks.ClusterStatusActive)))
				})

				It("should have called AWS EKS service once", func() {
					p.MockEKS().AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)
				})

				It("should have called AWS CFN ListStacks", func() {
					p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacks", 1)
				})
			})
		})

		When("the cluster is not ready", func() {
			BeforeEach(func() {
				clusterName = "test-cluster"
				logger.Level = 1

				p = mockprovider.NewMockProvider()

				c = &ClusterProvider{
					Provider: p,
				}

				p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
					return *input.Name == clusterName
				})).Return(&awseks.DescribeClusterOutput{
					Cluster: testutils.NewFakeCluster(clusterName, awseks.ClusterStatusDeleting),
				}, nil)
			})

			JustBeforeEach(func() {
				cluster, err = c.GetCluster(context.Background(), clusterName)
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have called AWS EKS service once", func() {
				p.MockEKS().AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)
			})

			It("should not call AWS CFN ListStacks", func() {
				p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacks", 0)
			})
		})

	})

	Describe("can get OIDC issuer URL and host fingerprint", func() {
		var (
			ctl *ClusterProvider
			cfg *api.ClusterConfig
			err error

			issuer = "https://exampleIssuer.eksctl.io/id/13EBFE0C5BD60778E91DFE559E02689C"
		)

		BeforeEach(func() {
			p := mockprovider.NewMockProvider()
			ctl = &ClusterProvider{
				Provider: p,
				Status:   &ProviderStatus{},
			}

			cfg = api.NewClusterConfig()

			describeClusterOutput := &awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster("testcluster", awseks.ClusterStatusActive),
			}

			describeClusterOutput.Cluster.Version = aws.String(api.Version1_21)

			describeClusterOutput.Cluster.Identity = &awseks.Identity{
				Oidc: &awseks.OIDC{
					Issuer: &issuer,
				},
			}

			p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				return true
			})).Return(describeClusterOutput, nil)
		})

		It("should get cluster, cache status and construct OIDC manager", func() {
			err = ctl.RefreshClusterStatus(cfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(cfg.Status.Endpoint).To(Equal("https://localhost/"))
			Expect(cfg.Status.CertificateAuthorityData).To(Equal([]byte("test\n")))

			Expect(ctl.ControlPlaneVersion()).To(Equal(api.Version1_21))

			_, err := ctl.NewOpenIDConnectManager(cfg)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	type platformVersionCase struct {
		platformVersion string
		expectedVersion int
		expectError     bool
	}

	DescribeTable("PlatformVersion", func(m *platformVersionCase) {
		actualVersion, err := PlatformVersion(m.platformVersion)
		if m.expectError {
			Expect(err).To(HaveOccurred())
		}
		Expect(actualVersion).To(Equal(m.expectedVersion))
	},
		Entry("eks.1", &platformVersionCase{
			platformVersion: "eks.1", expectedVersion: 1, expectError: false,
		}),
		Entry("eks.2", &platformVersionCase{
			platformVersion: "eks.2", expectedVersion: 2, expectError: false,
		}),
		Entry("eks.3", &platformVersionCase{
			platformVersion: "eks.3", expectedVersion: 3, expectError: false,
		}),
		Entry("eks.4", &platformVersionCase{
			platformVersion: "eks.4", expectedVersion: 4, expectError: false,
		}),
		Entry("eks.5", &platformVersionCase{
			platformVersion: "eks.5", expectedVersion: 5, expectError: false,
		}),
		Entry("eks.6", &platformVersionCase{
			platformVersion: "eks.6", expectedVersion: 6, expectError: false,
		}),
		Entry("eks.7", &platformVersionCase{
			platformVersion: "eks.7", expectedVersion: 7, expectError: false,
		}),
		Entry("eks. should raise an error", &platformVersionCase{
			platformVersion: "eks.", expectedVersion: -1, expectError: true,
		}),
		Entry("eks.invalid should raise an error", &platformVersionCase{
			platformVersion: "eks.invalid", expectedVersion: -1, expectError: true,
		}),
	)
})
