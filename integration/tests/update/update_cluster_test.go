// +build integration

package update

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

const (
	k8sUpdatePollInterval = "2s"
	k8sUpdatePollTimeout  = "3m"
)

var defaultCluster string
var noInstallCluster string
var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("up")
	defaultCluster = params.ClusterName
	noInstallCluster = params.NewClusterName("update")
}

func TestUpdate(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Update addons", func() {
	const (
		initNG = "kp-ng-0"
	)
	BeforeSuite(func() {
		params.KubeconfigTemp = false
		if params.KubeconfigPath == "" {
			wd, _ := os.Getwd()
			f, _ := ioutil.TempFile(wd, "kubeconfig-")
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

		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", defaultCluster,
			"--tags", "alpha.eksctl.io/description=eksctl integration test",
			"--nodegroup-name", initNG,
			"--node-labels", "ng-name="+initNG,
			"--nodes", "1",
			"--node-type", "t3.large",
			"--version", "1.15",
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

	// Chose 1.15 because the upgrade to 1.16 takes less time than upgrading to 1.17
	Context("cluster with version 1.15", func() {
		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			awsSession := NewSession(params.Region)

			Expect(awsSession).To(HaveExistingCluster(params.ClusterName, awseks.ClusterStatusActive, api.Version1_15))

			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initNG)))
		})

		It("should upgrade the control plane to version 1.16", func() {

			cmd := params.EksctlUpgradeCmd.
				WithArgs(
					"cluster",
					"--verbose", "4",
					"--name", params.ClusterName,
					"--approve",
				)
			Expect(cmd).To(RunSuccessfully())

			By(fmt.Sprintf("checking that control plane is updated to %v", "1.16"))
			config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
			Expect(err).ToNot(HaveOccurred())

			clientSet, err := kubernetes.NewForConfig(config)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func() string {
				serverVersion, err := clientSet.ServerVersion()
				Expect(err).ToNot(HaveOccurred())
				return fmt.Sprintf("%s.%s", serverVersion.Major, strings.TrimSuffix(serverVersion.Minor, "+"))
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(Equal("1.16"))
		})

		It("should upgrade kube-proxy", func() {
			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-kube-proxy",
				"--cluster", params.ClusterName,
				"--verbose", "4",
				"--approve",
			)
			Expect(cmd).To(RunSuccessfully())

			clientSet := getClientSet()
			Eventually(func() string {
				daemonSet, err := clientSet.AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), "kube-proxy", metav1.GetOptions{})
				Expect(err).ToNot(HaveOccurred())
				kubeProxyVersion, err := addons.ImageTag(daemonSet.Spec.Template.Spec.Containers[0].Image)
				Expect(err).ToNot(HaveOccurred())
				return kubeProxyVersion
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(Equal("v1.16.15-eksbuild.1"))
		})

		It("should upgrade aws-node", func() {
			rawClient := getRawClient()
			preUpdateAWSNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), "aws-node", metav1.GetOptions{})
			Expect(err).ToNot(HaveOccurred())
			preUpdateAWSNodeVersion, err := addons.ImageTag(preUpdateAWSNode.Spec.Template.Spec.Containers[0].Image)
			Expect(err).ToNot(HaveOccurred())
			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-aws-node",
				"--cluster", params.ClusterName,
				"--verbose", "4",
				"--approve",
			)
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() string {
				awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), "aws-node", metav1.GetOptions{})
				Expect(err).ToNot(HaveOccurred())
				awsNodeVersion, err := addons.ImageTag(awsNode.Spec.Template.Spec.Containers[0].Image)
				Expect(err).ToNot(HaveOccurred())
				return awsNodeVersion
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).ShouldNot(Equal(preUpdateAWSNodeVersion))
		})

		It("should upgrade coredns", func() {
			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-coredns",
				"--cluster", params.ClusterName,
				"--verbose", "4",
				"--approve",
			)
			Expect(cmd).To(RunSuccessfully())

			rawClient := getRawClient()
			Eventually(func() string {
				coreDNSDeployment, err := rawClient.ClientSet().AppsV1().Deployments(metav1.NamespaceSystem).Get(context.TODO(), "coredns", metav1.GetOptions{})
				Expect(err).ToNot(HaveOccurred())
				coreDNSVersion, err := addons.ImageTag(coreDNSDeployment.Spec.Template.Spec.Containers[0].Image)
				Expect(err).ToNot(HaveOccurred())
				return coreDNSVersion
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(Equal("v1.6.6-eksbuild.1"))
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
	ctl, err := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(cfg)
	Expect(err).ShouldNot(HaveOccurred())
	rawClient, err := ctl.NewRawClient(cfg)
	Expect(err).ToNot(HaveOccurred())
	return rawClient
}

func getClientSet() *kubernetes.Clientset {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   params.ClusterName,
			Region: params.Region,
		},
	}
	ctl, err := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(cfg)
	Expect(err).ShouldNot(HaveOccurred())
	clientSet, err := ctl.NewStdClientSet(cfg)
	Expect(err).ToNot(HaveOccurred())
	return clientSet
}
