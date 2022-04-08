package defaultaddons_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	da "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("KubeProxy", func() {
	var (
		clientSet    kubernetes.Interface
		input        da.AddonInput
		mockProvider *mockprovider.MockProvider
	)

	Context("IsKubeProxyUpToDate", func() {
		BeforeEach(func() {
			mockProvider = mockprovider.NewMockProvider()
			input = da.AddonInput{
				Region:              "eu-west-1",
				EKSAPI:              mockProvider.EKS(),
				ControlPlaneVersion: "1.19.1",
			}

			mockProvider.MockEKS().On("DescribeAddonVersions", &awseks.DescribeAddonVersionsInput{
				AddonName:         aws.String("kube-proxy"),
				KubernetesVersion: aws.String("1.19"),
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []*awseks.AddonInfo{
					{
						AddonName: aws.String("kube-proxy"),
						AddonVersions: []*awseks.AddonVersionInfo{
							{
								AddonVersion: aws.String("v1.17.0-eksbuild.1"),
							},
							{
								//latest, unordered list to ensure we sort correctly
								AddonVersion: aws.String("v1.18.1-eksbuild.2"),
							},
							{
								AddonVersion: aws.String("v1.18.1-eksbuild.1"),
							},
						},
					},
				},
			}, nil)
		})

		When("its not up-to-date", func() {
			BeforeEach(func() {
				rawClient := testutils.NewFakeRawClientWithSamples("testdata/sample-1.16-eksbuild.1.json")
				input.RawClient = rawClient
				clientSet = rawClient.ClientSet()
			})

			It("returns false", func() {
				needsUpdating, err := da.IsKubeProxyUpToDate(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(needsUpdating).To(BeFalse())
			})
		})

		When("when its up-to-date", func() {
			BeforeEach(func() {
				rawClient := testutils.NewFakeRawClientWithSamples("testdata/sample-1.19.json")
				input.RawClient = rawClient
				clientSet = rawClient.ClientSet()
			})

			It("returns true", func() {
				needsUpdating, err := da.IsKubeProxyUpToDate(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(needsUpdating).To(BeTrue())
			})
		})

		When("it doesn't exist", func() {
			BeforeEach(func() {
				rawClient := testutils.NewFakeRawClient()
				input.RawClient = rawClient
				clientSet = rawClient.ClientSet()
			})

			// if it doesn't exist it doesn't need updating, so its up to date ¯\_(ツ)_/¯ according to #2667
			It("returns true", func() {
				needsUpdating, err := da.IsKubeProxyUpToDate(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(needsUpdating).To(BeTrue())
			})
		})

		When("it has an existing invalid image tag", func() {
			BeforeEach(func() {
				rawClient := testutils.NewFakeRawClientWithSamples("testdata/sample-1.15-invalid-image.json")
				input.RawClient = rawClient
				clientSet = rawClient.ClientSet()
			})
			It("errors", func() {
				_, err := da.IsKubeProxyUpToDate(input)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("UpdateKubeProxyImageTag", func() {
		BeforeEach(func() {
			rawClient := testutils.NewFakeRawClientWithSamples("testdata/sample-1.19.json")
			clientSet = rawClient.ClientSet()
			mockProvider = mockprovider.NewMockProvider()
			input = da.AddonInput{
				RawClient:           rawClient,
				Region:              "eu-west-1",
				EKSAPI:              mockProvider.EKS(),
				ControlPlaneVersion: "1.19.1",
			}
		})

		When("nodeaffinity is not set on kube-proxy", func() {
			BeforeEach(func() {
				input.RawClient = testutils.NewFakeRawClientWithSamples("testdata/sample-old-kube-proxy-amd.json")
				mockProvider.MockEKS().On("DescribeAddonVersions", &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: aws.String("1.19"),
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []*awseks.AddonInfo{
						{
							AddonName: aws.String("kube-proxy"),
							AddonVersions: []*awseks.AddonVersionInfo{
								{
									AddonVersion: aws.String("v1.17.0-eksbuild.1"),
								},
								{
									//latest, unordered list to ensure we sort correctly
									AddonVersion: aws.String("v1.18.1-eksbuild.2"),
								},
								{
									AddonVersion: aws.String("v1.18.1-eksbuild.1"),
								},
							},
						},
					},
				}, nil)
			})
			It("errors", func() {
				_, err := da.UpdateKubeProxy(input, false)
				Expect(err).To(MatchError(ContainSubstring("NodeAffinity not configured on kube-proxy. Either manually update the proxy deployment, or switch to Managed Addons")))
			})
		})

		When("the version reported by EKS API is more up-to-date than the default cluster version", func() {
			BeforeEach(func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: aws.String("1.19"),
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []*awseks.AddonInfo{
						{
							AddonName: aws.String("kube-proxy"),
							AddonVersions: []*awseks.AddonVersionInfo{
								{
									AddonVersion: aws.String("v1.17.0-eksbuild.1"),
								},
								{
									//latest, unordered list to ensure we sort correctly
									AddonVersion: aws.String("v1.19.1-eksbuild.2"),
								},
								{
									AddonVersion: aws.String("v1.19.1-eksbuild.1"),
								},
							},
						},
					},
				}, nil)
			})

			It("uses the image version from the EKS api", func() {
				_, err := da.UpdateKubeProxy(input, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(kubeProxyImage(clientSet)).To(Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/eks/kube-proxy:v1.19.1-eksbuild.2"))
				Expect(kubeProxyNodeSelectorValues(clientSet)).To(ConsistOf("amd64", "arm64"))
			})
		})

		When("the version reported by EKS API is behind the default cluster version", func() {
			BeforeEach(func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: aws.String("1.19"),
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []*awseks.AddonInfo{
						{
							AddonName: aws.String("kube-proxy"),
							AddonVersions: []*awseks.AddonVersionInfo{
								{
									AddonVersion: aws.String("v1.17.0-eksbuild.1"),
								},
								{
									//latest, unordered list to ensure we sort correctly
									//behind the default-cluster version 1.18.1-eksbuild.1
									AddonVersion: aws.String("v1.18.0-eksbuild.2"),
								},
								{
									AddonVersion: aws.String("v1.18.0-eksbuild.1"),
								},
							},
						},
					},
				}, nil)
			})

			It("uses the default cluster version", func() {
				_, err := da.UpdateKubeProxy(input, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(kubeProxyImage(clientSet)).To(Equal("602401143452.dkr.ecr.eu-west-1.amazonaws.com/eks/kube-proxy:v1.19.1-eksbuild.1"))
			})
		})

		When("there are no versions returned", func() {
			BeforeEach(func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: aws.String("1.19"),
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []*awseks.AddonInfo{
						{
							AddonName:     aws.String("kube-proxy"),
							AddonVersions: []*awseks.AddonVersionInfo{},
						},
					},
				}, nil)
			})

			It("returns an error", func() {
				_, err := da.UpdateKubeProxy(input, false)
				Expect(err).To(MatchError(ContainSubstring("no versions available for \"kube-proxy\"")))
			})
		})

		When("there are no valid versions returned", func() {
			BeforeEach(func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: aws.String("1.19"),
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []*awseks.AddonInfo{
						{
							AddonName: aws.String("kube-proxy"),
							AddonVersions: []*awseks.AddonVersionInfo{
								{
									AddonVersion: aws.String("not-a.1valid-version!?"),
								},
							},
						},
					},
				}, nil)
			})

			It("returns an error", func() {
				_, err := da.UpdateKubeProxy(input, false)
				Expect(err).To(MatchError(ContainSubstring("failed to parse version")))
			})
		})

		When("describing the addon errors", func() {
			BeforeEach(func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: aws.String("1.19"),
				}).Return(&awseks.DescribeAddonVersionsOutput{}, fmt.Errorf("foo"))
			})

			It("returns an error", func() {
				_, err := da.UpdateKubeProxy(input, false)
				Expect(err).To(MatchError(ContainSubstring("failed to describe addon versions: foo")))
			})
		})
	})

})

func kubeProxyImage(clientSet kubernetes.Interface) string {
	kubeProxy, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), da.KubeProxy, metav1.GetOptions{})

	Expect(err).NotTo(HaveOccurred())
	Expect(kubeProxy).NotTo(BeNil())
	Expect(kubeProxy.Spec.Template.Spec.Containers).To(HaveLen(1))

	return kubeProxy.Spec.Template.Spec.Containers[0].Image
}

func kubeProxyNodeSelectorValues(clientSet kubernetes.Interface) []string {
	kubeProxy, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), da.KubeProxy, metav1.GetOptions{})

	Expect(err).NotTo(HaveOccurred())
	Expect(kubeProxy).NotTo(BeNil())

	for _, nodeSelector := range kubeProxy.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions {
		if nodeSelector.Key == "kubernetes.io/arch" {
			return nodeSelector.Values
		}
	}
	return []string{}
}
