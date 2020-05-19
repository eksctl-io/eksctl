// +build integration

package inferentia

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

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

func TestSuite(t *testing.T) {
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
				cfg := &api.ClusterConfig{
					Metadata: &api.ClusterMeta{
						Name:   defaultCluster,
						Region: params.Region,
					},
				}
				ctl := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
				err := ctl.RefreshClusterStatus(cfg)
				Expect(err).ShouldNot(HaveOccurred())

				clientSet, err := ctl.NewStdClientSet(cfg)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = clientSet.AppsV1().DaemonSets("kube-system").Get("neuron-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).ShouldNot(HaveOccurred())
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
				cfg := &api.ClusterConfig{
					Metadata: &api.ClusterMeta{
						Name:   noInstallCluster,
						Region: params.Region,
					},
				}
				ctl := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
				err := ctl.RefreshClusterStatus(cfg)
				Expect(err).ShouldNot(HaveOccurred())

				clientSet, err := ctl.NewStdClientSet(cfg)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = clientSet.AppsV1().DaemonSets("kube-system").Get("neuron-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).Should(BeNotFoundError())
			})
		})
	})
})
