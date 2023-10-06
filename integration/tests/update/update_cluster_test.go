//go:build integration
// +build integration

package update

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

const (
	k8sUpdatePollInterval = "30s"
	k8sUpdatePollTimeout  = "10m"
)

var (
	defaultCluster string
	params         *tests.Params
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("up")
	defaultCluster = params.ClusterName
}

func TestUpdate(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	eksVersion     string
	nextEKSVersion string
)

const (
	initNG = "kp-ng-0"
)

var _ = BeforeSuite(func() {
	params.KubeconfigTemp = false
	if params.KubeconfigPath == "" {
		wd, _ := os.Getwd()
		f, _ := os.CreateTemp(wd, "kubeconfig-")
		params.KubeconfigPath = f.Name()
		params.KubeconfigTemp = true
	}

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

	eksVersion, nextEKSVersion = clusterutils.GetCurrentAndNextVersionsForUpgrade(params.Version)

	cmd := params.EksctlCreateCmd.WithArgs(
		"cluster",
		"--verbose", "4",
		"--name", defaultCluster,
		"--tags", "alpha.eksctl.io/description=eksctl integration test",
		"--nodegroup-name", initNG,
		"--node-labels", "ng-name="+initNG,
		"--nodes", "1",
		"--node-type", "t3.large",
		"--version", eksVersion,
		"--kubeconfig", params.KubeconfigPath,
	)
	Expect(cmd).To(RunSuccessfully())
})
var _ = Describe("(Integration) Update addons", func() {

	Context("update cluster and addons", func() {
		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			config := NewConfig(params.Region)

			Expect(config).To(HaveExistingCluster(params.ClusterName, string(types.ClusterStatusActive), eksVersion))
			Expect(config).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(config).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initNG)))
		})

		It("should upgrade the control plane to the next version", func() {
			cmd := params.EksctlUpgradeCmd.
				WithArgs(
					"cluster",
					"--verbose", "4",
					"--name", params.ClusterName,
					"--approve",
				)
			Expect(cmd).To(RunSuccessfully())

			By(fmt.Sprintf("checking that control plane is updated to %v", nextEKSVersion))
			config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			clientSet, err := kubernetes.NewForConfig(config)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				serverVersion, err := clientSet.ServerVersion()
				Expect(err).NotTo(HaveOccurred())
				return fmt.Sprintf("%s.%s", serverVersion.Major, strings.TrimSuffix(serverVersion.Minor, "+"))
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(Equal(nextEKSVersion))
		})

		It("should upgrade kube-proxy", func() {
			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-kube-proxy",
				"--cluster", params.ClusterName,
				"--verbose", "4",
				"--approve",
			)
			Expect(cmd).To(RunSuccessfully())

			rawClient := getRawClient(context.Background())
			Eventually(func() string {
				daemonSet, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), "kube-proxy", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				kubeProxyVersion, err := addons.ImageTag(daemonSet.Spec.Template.Spec.Containers[0].Image)
				Expect(err).NotTo(HaveOccurred())
				v, err := version.NewVersion(kubeProxyVersion)
				Expect(err).NotTo(HaveOccurred())
				segments := v.Segments()
				Expect(len(segments)).To(BeNumerically(">=", 2))
				return fmt.Sprintf("%d.%d", segments[0], segments[1])
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(Equal(nextEKSVersion))
		})

		It("should upgrade aws-node", func() {
			rawClient := getRawClient(context.Background())
			getAWSNodeVersion := func() string {
				awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), "aws-node", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				imageTag, err := addons.ImageTag(awsNode.Spec.Template.Spec.Containers[0].Image)
				Expect(err).NotTo(HaveOccurred())
				return imageTag
			}
			preUpdateAWSNodeVersion := getAWSNodeVersion()

			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-aws-node",
				"--cluster", params.ClusterName,
				"--verbose", "4",
				"--approve",
			)
			Expect(cmd).To(RunSuccessfully())

			Eventually(getAWSNodeVersion, k8sUpdatePollTimeout, k8sUpdatePollInterval).ShouldNot(Equal(preUpdateAWSNodeVersion))
		})

		It("should upgrade coredns", func() {
			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-coredns",
				"--cluster", params.ClusterName,
				"--verbose", "4",
				"--approve",
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should upgrade the nodegroup to the next version", func() {
			cmd := params.EksctlUpgradeCmd.WithArgs(
				"nodegroup",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"--name", initNG,
				"--kubernetes-version", nextEKSVersion,
				"--timeout=60m", // wait for CF stacks to finish update
			)
			ExpectWithOffset(1, cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("nodegroup successfully upgraded")))
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

func getRawClient(ctx context.Context) *kubewrapper.RawClient {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   params.ClusterName,
			Region: params.Region,
		},
	}
	ctl, err := eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(ctx, cfg)
	Expect(err).ShouldNot(HaveOccurred())
	rawClient, err := ctl.NewRawClient(cfg)
	Expect(err).NotTo(HaveOccurred())
	return rawClient
}
