package eks_test

import (
	"fmt"
	"io/ioutil"
	"os"

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
		c      *ClusterProvider
		p      *mockprovider.MockProvider
		output string
	)

	BeforeEach(func() {
		output = "json"
	})

	Describe("ListAll", func() {
		Context("With a cluster name", func() {
			var (
				clusterName string
				err         error
			)

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
					err = c.ListClusters(clusterName, 100, output, false)
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EKS service once", func() {
					Expect(p.MockEKS().AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)).To(BeTrue())
				})

				It("should not call AWS CFN ListStacksPages", func() {
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 0)).To(BeTrue())
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
					err = c.ListClusters(clusterName, 100, output, false)
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EKS service once", func() {
					Expect(p.MockEKS().AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)).To(BeTrue())
				})

				It("should have called AWS CFN ListStacksPages", func() {
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 1)).To(BeTrue())
				})
			})
		})

		Context("with a cluster name but cluster isn't ready", func() {
			var (
				clusterName    string
				err            error
				originalStdout *os.File
				reader         *os.File
				writer         *os.File
			)

			BeforeEach(func() {
				originalStdout = os.Stdout
				reader, writer, _ = os.Pipe()
				os.Stdout = writer

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
				err = c.ListClusters(clusterName, 100, output, false)
			})

			AfterEach(func() {
				os.Stdout = originalStdout
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have called AWS EKS service once", func() {
				Expect(p.MockEKS().AssertNumberOfCalls(GinkgoT(), "DescribeCluster", 1)).To(BeTrue())
			})

			It("should not call AWS CFN ListStacksPages", func() {
				Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "ListStacksPages", 0)).To(BeTrue())
			})

			It("the output should equal the golden file singlecluster_deleting.golden", func() {
				writer.Close()
				g, err := ioutil.ReadFile("testdata/singlecluster_deleting.golden")
				if err != nil {
					GinkgoT().Fatalf("failed reading .golden: %s", err)
				}

				actualOutput, _ := ioutil.ReadAll(reader)

				Expect(actualOutput).Should(MatchJSON(string(g)))
			})
		})

		Context("with no cluster name", func() {
			var (
				err       error
				chunkSize int
			)

			Context("and chunk-size of 1", func() {
				var (
					callNumber int
				)
				BeforeEach(func() {
					chunkSize = 1
					callNumber = 0

					p = mockprovider.NewMockProvider()

					c = &ClusterProvider{
						Provider: p,
					}

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
				})

				JustBeforeEach(func() {
					err = c.ListClusters("", chunkSize, output, false)
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EKS service twice", func() {
					Expect(p.MockEKS().AssertNumberOfCalls(GinkgoT(), "ListClusters", 2)).To(BeTrue())
				})
			})
			Context("and chunk-size of 100", func() {
				BeforeEach(func() {
					chunkSize = 100

					p = mockprovider.NewMockProvider()

					c = &ClusterProvider{
						Provider: p,
					}

					mockResultFn := func(_ *awseks.ListClustersInput) *awseks.ListClustersOutput {
						output := &awseks.ListClustersOutput{
							Clusters: []*string{aws.String("cluster-1"), aws.String("cluster-2")},
						}
						return output
					}

					p.MockEKS().On("ListClusters", mock.MatchedBy(func(input *awseks.ListClustersInput) bool {
						return *input.MaxResults == int64(chunkSize)
					})).Return(mockResultFn, nil)
				})

				JustBeforeEach(func() {
					err = c.ListClusters("", chunkSize, output, false)
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EKS service once", func() {
					Expect(p.MockEKS().AssertNumberOfCalls(GinkgoT(), "ListClusters", 1)).To(BeTrue())
				})
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
	)
})
