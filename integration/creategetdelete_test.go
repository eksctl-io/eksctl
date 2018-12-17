// +build integration

package integration_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/weaveworks/eksctl/pkg/testutils/aws"
	. "github.com/weaveworks/eksctl/pkg/testutils/matchers"
	"github.com/weaveworks/eksctl/pkg/utils"
)

type tInterface interface {
	GinkgoTInterface
	Helper()
}

type tHelper struct{ GinkgoTInterface }

func (t *tHelper) Helper()      { return }
func (t *tHelper) Name() string { return "eksctl-test" }

func newKubeTest() (*harness.Test, error) {
	t := &tHelper{GinkgoT()}
	l := &logger.TestLogger{}
	h := harness.New(harness.Options{Logger: l.ForTest(t)})
	if err := h.Setup(); err != nil {
		return nil, err
	}
	if err := h.SetKubeconfig(kubeconfigPath); err != nil {
		return nil, err
	}
	return h.NewTest(t), nil
}

var _ = Describe("(Integration) Create, Get, Scale & Delete", func() {

	BeforeSuite(func() {
		kubeconfigTemp = false
		if kubeconfigPath == "" {
			wd, _ := os.Getwd()
			f, _ := ioutil.TempFile(wd, "kubeconfig-")
			kubeconfigPath = f.Name()
			kubeconfigTemp = true
		}
	})

	AfterSuite(func() {
		gexec.KillAndWait()
		if kubeconfigTemp {
			os.Remove(kubeconfigPath)
		}
	})

	Describe("when creating a cluster with 1 node", func() {
		It("should not return an error", func() {
			if !doCreate {
				fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", clusterName)
				return
			}

			fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", kubeconfigPath)

			if clusterName == "" {
				clusterName = utils.ClusterName("", "")
			}

			args := []string{"create", "cluster",
				"--name", clusterName,
				"--tags", "eksctl.cluster.k8s.io/v1alpha1/description=eksctl integration test",
				"--node-type", "t2.medium",
				"--nodes", "1",
				"--region", region,
				"--kubeconfig", kubeconfigPath,
			}

			command := exec.Command(eksctlPath, args...)
			cmdSession, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

			if err != nil {
				Fail(fmt.Sprintf("error starting process: %v", err), 1)
			}

			cmdSession.Wait(createTimeout)
			Expect(cmdSession.ExitCode()).Should(Equal(0))
		})

		awsSession := aws.NewSession(region)

		It("should have created an EKS cluster", func() {
			Expect(awsSession).To(HaveExistingCluster(clusterName, awseks.ClusterStatusActive, "1.11"))
		})

		It("should have the required cloudformation stacks", func() {
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", clusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%d", clusterName, 0)))
		})

		It("should have created a valid kubectl config file", func() {
			config, err := clientcmd.LoadFromFile(kubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())

			err = clientcmd.ConfirmUsable(*config, "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(config.CurrentContext).To(ContainSubstring("eksctl"))
			Expect(config.CurrentContext).To(ContainSubstring(clusterName))
			Expect(config.CurrentContext).To(ContainSubstring(region))
		})

		Context("and we create a deployment using kubectl", func() {
			var (
				err  error
				test *harness.Test
			)

			BeforeEach(func() {
				test, err = newKubeTest()
				Expect(err).ShouldNot(HaveOccurred())
				test.CreateNamespace(test.Namespace)
			})

			AfterEach(func() {
				test.Close()
			})

			It("should deploy the service to the cluster", func() {
				d := test.CreateDeploymentFromFile(test.Namespace, "podinfo.yaml")
				test.WaitForDeploymentReady(d, 1*time.Minute)

				pods := test.ListPodsFromDeployment(d)
				Expect(len(pods.Items)).To(Equal(2))

				// For each pod of the Deployment, check we receive a sensible response to a
				// GET request on /version.
				for _, pod := range pods.Items {
					Expect(pod.Namespace).To(Equal(test.Namespace))

					req := test.PodProxyGet(&pod, "", "/version")
					fmt.Fprintf(GinkgoWriter, "url = %#v", req.URL())

					var js interface{}
					test.PodProxyGetJSON(&pod, "", "/version", &js)

					Expect(js.(map[string]interface{})).To(HaveKeyWithValue("version", "1.0.1"))
				}
			})
		})

		Context("and listing clusters", func() {
			var cmdSession *gexec.Session
			It("should not return an error", func() {
				var err error
				args := []string{"get", "clusters", "--region", region}

				command := exec.Command(eksctlPath, args...)
				cmdSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)

				if err != nil {
					Fail(fmt.Sprintf("error starting process: %v", err), 1)
				}

				cmdSession.Wait(getTimeout)
				Expect(cmdSession.ExitCode()).Should(Equal(0))
			})

			It("should return the previously created cluster", func() {
				Expect(string(cmdSession.Buffer().Contents())).To(ContainSubstring(clusterName))
			})
		})

		Context("and scale the cluster", func() {
			It("should not return an error", func() {
				args := []string{"scale", "nodegroup",
					"--name", clusterName,
					"--region", region,
					"--nodes", "2",
				}

				command := exec.Command(eksctlPath, args...)
				cmdSession, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				if err != nil {
					Fail(fmt.Sprintf("error starting process: %v", err), 1)
				}

				cmdSession.Wait(scaleTimeout)
				Expect(cmdSession.ExitCode()).Should(Equal(0))
			})

			It("should make it 2 nodes total", func() {
				test, err := newKubeTest()
				Expect(err).ShouldNot(HaveOccurred())
				defer test.Close()

				test.WaitForNodesReady(2, scaleTimeout)

				nodes := test.ListNodes(metav1.ListOptions{})

				Expect(len(nodes.Items)).To(Equal(2))
			})
		})

		Context("and deleting the cluster", func() {
			It("should not return an error", func() {
				if !doDelete {
					Skip("will not delete cluster " + clusterName)
				}

				args := []string{"delete", "cluster",
					"--name", clusterName,
					"--region", region,
					"--wait",
				}

				command := exec.Command(eksctlPath, args...)
				cmdSession, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)

				if err != nil {
					Fail(fmt.Sprintf("error starting process: %v", err), 1)
				}

				cmdSession.Wait(deleteTimeout)
				Expect(cmdSession.ExitCode()).Should(Equal(0))
			})

			awsSession := aws.NewSession(region)

			It("should have deleted the EKS cluster", func() {
				if !doDelete {
					Skip("will not delete cluster " + clusterName)
				}

				Expect(awsSession).ToNot(HaveExistingCluster(clusterName, awseks.ClusterStatusActive, "1.11"))
			})

			It("should have deleted the required cloudformation stacks", func() {
				if !doDelete {
					Skip("will not delete cluster " + clusterName)
				}

				Expect(awsSession).ToNot(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", clusterName)))
				Expect(awsSession).ToNot(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%d", clusterName, 0)))
			})
		})
	})
})
