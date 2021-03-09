package eks_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
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

	Describe("ListClusters", func() {
		var (
			err            error
			chunkSize      int
			listAllRegions bool
			clusters       []*api.ClusterConfig
		)

		BeforeEach(func() {
			p = mockprovider.NewMockProvider()
			c = &ClusterProvider{
				Provider: p,
			}
			listAllRegions = false
		})

		When("and chunk-size of 1", func() {
			When("the clusters are not eksctl created", func() {

				var (
					callNumber int
				)
				BeforeEach(func() {
					chunkSize = 1
					callNumber = 0

					mockResultFn := func(_ *awseks.ListClustersInput) *awseks.ListClustersOutput {
						clusterName := fmt.Sprintf("cluster-%d", callNumber)
						output := &awseks.ListClustersOutput{
							Clusters: []*string{aws.String(clusterName)},
						}
						if callNumber == 0 {
							output.NextToken = aws.String("SOMERANDOMTOKEN")
						}

						callNumber++
						return output
					}

					p.MockEKS().On("ListClusters", mock.MatchedBy(func(input *awseks.ListClustersInput) bool {
						return *input.MaxResults == int64(chunkSize)
					})).Return(mockResultFn, nil)

					p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(nil)
				})

				JustBeforeEach(func() {
					clusters, err = c.ListClusters(chunkSize, listAllRegions)
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should return a list of clusters", func() {
					Expect(clusters).To(ConsistOf(
						&api.ClusterConfig{
							Metadata: &api.ClusterMeta{
								Name:   "cluster-0",
								Region: c.Provider.Region(),
							},
							Status: &api.ClusterStatus{
								EKSCTLCreated: "False",
							},
						},
						&api.ClusterConfig{
							Metadata: &api.ClusterMeta{
								Name:   "cluster-1",
								Region: c.Provider.Region(),
							},
							Status: &api.ClusterStatus{
								EKSCTLCreated: "False",
							},
						},
					))
				})

				It("should have called AWS EKS service twice", func() {
					p.MockEKS().AssertNumberOfCalls(GinkgoT(), "ListClusters", 2)
				})

				It("should check if the clusters are eksctl created", func() {
					p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 1)
				})
			})
		})

		Context("and chunk-size of 100", func() {
			BeforeEach(func() {
				chunkSize = 100

				mockResultFn := func(_ *awseks.ListClustersInput) *awseks.ListClustersOutput {
					output := &awseks.ListClustersOutput{
						Clusters: []*string{aws.String("cluster-1"), aws.String("cluster-2")},
					}
					return output
				}

				p.MockEKS().On("ListClusters", mock.MatchedBy(func(input *awseks.ListClustersInput) bool {
					return *input.MaxResults == int64(chunkSize)
				})).Return(mockResultFn, nil)

				p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(nil)
			})

			JustBeforeEach(func() {
				clusters, err = c.ListClusters(chunkSize, listAllRegions)
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have called AWS EKS service once", func() {
				p.MockEKS().AssertNumberOfCalls(GinkgoT(), "ListClusters", 1)
			})
		})

		Context("`listAllRegions` is true", func() {
			BeforeEach(func() {
				chunkSize = 100
				listAllRegions = true

				p.MockEC2().On("DescribeRegions", mock.Anything).Return(&ec2.DescribeRegionsOutput{}, nil)
			})

			JustBeforeEach(func() {
				clusters, err = c.ListClusters(chunkSize, listAllRegions)
			})

			It("should have called AWS EC2 service once", func() {
				p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeRegions", 1)
			})
		})
	})

	Describe("ListClusters when no clusters exist", func() {
		It("should return empty slice", func() {
			p = mockprovider.NewMockProvider()

			c = &ClusterProvider{
				Provider: p,
				Status:   &ProviderStatus{},
			}

			mockResultFn := func(_ *awseks.ListClustersInput) *awseks.ListClustersOutput {
				output := &awseks.ListClustersOutput{
					Clusters: []*string{},
				}
				return output
			}

			chunkSize := 1

			p.MockEKS().On("ListClusters", mock.MatchedBy(func(input *awseks.ListClustersInput) bool {
				return *input.MaxResults == int64(chunkSize)
			})).Return(mockResultFn, nil)

			p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(nil)

			clusters, err := c.ListClusters(chunkSize, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(clusters).ToNot(BeNil())
			Expect(len(clusters)).To(Equal(0))
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

	Describe("CanDelete", func() {
		cfg := api.NewClusterConfig()
		It("not yet created clusters are deletable", func() {
			p = mockprovider.NewMockProvider()

			c = &ClusterProvider{
				Provider: p,
				Status:   &ProviderStatus{},
			}

			p.MockEKS().On("DescribeCluster", mock.Anything).
				Return(nil, awserr.New(awseks.ErrCodeResourceNotFoundException, "", nil))

			canDelete, err := c.CanDelete(cfg)
			Expect(canDelete).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		})
		It("forwards API errors", func() {
			p = mockprovider.NewMockProvider()

			c = &ClusterProvider{
				Provider: p,
				Status:   &ProviderStatus{},
			}

			p.MockEKS().On("DescribeCluster", mock.Anything).
				Return(nil, awserr.New(awseks.ErrCodeBadRequestException, "", nil))

			_, err := c.CanDelete(cfg)
			Expect(err).To(HaveOccurred())
		})
	})

	type managedNodesSupportCase struct {
		platformVersion   string
		kubernetesVersion string

		expectError bool
		supports    bool
	}

	DescribeTable("EKS managed nodes support", func(m *managedNodesSupportCase) {
		var platform *string
		if m.platformVersion != "" {
			platform = aws.String(m.platformVersion)
		}
		cluster := &awseks.Cluster{
			Version:         aws.String(m.kubernetesVersion),
			PlatformVersion: platform,
		}

		supportsManagedNodes, err := ClusterSupportsManagedNodes(cluster)
		if m.expectError {
			Expect(err).To(HaveOccurred())
		}
		Expect(supportsManagedNodes).To(Equal(m.supports))
	},
		Entry("with minimum required versions", &managedNodesSupportCase{
			platformVersion:   "eks.3",
			kubernetesVersion: "1.14",
			expectError:       false,
			supports:          true,
		}),
		Entry("with newer version", &managedNodesSupportCase{
			platformVersion:   "eks.6",
			kubernetesVersion: "1.14",
			expectError:       false,
			supports:          true,
		}),
		Entry("with unsupported platform version", &managedNodesSupportCase{
			platformVersion:   "eks.2",
			kubernetesVersion: "1.14",
			expectError:       false,
			supports:          false,
		}),
		Entry("with unsupported Kubernetes version", &managedNodesSupportCase{
			platformVersion:   "eks.3",
			kubernetesVersion: "1.13",
			expectError:       false,
			supports:          false,
		}),
		Entry("with unsupported Kubernetes version", &managedNodesSupportCase{
			platformVersion:   "eks.6",
			kubernetesVersion: "1.13",
			expectError:       false,
			supports:          false,
		}),
		Entry("with invalid platform version", &managedNodesSupportCase{
			platformVersion:   "eks.invalid",
			kubernetesVersion: "1.14",
			expectError:       true,
			supports:          false,
		}),
		Entry("with invalid platform version", &managedNodesSupportCase{
			platformVersion:   " eks.3",
			kubernetesVersion: "1.14",
			expectError:       true,
			supports:          false,
		}),
		Entry("with invalid Kubernetes version", &managedNodesSupportCase{
			platformVersion:   "eks.5",
			kubernetesVersion: "1.",
			expectError:       true,
			supports:          false,
		}),
		Entry("with non-existent platform version", &managedNodesSupportCase{
			kubernetesVersion: "1.",
			expectError:       false,
			supports:          false,
		}),
		Entry("with 1.15", &managedNodesSupportCase{
			kubernetesVersion: "1.15",
			expectError:       false,
			supports:          true,
		}),
		Entry("with 1.16", &managedNodesSupportCase{
			kubernetesVersion: "1.16",
			expectError:       false,
			supports:          true,
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

	type kmsSupportCase struct {
		key               string
		errSubstr         string
		kubernetesVersion string
	}

	DescribeTable("KMS validation", func(k kmsSupportCase) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Version = k.kubernetesVersion

		clusterConfig.SecretsEncryption = &api.SecretsEncryption{
			KeyARN: k.key,
		}
		err := ValidateFeatureCompatibility(clusterConfig, nil)
		if k.errSubstr != "" {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(k.errSubstr))
		} else {
			Expect(err).ToNot(HaveOccurred())
		}
	},
		Entry("Invalid ARN", kmsSupportCase{
			key:               "invalid:arn",
			errSubstr:         "invalid ARN",
			kubernetesVersion: "1.14",
		}),
		Entry("Valid ARN", kmsSupportCase{
			key:               "arn:aws:kms:us-west-2:000000000000:key/12345-12345",
			errSubstr:         "",
			kubernetesVersion: "1.14",
		}),
		Entry("Supports K8s 1.13", kmsSupportCase{
			key:               "arn:aws:kms:us-west-2:000000000000:key/12345-12345",
			errSubstr:         "",
			kubernetesVersion: "1.13",
		}),
		Entry("Unsupported K8s version", kmsSupportCase{
			key:               "arn:aws:kms:us-west-2:000000000000:key/12345-12345",
			errSubstr:         "KMS is only supported for EKS version 1.13 and above",
			kubernetesVersion: "1.12",
		}),
	)
})
