// +build integration

package inferentia

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"k8s.io/client-go/kubernetes"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/file"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var defaultCluster string
var noInstallCluster string
var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("inf1")
	defaultCluster = params.ClusterName
	noInstallCluster = params.NewClusterName("inf1-no-plugin")
}

func TestInferentia(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Inferentia nodes", func() {
	const (
		initNG = "inf1-ng-0"
		newNG  = "inf1-ng-1"
	)
	BeforeSuite(func() {
		params.KubeconfigTemp = false
		if params.KubeconfigPath == "" {
			wd, _ := os.Getwd()
			f, _ := ioutil.TempFile(wd, "kubeconfig-")
			params.KubeconfigPath = f.Name()
			params.KubeconfigTemp = true
		}
	})

	AfterSuite(func() {
		params.DeleteClusters()
		gexec.KillAndWait()
		if params.KubeconfigTemp {
			os.Remove(params.KubeconfigPath)
		}
		os.RemoveAll(params.TestDirectory)
	})

	When("creating a cluster with inf1 nodes", func() {
		Context("by default", func() {
			It("should not return an error", func() {
				if params.SkipCreate {
					fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", defaultCluster)
					if !file.Exists(params.KubeconfigPath) {
						// Generate the Kubernetes configuration that eksctl create
						// would have generated otherwise:
						cmd := params.EksctlUtilsCmd.WithArgs(
							"write-kubeconfig",
							"--verbose", "4",
							"--cluster", defaultCluster,
							"--kubeconfig", params.KubeconfigPath,
						)
						Expect(cmd).To(RunSuccessfully())
					}
					return
				}

				fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

				cmd := params.EksctlCreateCmd.WithArgs(
					"cluster",
					"--verbose", "4",
					"--name", defaultCluster,
					"--tags", "alpha.eksctl.io/description=eksctl integration test",
					"--nodegroup-name", initNG,
					"--node-labels", "ng-name="+initNG,
					"--nodes", "1",
					"--node-type", "inf1.xlarge",
					"--version", params.Version,
					"--kubeconfig", params.KubeconfigPath,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should have installed the neuron device plugin", func() {
				clientSet := newClientSet()
				_, err := clientSet.AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should not have installed the nvidia device plugin", func() {
				_, err := newClientSet().AppsV1().DaemonSets("kube-system").Get(context.TODO(), "nvidia-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).Should(BeNotFoundError())
			})
		})
		Context("with --install-neuron-plugin=false", func() {
			It("should not return an error", func() {
				if params.SkipCreate {
					fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", noInstallCluster)
					if !file.Exists(params.KubeconfigPath) {
						// Generate the Kubernetes configuration that eksctl create
						// would have generated otherwise:
						cmd := params.EksctlUtilsCmd.WithArgs(
							"write-kubeconfig",
							"--verbose", "4",
							"--cluster", noInstallCluster,
							"--kubeconfig", params.KubeconfigPath,
						)
						Expect(cmd).To(RunSuccessfully())
					}
					return
				}

				fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

				cmd := params.EksctlCreateCmd.WithArgs(
					"cluster",
					"--verbose", "4",
					"--name", noInstallCluster,
					"--tags", "alpha.eksctl.io/description=eksctl integration test",
					"--install-neuron-plugin=false",
					"--nodegroup-name", initNG,
					"--node-labels", "ng-name="+initNG,
					"--nodes", "1",
					"--node-type", "inf1.xlarge",
					"--version", params.Version,
					"--kubeconfig", params.KubeconfigPath,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should not have installed the neuron device plugin", func() {
				_, err := newClientSet().AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).Should(BeNotFoundError())
			})

			When("adding a nodegroup by default", func() {
				It("should install without error", func() {
					cmd := params.EksctlCreateCmd.WithArgs(
						"nodegroup",
						"--cluster", noInstallCluster,
						"--nodes", "1",
						"--verbose", "4",
						"--name", newNG,
						"--tags", "alpha.eksctl.io/description=eksctl integration test",
						"--node-labels", "ng-name="+newNG,
						"--nodes", "1",
						"--node-type", "inf1.xlarge",
						"--version", params.Version,
					)
					Expect(cmd).To(RunSuccessfully())
				})
				It("should install the neuron device plugin", func() {
					_, err := newClientSet().AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})
})

func newClientSet() *kubernetes.Clientset {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   defaultCluster,
			Region: params.Region,
		},
	}
	ctl, err := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(cfg)
	Expect(err).ShouldNot(HaveOccurred())

	clientSet, err := ctl.NewStdClientSet(cfg)
	Expect(err).ShouldNot(HaveOccurred())
	return clientSet
}
