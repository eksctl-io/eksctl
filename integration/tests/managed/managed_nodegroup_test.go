//go:build integration
// +build integration

package managed

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	harness "github.com/dlespiau/kube-test-harness"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
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

const initialAl2Nodegroup = "al2-1"

var _ = BeforeSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

	cmd := params.EksctlCreateCmd.WithArgs(
		"cluster",
		"--verbose", "4",
		"--name", params.ClusterName,
		"--tags", "alpha.eksctl.io/description=eksctl integration test",
		"--managed",
		"--nodegroup-name", initialAl2Nodegroup,
		"--node-labels", "ng-name="+initialAl2Nodegroup,
		"--nodes", "2",
		"--instance-types", "t3a.xlarge",
		"--version", params.Version,
		"--kubeconfig", params.KubeconfigPath,
	)
	Expect(cmd).To(RunSuccessfully())
})

var _ = Describe("(Integration) Create Managed Nodegroups", func() {

	const (
		bottlerocketNodegroup = "bottlerocket-1"
		ubuntuNodegroup       = "ubuntu-1"
		newPublicNodeGroup    = "ng-public-1"
		newPrivateNodeGroup   = "ng-private-1"
	)

	makeClusterConfig := func() *api.ClusterConfig {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = params.ClusterName
		clusterConfig.Metadata.Region = params.Region
		clusterConfig.Metadata.Version = params.Version
		return clusterConfig
	}

	defaultTimeout := 20 * time.Minute

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

	Context("cluster with 1 al2 managed nodegroup", func() {
		Context("and add two managed nodegroups (one public and one private)", func() {
			It("should not return an error for public node group", func() {
				cmd := params.EksctlCreateCmd.WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--nodes", "4",
					"--managed",
					"--instance-types", "t3a.xlarge",
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
					"--instance-types", "t3a.xlarge",
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

			// clean up
			It("should not return an error when deleting nodegroups with taints applied", func() {
				cmd := params.EksctlDeleteCmd.WithArgs(
					"nodegroup",
					"--verbose", "4",
					"--cluster", params.ClusterName,
					"taints",
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})

		It("supports adding bottlerocket and ubuntu nodegroups with additional volumes", func() {
			clusterConfig := makeClusterConfig()
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         bottlerocketNodegroup,
						VolumeSize:   aws.Int(35),
						AMIFamily:    "Bottlerocket",
						InstanceType: "t3a.xlarge",
						Bottlerocket: &api.NodeGroupBottlerocket{
							EnableAdminContainer: api.Enabled(),
							Settings: &api.InlineDocument{
								"motd": "Bottlerocket is the future",
								"network": map[string]string{
									"hostname": "custom-bottlerocket-host",
								},
							},
						},
						AdditionalVolumes: []*api.VolumeMapping{
							{
								VolumeName: aws.String("/dev/sdb"),
							},
						},
					},
					Taints: []api.NodeGroupTaint{
						{
							Key:    "key2",
							Value:  "value2",
							Effect: "PreferNoSchedule",
						},
					},
				},
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         ubuntuNodegroup,
						VolumeSize:   aws.Int(25),
						AMIFamily:    "Ubuntu2004",
						InstanceType: "t3a.xlarge",
					},
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

			tests.AssertNodeVolumes(params.KubeconfigPath, params.Region, bottlerocketNodegroup, "/dev/sdb")

			By("correctly configuring the bottlerocket nodegroup")
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			nodeList := kubeTest.ListNodes(metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", "eks.amazonaws.com/nodegroup", bottlerocketNodegroup),
			})
			Expect(nodeList.Items).NotTo(BeEmpty())
			for _, node := range nodeList.Items {
				Expect(node.Status.NodeInfo.OSImage).To(ContainSubstring("Bottlerocket"))
				Expect(node.Labels["kubernetes.io/hostname"]).To(Equal("custom-bottlerocket-host"))
			}
			kubeTest.Close()
		})

		It("should have created an EKS cluster and 4 CloudFormation stacks", func() {
			awsSession := NewConfig(params.Region)

			Expect(awsSession).To(HaveExistingCluster(params.ClusterName, string(ekstypes.ClusterStatusActive), params.Version))

			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initialAl2Nodegroup)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, bottlerocketNodegroup)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, ubuntuNodegroup)))
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
				checkNg := func(ngName string) {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"nodegroup-health",
						"--cluster", params.ClusterName,
						"--name", ngName,
					)

					Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("active")))
				}

				checkNg(initialAl2Nodegroup)
				checkNg(bottlerocketNodegroup)
				checkNg(ubuntuNodegroup)
			})
		})

		Context("and scale the initial nodegroup", func() {
			It("should not return an error", func() {
				cmd := params.EksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", params.ClusterName,
					"--nodes-min", "2",
					"--nodes", "3",
					"--nodes-max", "4",
					"--name", initialAl2Nodegroup,
				)
				Expect(cmd).To(RunSuccessfully())
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

				ctx := context.Background()
				clusterProvider, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, clusterConfig)
				Expect(err).NotTo(HaveOccurred())
				ctl := clusterProvider.AWSProvider
				out, err := ctl.EKS().DescribeNodegroup(ctx, &awseks.DescribeNodegroupInput{
					ClusterName:   &params.ClusterName,
					NodegroupName: aws.String("update-config-ng"),
				})
				Expect(err).NotTo(HaveOccurred())
				Eventually(out.Nodegroup.UpdateConfig.MaxUnavailable).Should(Equal(aws.Int32(2)))

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

				out, err = ctl.EKS().DescribeNodegroup(ctx, &awseks.DescribeNodegroupInput{
					ClusterName:   &params.ClusterName,
					NodegroupName: aws.String("update-config-ng"),
				})
				Expect(err).NotTo(HaveOccurred())
				Eventually(out.Nodegroup.UpdateConfig.MaxUnavailable).Should(Equal(aws.Int32(1)))
			})

			// clean up
			It("should not return an error when deleting nodegroup with an update config", func() {
				cmd := params.EksctlDeleteCmd.WithArgs(
					"nodegroup",
					"--verbose", "4",
					"--cluster", params.ClusterName,
					"update-config-ng",
				)
				Expect(cmd).To(RunSuccessfully())
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

				upgradeNg := func(ngName string) {
					By(fmt.Sprintf("upgrading nodegroup %s to Kubernetes version %s", ngName, nextVersion))
					cmd = params.EksctlUpgradeCmd.WithArgs(
						"nodegroup",
						"--verbose", "4",
						"--cluster", params.ClusterName,
						"--name", ngName,
						"--kubernetes-version", nextVersion,
						"--timeout=60m", // wait for CF stacks to finish update
					)
					ExpectWithOffset(1, cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("nodegroup successfully upgraded")))
				}

				upgradeNg(initialAl2Nodegroup)
				upgradeNg(bottlerocketNodegroup)
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
