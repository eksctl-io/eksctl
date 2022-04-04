//go:build integration
// +build integration

package update

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/weaveworks/eksctl/pkg/eks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

const (
	k8sUpdatePollInterval = "2s"
	k8sUpdatePollTimeout  = "3m"
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

var _ = Describe("(Integration) Update addons", func() {
	const (
		initNG = "kp-ng-0"
	)
	var (
		eksVersion     string
		nextEKSVersion string
	)

	BeforeSuite(func() {
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

		supportedVersions := api.SupportedVersions()
		if len(supportedVersions) < 2 {
			Fail("Update cluster test requires at least two supported EKS versions")
		}

		// Use the lowest supported version
		eksVersion, nextEKSVersion = supportedVersions[0], supportedVersions[1]

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

	AfterSuite(func() {
		params.DeleteClusters()
		gexec.KillAndWait()
		if params.KubeconfigTemp {
			os.Remove(params.KubeconfigPath)
		}
		os.RemoveAll(params.TestDirectory)
	})

	Context(fmt.Sprintf("cluster with version %s", eksVersion), func() {
		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			awsSession := NewSession(params.Region)

			Expect(awsSession).To(HaveExistingCluster(params.ClusterName, awseks.ClusterStatusActive, eksVersion))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initNG)))
		})

		It(fmt.Sprintf("should upgrade the control plane to version %s", nextEKSVersion), func() {

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

			rawClient := getRawClient()
			kubernetesVersion, err := rawClient.ServerVersion()
			Expect(err).NotTo(HaveOccurred())
			Eventually(func() string {
				daemonSet, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), "kube-proxy", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				kubeProxyVersion, err := addons.ImageTag(daemonSet.Spec.Template.Spec.Containers[0].Image)
				Expect(err).NotTo(HaveOccurred())
				return kubeProxyVersion
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(ContainSubstring(fmt.Sprintf("v%s-eksbuild.", kubernetesVersion)))
		})

		It("should upgrade aws-node", func() {
			rawClient := getRawClient()
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

	})
})

func getRawClient() *kubewrapper.RawClient {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   params.ClusterName,
			Region: params.Region,
		},
	}
	ctl, err := eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(cfg)
	Expect(err).ShouldNot(HaveOccurred())
	rawClient, err := ctl.NewRawClient(cfg)
	Expect(err).NotTo(HaveOccurred())
	return rawClient
}
