//go:build integration
// +build integration

package managed

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	harness "github.com/dlespiau/kube-test-harness"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/weaveworks/eksctl/pkg/eks"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/names"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	k8sUpdatePollInterval = "2s"
	k8sUpdatePollTimeout  = "3m"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("managed")
	supportedVersions := api.SupportedVersions()
	if len(supportedVersions) < 2 {
		panic("managed nodegroup tests require at least two supported Kubernetes versions to run")
	}
	params.Version = supportedVersions[len(supportedVersions)-2]
}

func TestManaged(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Create Managed Nodegroups", func() {

	const (
		initialNodeGroup    = "managed-ng-0"
		newPublicNodeGroup  = "ng-public-1"
		newPrivateNodeGroup = "ng-private-1"
	)

	makeClusterConfig := func() *api.ClusterConfig {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = params.ClusterName
		clusterConfig.Metadata.Region = params.Region
		clusterConfig.Metadata.Version = params.Version
		return clusterConfig
	}

	defaultTimeout := 20 * time.Minute

	BeforeSuite(func() {
		fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", params.ClusterName,
			"--tags", "alpha.eksctl.io/description=eksctl integration test",
			"--managed",
			"--nodegroup-name", initialNodeGroup,
			"--node-labels", "ng-name="+initialNodeGroup,
			"--nodes", "2",
			"--version", params.Version,
			"--kubeconfig", params.KubeconfigPath,
		)
		Expect(cmd).To(RunSuccessfully())
	})

	DescribeTable("Bottlerocket and Ubuntu support", func(ng *api.ManagedNodeGroup) {
		clusterConfig := makeClusterConfig()
		clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{ng}
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"nodegroup",
				"--config-file", "-",
				"--verbose", "4",
			).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.Reader(clusterConfig))

		Expect(cmd).To(RunSuccessfully())
	},
		Entry("Bottlerocket", &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:       "bottlerocket",
				VolumeSize: aws.Int(35),
				AMIFamily:  "Bottlerocket",
			},
			Taints: []api.NodeGroupTaint{
				{
					Key:    "key2",
					Value:  "value2",
					Effect: "PreferNoSchedule",
				},
			},
		}),

		Entry("Ubuntu", &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:       "ubuntu",
				VolumeSize: aws.Int(25),
				AMIFamily:  "Ubuntu2004",
			},
		}),
	)

	type managedCLIEntry struct {
		createArgs []string

		expectedErr string
	}

	DescribeTable("Managed CLI features", func(m managedCLIEntry) {
		runAssertions := func(createArgs ...string) {
			cmd := params.EksctlCreateCmd.
				WithArgs(createArgs...).
				WithArgs(m.createArgs...)

			if m.expectedErr != "" {
				session := cmd.Run()
				Expect(session.ExitCode()).NotTo(Equal(0))
				output := session.Err.Contents()
				Expect(string(output)).To(ContainSubstring(m.expectedErr))
				return
			}

			Expect(cmd).To(RunSuccessfully())
		}

		// Run the same assertions for both `create cluster` and `create nodegroup`
		runAssertions("cluster")
		runAssertions(
			"nodegroup",
			"--cluster", params.ClusterName,
		)
	},
		Entry("Windows AMI", managedCLIEntry{
			createArgs: []string{
				"--node-ami-family=WindowsServer2019FullContainer",
			},
			expectedErr: "Windows is not supported for managed nodegroups; eksctl now creates " +
				"managed nodegroups by default, to use a self-managed nodegroup, pass --managed=false",
		}),

		Entry("Windows AMI with dry-run", managedCLIEntry{
			createArgs: []string{
				"--node-ami-family=WindowsServer2019FullContainer",
				"--dry-run",
			},
			expectedErr: "Windows is not supported for managed nodegroups; eksctl now creates " +
				"managed nodegroups by default, to use a self-managed nodegroup, pass --managed=false",
		}),

		Entry("Bottlerocket with dry-run", managedCLIEntry{
			createArgs: []string{
				"--node-ami-family=Bottlerocket",
				"--instance-prefix=bottle",
				"--instance-name=rocket",
				"--dry-run",
			},
		}),

		Entry("Ubuntu with dry-run", managedCLIEntry{
			createArgs: []string{
				"--node-ami-family=Ubuntu2004",
				"--dry-run",
			},
		}),
	)

	Describe("Bottlerocket nodegroup", func() {
		var kubeTest *harness.Test

		BeforeEach(func() {
			var err error
			kubeTest, err = kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			kubeTest.Close()
		})

		assertCreateBottlerocket := func(ng *api.ManagedNodeGroup) *corev1.NodeList {
			clusterConfig := makeClusterConfig()
			ng.Name = names.ForNodeGroup("", "")

			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{ng}
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))

			Expect(cmd).To(RunSuccessfully())

			nodeList := kubeTest.ListNodes(metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", "eks.amazonaws.com/nodegroup", ng.Name),
			})
			Expect(nodeList.Items).NotTo(BeEmpty())
			for _, node := range nodeList.Items {
				Expect(node.Status.NodeInfo.OSImage).To(ContainSubstring("Bottlerocket"))
			}
			return nodeList
		}

		It("should create a standard nodegroup", func() {
			ng := &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					VolumeSize: aws.Int(35),
					AMIFamily:  "Bottlerocket",
					Labels: map[string]string{
						"ami-family": "bottlerocket",
					},
					Bottlerocket: &api.NodeGroupBottlerocket{
						EnableAdminContainer: api.Enabled(),
					},
					ScalingConfig: &api.ScalingConfig{
						DesiredCapacity: aws.Int(1),
					},
				},
				Taints: []api.NodeGroupTaint{
					{
						Key:    "key1",
						Value:  "value1",
						Effect: "PreferNoSchedule",
					},
				},
			}

			assertCreateBottlerocket(ng)
		})

		It("should create a nodegroup with custom Bottlerocket settings", func() {
			ng := &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					VolumeSize: aws.Int(20),
					AMIFamily:  "Bottlerocket",
					Labels: map[string]string{
						"ami-family": "bottlerocket",
					},
					Bottlerocket: &api.NodeGroupBottlerocket{
						EnableAdminContainer: api.Enabled(),
						Settings: &api.InlineDocument{
							"motd": "Bottlerocket is the future",
							"network": map[string]string{
								"hostname": "custom-bottlerocket-host",
							},
						},
					},
				},
				Taints: []api.NodeGroupTaint{
					{
						Key:    "key1",
						Value:  "value1",
						Effect: "PreferNoSchedule",
					},
				},
			}

			nodeList := assertCreateBottlerocket(ng)
			for _, node := range nodeList.Items {
				Expect(node.Labels["kubernetes.io/hostname"]).To(Equal("custom-bottlerocket-host"))
			}
		})
	})

	Context("cluster with 1 managed nodegroup", func() {
		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			awsSession := NewSession(params.Region)

			Expect(awsSession).To(HaveExistingCluster(params.ClusterName, awseks.ClusterStatusActive, params.Version))

			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initialNodeGroup)))
		})

		It("should have created a valid kubectl config file", func() {
			config, err := clientcmd.LoadFromFile(params.KubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())

			err = clientcmd.ConfirmUsable(*config, "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(config.CurrentContext).To(ContainSubstring("eksctl"))
			Expect(config.CurrentContext).To(ContainSubstring(params.ClusterName))
			Expect(config.CurrentContext).To(ContainSubstring(params.Region))
		})

		Context("and listing clusters", func() {
			It("should return the previously created cluster", func() {
				cmd := params.EksctlGetCmd.WithArgs("clusters", "--all-regions")
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(params.ClusterName)))
			})
		})

		Context("and checking the nodegroup health", func() {
			It("should return healthy", func() {
				cmd := params.EksctlUtilsCmd.WithArgs(
					"nodegroup-health",
					"--cluster", params.ClusterName,
					"--name", initialNodeGroup,
				)

				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("active")))
			})
		})

		Context("and scale the initial nodegroup", func() {
			It("should not return an error", func() {
				cmd := params.EksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", params.ClusterName,
					"--nodes-min", "2",
					"--nodes", "3",
					"--nodes-max", "4",
					"--name", initialNodeGroup,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})

		Context("and add two managed nodegroups (one public and one private)", func() {
			It("should not return an error for public node group", func() {
				cmd := params.EksctlCreateCmd.WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--nodes", "4",
					"--managed",
					newPublicNodeGroup,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should not return an error for private node group", func() {
				cmd := params.EksctlCreateCmd.WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--nodes", "2",
					"--managed",
					"--node-private-networking",
					newPrivateNodeGroup,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			Context("create test workloads", func() {
				var (
					err  error
					test *harness.Test
				)

				BeforeEach(func() {
					test, err = kube.NewTest(params.KubeconfigPath)
					Expect(err).ShouldNot(HaveOccurred())
				})

				AfterEach(func() {
					test.Close()
					Eventually(func() int {
						return len(test.ListPods(test.Namespace, metav1.ListOptions{}).Items)
					}, "3m", "1s").Should(BeZero())
				})

				It("should deploy podinfo service to the cluster and access it via proxy", func() {
					d := test.CreateDeploymentFromFile(test.Namespace, "../../data/podinfo.yaml")
					test.WaitForDeploymentReady(d, defaultTimeout)

					pods := test.ListPodsFromDeployment(d)
					Expect(len(pods.Items)).To(Equal(2))

					// For each pod of the Deployment, check we receive a sensible response to a
					// GET request on /version.
					for _, pod := range pods.Items {
						Expect(pod.Namespace).To(Equal(test.Namespace))

						req := test.PodProxyGet(&pod, "", "/version")
						fmt.Fprintf(GinkgoWriter, "url = %#v", req.URL())

						var js map[string]interface{}
						test.PodProxyGetJSON(&pod, "", "/version", &js)

						Expect(js).To(HaveKeyWithValue("version", "1.5.1"))
					}
				})

				It("should have functional DNS", func() {
					d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-dns.yaml")
					test.WaitForDaemonSetReady(d, defaultTimeout)
					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

				It("should have access to HTTP(S) sites", func() {
					d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-http.yaml")
					test.WaitForDaemonSetReady(d, defaultTimeout)
					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

			})

			Context("and delete the managed public nodegroup", func() {
				It("should not return an error", func() {
					cmd := params.EksctlDeleteCmd.WithArgs(
						"nodegroup",
						"--verbose", "4",
						"--cluster", params.ClusterName,
						newPublicNodeGroup,
					)
					Expect(cmd).To(RunSuccessfully())
				})
			})

			Context("and delete the managed private nodegroup", func() {
				It("should not return an error", func() {
					cmd := params.EksctlDeleteCmd.WithArgs(
						"nodegroup",
						"--verbose", "4",
						"--cluster", params.ClusterName,
						newPrivateNodeGroup,
					)
					Expect(cmd).To(RunSuccessfully())
				})
			})
		})

		Context("and upgrading a nodegroup", func() {
			It("should upgrade to the next Kubernetes version", func() {
				By("updating the control plane version")
				cmd := params.EksctlUpgradeCmd.
					WithArgs(
						"cluster",
						"--verbose", "4",
						"--name", params.ClusterName,
						"--approve",
					)
				Expect(cmd).To(RunSuccessfully())

				var nextVersion string
				{
					supportedVersions := api.SupportedVersions()
					nextVersion = supportedVersions[len(supportedVersions)-1]
				}
				By(fmt.Sprintf("checking that control plane is updated to %v", nextVersion))
				config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
				Expect(err).NotTo(HaveOccurred())

				clientset, err := kubernetes.NewForConfig(config)
				Expect(err).NotTo(HaveOccurred())
				Eventually(func() string {
					serverVersion, err := clientset.ServerVersion()
					Expect(err).NotTo(HaveOccurred())
					return fmt.Sprintf("%s.%s", serverVersion.Major, strings.TrimSuffix(serverVersion.Minor, "+"))
				}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(Equal(nextVersion))

				By(fmt.Sprintf("upgrading nodegroup %s to Kubernetes version %s", initialNodeGroup, nextVersion))
				cmd = params.EksctlUpgradeCmd.WithArgs(
					"nodegroup",
					"--verbose", "4",
					"--cluster", params.ClusterName,
					"--name", initialNodeGroup,
					"--kubernetes-version", nextVersion,
				)
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("nodegroup successfully upgraded")))
			})
		})

		Context("and creating a nodegroup with taints", func() {
			It("should create nodegroups with taints applied", func() {
				taints := []api.NodeGroupTaint{
					{
						Key:    "key1",
						Value:  "value1",
						Effect: "NoSchedule",
					},
					{
						Key:    "key2",
						Effect: "NoSchedule",
					},
					{
						Key:    "key3",
						Value:  "value2",
						Effect: "NoExecute",
					},
				}
				clusterConfig := makeClusterConfig()
				clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name: "taints",
						},
						Taints: taints,
					},
				}

				cmd := params.EksctlCreateCmd.
					WithArgs(
						"nodegroup",
						"--config-file", "-",
						"--verbose", "4",
					).
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.Reader(clusterConfig))
				Expect(cmd).To(RunSuccessfully())

				config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
				Expect(err).NotTo(HaveOccurred())
				clientset, err := kubernetes.NewForConfig(config)
				Expect(err).NotTo(HaveOccurred())

				mapTaints := func(taints []api.NodeGroupTaint) []corev1.Taint {
					var ret []corev1.Taint
					for _, t := range taints {
						ret = append(ret, corev1.Taint{
							Key:    t.Key,
							Value:  t.Value,
							Effect: t.Effect,
						})
					}
					return ret
				}

				tests.AssertNodeTaints(tests.ListNodes(clientset, "taints"), mapTaints(taints))
			})
		})

		Context("and creating a nodegroup with an update config", func() {
			It("defining the UpdateConfig field in the cluster config", func() {
				By("creating it")
				updateConfig := &api.NodeGroupUpdateConfig{
					MaxUnavailable: aws.Int(2),
				}
				clusterConfig := makeClusterConfig()
				clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name: "update-config-ng",
						},
						UpdateConfig: updateConfig,
					},
				}
				cmd := params.EksctlCreateCmd.
					WithArgs(
						"nodegroup",
						"--config-file", "-",
						"--verbose", "4",
					).
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.Reader(clusterConfig))
				Expect(cmd).To(RunSuccessfully())

				clusterProvider, err := eks.New(&api.ProviderConfig{Region: params.Region}, clusterConfig)
				Expect(err).NotTo(HaveOccurred())
				ctl := clusterProvider.Provider
				out, err := ctl.EKS().DescribeNodegroup(&awseks.DescribeNodegroupInput{
					ClusterName:   &params.ClusterName,
					NodegroupName: aws.String("update-config-ng"),
				})
				Expect(err).NotTo(HaveOccurred())
				Eventually(out.Nodegroup.UpdateConfig.MaxUnavailable).Should(Equal(aws.Int64(2)))

				By("and updating the nodegroup's UpdateConfig")
				clusterConfig.ManagedNodeGroups[0].Spot = true
				clusterConfig.ManagedNodeGroups[0].UpdateConfig = &api.NodeGroupUpdateConfig{
					MaxUnavailable: aws.Int(1),
				}

				cmd = params.EksctlUpdateCmd.
					WithArgs(
						"nodegroup",
						"--config-file", "-",
						"--verbose", "4",
					).
					WithoutArg("--region", params.Region).
					WithStdin(clusterutils.Reader(clusterConfig))

				Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
					ContainElement(ContainSubstring("unchanged fields for nodegroup update-config-ng: the following fields remain unchanged; they are not supported by `eksctl update nodegroup`: Spot")),
				))

				out, err = ctl.EKS().DescribeNodegroup(&awseks.DescribeNodegroupInput{
					ClusterName:   &params.ClusterName,
					NodegroupName: aws.String("update-config-ng"),
				})
				Expect(err).NotTo(HaveOccurred())
				Eventually(out.Nodegroup.UpdateConfig.MaxUnavailable).Should(Equal(aws.Int64(1)))
			})
		})

		Context("and deleting the cluster", func() {
			It("should not return an error", func() {
				cmd := params.EksctlDeleteClusterCmd.WithArgs(
					"--name", params.ClusterName,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
