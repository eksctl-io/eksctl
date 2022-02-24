package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
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
					cluster, err = c.GetCluster(clusterName)
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

				It("should not call AWS CFN ListStacksPages", func() {
					p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 0)
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

					p.MockCloudFormation().On("ListStacksPages", mock.MatchedBy(func(input *cfn.ListStacksInput) bool {
						matches := 0
						for i := range input.StackStatusFilter {
							if *input.StackStatusFilter[i] == expectedStatusFilter[i] {
								matches++
							}
						}
						return matches == len(expectedStatusFilter)
					}), mock.Anything).Return(nil)
				})

				JustBeforeEach(func() {
					cluster, err = c.GetCluster(clusterName)
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

				It("should have called AWS CFN ListStacksPages", func() {
					p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 1)
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
				cluster, err = c.GetCluster(clusterName)
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have called AWS EKS service once", func() {
				p.MockEKS().AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)
			})

			It("should not call AWS CFN ListStacksPages", func() {
				p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 0)
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

			describeClusterOutput.Cluster.Version = aws.String(api.Version1_13)

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

			Expect(ctl.ControlPlaneVersion()).To(Equal(api.Version1_13))

			_, err := ctl.NewOpenIDConnectManager(cfg)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	type managedNodesSupportCase struct {
		platformVersion string

		expectError bool
		supports    bool
	}

	DescribeTable("EKS managed nodes support", func(m *managedNodesSupportCase) {
		var platform *string
		if m.platformVersion != "" {
			platform = aws.String(m.platformVersion)
		}
		cluster := &awseks.Cluster{
			PlatformVersion: platform,
		}

		supportsManagedNodes, err := ClusterSupportsManagedNodes(cluster)
		if m.expectError {
			Expect(err).To(HaveOccurred())
		}
		Expect(supportsManagedNodes).To(Equal(m.supports))
	},
		Entry("with unsupported platform version", &managedNodesSupportCase{
			platformVersion: "eks.2",
			expectError:     false,
			supports:        false,
		}),
		Entry("with invalid platform version", &managedNodesSupportCase{
			platformVersion: " eks.3",
			expectError:     true,
			supports:        false,
		}),
	)

	type fargateSupportCase struct {
		platformVersion   string
		kubernetesVersion string
		supportsFargate   bool
		expectError       bool
	}

	DescribeTable("ClusterSupportsFargate", func(m *fargateSupportCase) {
		cluster := &awseks.Cluster{
			Version:         aws.String(m.kubernetesVersion),
			PlatformVersion: aws.String(m.platformVersion),
		}
		supportsFargate, err := ClusterSupportsFargate(cluster)
		if m.expectError {
			Expect(err).To(HaveOccurred())
		}
		Expect(supportsFargate).To(Equal(m.supportsFargate))
	},
		Entry("eks.1 does NOT support Fargate", &fargateSupportCase{
			platformVersion: "eks.1", kubernetesVersion: "1.14", supportsFargate: false, expectError: false,
		}),
		Entry("eks.2 does NOT support Fargate", &fargateSupportCase{
			platformVersion: "eks.2", kubernetesVersion: "1.14", supportsFargate: false, expectError: false,
		}),
		Entry("eks.3 does NOT support Fargate", &fargateSupportCase{
			platformVersion: "eks.3", kubernetesVersion: "1.14", supportsFargate: false, expectError: false,
		}),
		Entry("eks.4 does NOT support Fargate", &fargateSupportCase{
			platformVersion: "eks.4", kubernetesVersion: "1.14", supportsFargate: false, expectError: false,
		}),
		Entry("eks.5 is the minimum version which supports Fargate", &fargateSupportCase{
			platformVersion: "eks.5", kubernetesVersion: "1.14", supportsFargate: true, expectError: false,
		}),
		Entry("1.14 is the minimum version which supports Fargate", &fargateSupportCase{
			platformVersion: "eks.5", kubernetesVersion: "1.13", supportsFargate: false, expectError: false,
		}),
		Entry("eks.6 supports Fargate", &fargateSupportCase{
			platformVersion: "eks.6", kubernetesVersion: "1.14", supportsFargate: true, expectError: false,
		}),
		Entry("eks.7 supports Fargate", &fargateSupportCase{
			platformVersion: "eks.7", kubernetesVersion: "1.14", supportsFargate: true, expectError: false,
		}),
		Entry("eks. should raise an error", &fargateSupportCase{
			platformVersion: "eks.", kubernetesVersion: "1.14", supportsFargate: false, expectError: true,
		}),
		Entry("eks.invalid should raise an error", &fargateSupportCase{
			platformVersion: "eks.invalid", kubernetesVersion: "1.14", supportsFargate: false, expectError: true,
		}),
		Entry("invalid Kubernetes version should raise an error", &fargateSupportCase{
			platformVersion: "eks.5", kubernetesVersion: "1.", supportsFargate: false, expectError: true,
		}),
		Entry("should support 1.15 for all platform versions", &fargateSupportCase{
			platformVersion: "eks.1", kubernetesVersion: "1.15", supportsFargate: true, expectError: false,
		}),
		Entry("should support 1.16 for all platform versions", &fargateSupportCase{
			platformVersion: "eks.1", kubernetesVersion: "1.16", supportsFargate: true, expectError: false,
		}),
	)

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
