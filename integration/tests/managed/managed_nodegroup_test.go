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
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	"github.com/weaveworks/eksctl/integration/utilities/unowned"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	k8sUpdatePollInterval = "2s"
	k8sUpdatePollTimeout  = "3m"
)

var (
	params *tests.Params
	cfg    *api.ClusterConfig
)

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

	defaultTimeout := 20 * time.Minute
	var unownedCluster *unowned.Cluster

	BeforeSuite(func() {
		fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)
		cfg = api.NewClusterConfig()
		cfg.Metadata = &api.ClusterMeta{
			Name:    params.ClusterName,
			Region:  params.Region,
			Version: params.Version,
		}

		if params.UnownedCluster {
			unownedCluster = unowned.NewCluster(cfg)
			cfg.VPC = unownedCluster.VPC

			cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: initialNodeGroup,
						ScalingConfig: &api.ScalingConfig{
							DesiredCapacity: aws.Int(2),
						},
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
				WithStdinJSONContent(cfg)
			Expect(cmd).To(RunSuccessfully())

			cmd = params.EksctlUtilsCmd.WithArgs(
				"write-kubeconfig",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"--kubeconfig", params.KubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())
		} else {
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
		}
	})

	AfterSuite(func() {
		params.DeleteClusters()
		if params.UnownedCluster {
			unownedCluster.DeleteStack()
		}
	})

	Context("cluster with 1 managed nodegroup", func() {
		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			awsSession := NewSession(params.Region)

			Expect(awsSession).To(HaveExistingCluster(params.ClusterName, awseks.ClusterStatusActive, params.Version))

			if !params.UnownedCluster {
				Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			}

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
				cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name: newPublicNodeGroup,
							ScalingConfig: &api.ScalingConfig{
								DesiredCapacity: aws.Int(4),
							},
							PrivateNetworking: false,
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
					WithStdinJSONContent(cfg)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should not return an error for private node group", func() {
				cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name: newPrivateNodeGroup,
							ScalingConfig: &api.ScalingConfig{
								DesiredCapacity: aws.Int(2),
							},
							PrivateNetworking: true,
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
					WithStdinJSONContent(cfg)
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
				Expect(err).ToNot(HaveOccurred())

				clientset, err := kubernetes.NewForConfig(config)
				Expect(err).ToNot(HaveOccurred())
				Eventually(func() string {
					serverVersion, err := clientset.ServerVersion()
					Expect(err).ToNot(HaveOccurred())
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
