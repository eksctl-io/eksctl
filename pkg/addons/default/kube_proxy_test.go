package defaultaddons_test

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	da "github.com/weaveworks/eksctl/pkg/addons/default"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("KubeProxy", func() {
	var (
		clientSet           kubernetes.Interface
		input               da.AddonInput
		mockProvider        *mockprovider.MockProvider
		kubernetesVersion   = aws.String("1.23")
		controlPlaneVersion = "1.23.1"
	)

	Context("UpdateKubeProxyImageTag", func() {
		BeforeEach(func() {
			rawClient := testutils.NewFakeRawClientWithSamples("testdata/sample-1.23.json")
			clientSet = rawClient.ClientSet()
			mockProvider = mockprovider.NewMockProvider()
			input = da.AddonInput{
				RawClient:             rawClient,
				Region:                "eu-west-1",
				AddonVersionDescriber: mockProvider.EKS(),
				ControlPlaneVersion:   controlPlaneVersion,
			}
		})

		When("nodeAffinity is not set on kube-proxy", func() {
			BeforeEach(func() {
				input.RawClient = testutils.NewFakeRawClientWithSamples("testdata/sample-old-kube-proxy-amd.json")
				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: kubernetesVersion,
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []ekstypes.AddonInfo{
						{
							AddonName: aws.String("kube-proxy"),
							AddonVersions: []ekstypes.AddonVersionInfo{
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
				_, err := da.UpdateKubeProxy(context.Background(), input, false)
				Expect(err).To(MatchError(ContainSubstring("NodeAffinity not configured on kube-proxy. Either manually update the proxy deployment, or switch to Managed Addons")))
			})
		})

		type versionUpdateEntry struct {
			addonOutput ekstypes.AddonInfo

			expectedImageTag string
		}
		DescribeTable("kube-proxy version update", func(e versionUpdateEntry) {
			mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, &awseks.DescribeAddonVersionsInput{
				AddonName:         aws.String("kube-proxy"),
				KubernetesVersion: kubernetesVersion,
			}).Return(&awseks.DescribeAddonVersionsOutput{
				Addons: []ekstypes.AddonInfo{e.addonOutput},
			}, nil)

			_, err := da.UpdateKubeProxy(context.Background(), input, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(kubeProxyImage(clientSet)).To(Equal(fmt.Sprintf("602401143452.dkr.ecr.eu-west-1.amazonaws.com/eks/kube-proxy:%s", e.expectedImageTag)))
			Expect(kubeProxyNodeSelectorValues(clientSet)).To(ConsistOf("amd64", "arm64"))
		},
			Entry("a more up-to-date version should use a minimal container image", versionUpdateEntry{
				addonOutput: ekstypes.AddonInfo{
					AddonName: aws.String("kube-proxy"),
					AddonVersions: []ekstypes.AddonVersionInfo{
						{
							AddonVersion: aws.String("v1.17.0-eksbuild.1"),
						},
						{
							// Latest, unordered list to ensure we sort correctly.
							AddonVersion: aws.String("v1.23.1-eksbuild.2"),
						},
						{
							AddonVersion: aws.String("v1.23.1-eksbuild.1"),
						},
					},
				},
				expectedImageTag: "v1.23.1-minimal-eksbuild.2",
			}),

			Entry("a more up-to-date version that lacks a pre-release version should be returned unchanged", versionUpdateEntry{
				addonOutput: ekstypes.AddonInfo{
					AddonName: aws.String("kube-proxy"),
					AddonVersions: []ekstypes.AddonVersionInfo{
						{
							AddonVersion: aws.String("v1.17.0"),
						},
						{
							AddonVersion: aws.String("v1.23.2"),
						},
						{
							AddonVersion: aws.String("v1.23.1"),
						},
					},
				},

				expectedImageTag: "v1.23.2",
			}),

			Entry("a more up-to-date version that lacks a `v` prefix should not have a `v` prefix", versionUpdateEntry{
				addonOutput: ekstypes.AddonInfo{
					AddonName: aws.String("kube-proxy"),
					AddonVersions: []ekstypes.AddonVersionInfo{
						{
							AddonVersion: aws.String("1.17.0-eksbuild.1"),
						},
						{
							// Latest, unordered list to ensure we sort correctly.
							AddonVersion: aws.String("1.23.1-eksbuild.2"),
						},
						{
							AddonVersion: aws.String("1.23.1-eksbuild.1"),
						},
					},
				},

				expectedImageTag: "1.23.1-minimal-eksbuild.2",
			}),

			Entry("version that is behind the default cluster version should not be used", versionUpdateEntry{
				addonOutput: ekstypes.AddonInfo{
					AddonName: aws.String("kube-proxy"),
					AddonVersions: []ekstypes.AddonVersionInfo{
						{
							AddonVersion: aws.String("v1.17.0-eksbuild.1"),
						},
						{
							// Latest, unordered list to ensure we sort correctly
							// behind the default-cluster version 1.18.1-eksbuild.1.
							AddonVersion: aws.String("v1.18.0-eksbuild.2"),
						},
						{
							AddonVersion: aws.String("v1.18.0-eksbuild.1"),
						},
					},
				},

				expectedImageTag: "v1.23.1-eksbuild.1",
			}),
		)

		When("there are no versions returned", func() {
			BeforeEach(func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: kubernetesVersion,
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []ekstypes.AddonInfo{
						{
							AddonName:     aws.String("kube-proxy"),
							AddonVersions: []ekstypes.AddonVersionInfo{},
						},
					},
				}, nil)
			})

			It("returns an error", func() {
				_, err := da.UpdateKubeProxy(context.Background(), input, false)
				Expect(err).To(MatchError(ContainSubstring("no versions available for \"kube-proxy\"")))
			})
		})

		When("there are no valid versions returned", func() {
			BeforeEach(func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: kubernetesVersion,
				}).Return(&awseks.DescribeAddonVersionsOutput{
					Addons: []ekstypes.AddonInfo{
						{
							AddonName: aws.String("kube-proxy"),
							AddonVersions: []ekstypes.AddonVersionInfo{
								{
									AddonVersion: aws.String("not-a.1valid-version!?"),
								},
							},
						},
					},
				}, nil)
			})

			It("returns an error", func() {
				_, err := da.UpdateKubeProxy(context.Background(), input, false)
				Expect(err).To(MatchError(ContainSubstring("failed to parse version")))
			})
		})

		When("describing the addon errors", func() {
			BeforeEach(func() {
				mockProvider.MockEKS().On("DescribeAddonVersions", mock.Anything, &awseks.DescribeAddonVersionsInput{
					AddonName:         aws.String("kube-proxy"),
					KubernetesVersion: kubernetesVersion,
				}).Return(&awseks.DescribeAddonVersionsOutput{}, fmt.Errorf("foo"))
			})

			It("returns an error", func() {
				_, err := da.UpdateKubeProxy(context.Background(), input, false)
				Expect(err).To(MatchError(ContainSubstring("failed to describe addon versions: foo")))
			})
		})
	})

})

func kubeProxyImage(clientSet kubernetes.Interface) string {
	kubeProxy, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.Background(), da.KubeProxy, metav1.GetOptions{})

	Expect(err).NotTo(HaveOccurred())
	Expect(kubeProxy).NotTo(BeNil())
	Expect(kubeProxy.Spec.Template.Spec.Containers).To(HaveLen(1))

	return kubeProxy.Spec.Template.Spec.Containers[0].Image
}

func kubeProxyNodeSelectorValues(clientSet kubernetes.Interface) []string {
	kubeProxy, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.Background(), da.KubeProxy, metav1.GetOptions{})

	Expect(err).NotTo(HaveOccurred())
	Expect(kubeProxy).NotTo(BeNil())

	for _, nodeSelector := range kubeProxy.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions {
		if nodeSelector.Key == corev1.LabelArchStable {
			return nodeSelector.Values
		}
	}
	return []string{}
}
