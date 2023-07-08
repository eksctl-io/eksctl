//go:build integration
// +build integration

package inferentia

import (
	"context"
	"os"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var (
	params                  *tests.Params
	clusterWithNeuronPlugin string
	clusterWithoutPlugin    string
	selectedNodeType        string
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("inf")
	// Hardcoding the region for this integration test, as inf2.xlarge instance type support is limited.
	params.Region = api.RegionUSEast2
	clusterWithNeuronPlugin = params.ClusterName
	clusterWithoutPlugin = params.NewClusterName("inf-no-plugin")

	selectedNodeType = "inf2.xlarge"
	currentDay := time.Now().Day()
	if currentDay%2 == 0 {
		selectedNodeType = "inf1.xlarge"
	}
}

func TestInferentia(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const initNG = "inf-ng-0"

var _ = BeforeSuite(func() {
	params.KubeconfigTemp = false
	if params.KubeconfigPath == "" {
		wd, _ := os.Getwd()
		f, _ := os.CreateTemp(wd, "kubeconfig-")
		params.KubeconfigPath = f.Name()
		params.KubeconfigTemp = true
	}

	if !params.SkipCreate {
		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", clusterWithoutPlugin,
			"--tags", "alpha.eksctl.io/description=eksctl integration test",
			"--install-neuron-plugin=false",
			"--nodegroup-name", initNG,
			"--node-labels", "ng-name="+initNG,
			"--nodes", "1",
			"--node-type", selectedNodeType,
			"--version", params.Version,
			"--region", params.Region,
			"--zones", "us-east-2a,us-east-2b",
			"--kubeconfig", params.KubeconfigPath,
		)
		Expect(cmd).To(RunSuccessfully())

		cmd = params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", clusterWithNeuronPlugin,
			"--tags", "alpha.eksctl.io/description=eksctl integration test",
			"--nodegroup-name", initNG,
			"--node-labels", "ng-name="+initNG,
			"--nodes", "1",
			"--node-type", selectedNodeType,
			"--version", params.Version,
			"--region", params.Region,
			"--zones", "us-east-2a,us-east-2b",
			"--kubeconfig", params.KubeconfigPath,
		)
		Expect(cmd).To(RunSuccessfully())
	}
})

var _ = Describe("(Integration) Inferentia nodes", func() {
	const (
		newNG = "inf-ng-1"
	)

	Context("cluster with inf nodes", func() {
		Context("by default", func() {
			BeforeEach(func() {
				cmd := params.EksctlUtilsCmd.WithArgs(
					"write-kubeconfig",
					"--verbose", "4",
					"--region", params.Region,
					"--cluster", clusterWithoutPlugin,
					"--kubeconfig", params.KubeconfigPath,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should have installed the neuron device plugin", func() {
				clientSet := newClientSet(clusterWithNeuronPlugin)
				_, err := clientSet.AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should not have installed the nvidia device plugin", func() {
				_, err := newClientSet(clusterWithNeuronPlugin).AppsV1().DaemonSets("kube-system").Get(context.TODO(), "nvidia-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).Should(BeNotFoundError())
			})
		})

		Context("with --install-neuron-plugin=false", func() {
			BeforeEach(func() {
				cmd := params.EksctlUtilsCmd.WithArgs(
					"write-kubeconfig",
					"--verbose", "4",
					"--region", params.Region,
					"--cluster", clusterWithoutPlugin,
					"--kubeconfig", params.KubeconfigPath,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should not have installed the neuron device plugin", func() {
				_, err := newClientSet(clusterWithoutPlugin).AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
				Expect(err).Should(BeNotFoundError())
			})

			When("adding an unmanaged nodegroup by default", func() {
				params.LogStacksEventsOnFailureForCluster(clusterWithoutPlugin)
				It("should install without error", func() {
					cmd := params.EksctlCreateCmd.WithArgs(
						"nodegroup",
						"--region", params.Region,
						"--cluster", clusterWithoutPlugin,
						"--managed=false",
						"--nodes", "1",
						"--verbose", "4",
						"--name", newNG,
						"--tags", "alpha.eksctl.io/description=eksctl integration test",
						"--node-labels", "ng-name="+newNG,
						"--nodes", "1",
						"--node-type", selectedNodeType,
						"--version", params.Version,
					)
					Expect(cmd).To(RunSuccessfully())
				})
				It("should install the neuron device plugin", func() {
					_, err := newClientSet(clusterWithoutPlugin).AppsV1().DaemonSets("kube-system").Get(context.TODO(), "neuron-device-plugin-daemonset", metav1.GetOptions{})
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
	gexec.KillAndWait()
	if params.KubeconfigTemp {
		os.Remove(params.KubeconfigPath)
	}
	os.RemoveAll(params.TestDirectory)
})

func newClientSet(name string) kubernetes.Interface {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   name,
			Region: params.Region,
		},
	}
	ctx := context.Background()
	ctl, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(ctx, cfg)
	Expect(err).ShouldNot(HaveOccurred())

	clientSet, err := ctl.NewStdClientSet(cfg)
	Expect(err).ShouldNot(HaveOccurred())
	return clientSet
}
