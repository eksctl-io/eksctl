// +build integration

package integration_test

import (
	"fmt"
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	harness "github.com/dlespiau/kube-test-harness"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/utils/names"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
)

var _ = Describe("(Integration) Create Managed Nodegroups", func() {

	const (
		initialNodeGroup = "managed-ng-0"
		newNodeGroup     = "ng-1"
	)

	defaultTimeout := 10 * time.Minute

	Describe("when creating a cluster with 1 managed nodegroup", func() {
		// always generate a unique name
		managedClusterName := "managed-" + names.ForCluster("", "")

		It("should not return an error", func() {
			fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", kubeconfigPath)

			cmd := eksctlCreateCmd.WithArgs(
				"cluster",
				"--verbose", "4",
				"--name", managedClusterName,
				"--tags", "alpha.eksctl.io/description=eksctl integration test",
				"--nodegroup-name", initialNodeGroup,
				"--node-labels", "ng-name="+initialNodeGroup,
				"--nodes", "2",
				"--version", version,
				"--kubeconfig", kubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			awsSession := NewSession(region)

			Expect(awsSession).To(HaveExistingCluster(managedClusterName, awseks.ClusterStatusActive, version))

			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", managedClusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", managedClusterName, initialNodeGroup)))
		})

		It("should have created a valid kubectl config file", func() {
			config, err := clientcmd.LoadFromFile(kubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())

			err = clientcmd.ConfirmUsable(*config, "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(config.CurrentContext).To(ContainSubstring("eksctl"))
			Expect(config.CurrentContext).To(ContainSubstring(managedClusterName))
			Expect(config.CurrentContext).To(ContainSubstring(region))
		})

		Context("and listing clusters", func() {
			It("should return the previously created cluster", func() {
				cmd := eksctlGetCmd.WithArgs("clusters", "--all-regions")
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(managedClusterName)))
			})
		})

		Context("and scale the initial nodegroup", func() {
			It("should not return an error", func() {
				cmd := eksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", managedClusterName,
					"--nodes", "3",
					"--name", initialNodeGroup,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})

		Context("and add a managed nodegroup", func() {
			It("should not return an error", func() {
				cmd := eksctlCreateCmd.WithArgs(
					"nodegroup",
					"--cluster", managedClusterName,
					"--nodes", "4",
					"--managed",
					newNodeGroup,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			Context("create test workloads", func() {
				var (
					err  error
					test *harness.Test
				)

				BeforeEach(func() {
					test, err = newKubeTest()
					Expect(err).ShouldNot(HaveOccurred())
				})

				AfterEach(func() {
					test.Close()
					Eventually(func() int {
						return len(test.ListPods(test.Namespace, metav1.ListOptions{}).Items)
					}, "3m", "1s").Should(BeZero())
				})

				It("should deploy podinfo service to the cluster and access it via proxy", func() {
					d := test.CreateDeploymentFromFile(test.Namespace, "podinfo.yaml")
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
					d := test.CreateDaemonSetFromFile(test.Namespace, "test-dns.yaml")

					test.WaitForDaemonSetReady(d, defaultTimeout)

					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

				It("should have access to HTTP(S) sites", func() {
					d := test.CreateDaemonSetFromFile(test.Namespace, "test-http.yaml")

					test.WaitForDaemonSetReady(d, defaultTimeout)

					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

			})

			Context("and delete the second nodegroup", func() {
				It("should not return an error", func() {
					cmd := eksctlDeleteCmd.WithArgs(
						"nodegroup",
						"--verbose", "4",
						"--cluster", managedClusterName,
						newNodeGroup,
					)
					Expect(cmd).To(RunSuccessfully())
				})
			})
		})

		Context("and deleting the cluster", func() {
			It("should not return an error", func() {
				cmd := eksctlDeleteClusterCmd.WithArgs(
					"--name", managedClusterName,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})
	})

})
