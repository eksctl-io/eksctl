//go:build integration
// +build integration

//revive:disable Not changing package name
package bare_cluster

import (
	"context"
	"testing"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	testing.Init()
	params = tests.NewParams("bare-cluster")
}

func TestBareCluster(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("Bare Clusters", Ordered, func() {
	var clusterConfig *api.ClusterConfig

	BeforeAll(func() {
		By("creating a cluster with only VPC CNI and no other default networking addons")
		clusterConfig = api.NewClusterConfig()
		clusterConfig.Metadata.Name = params.ClusterName
		clusterConfig.Metadata.Region = params.Region
		clusterConfig.AddonsConfig.DisableDefaultAddons = true
		clusterConfig.Addons = []*api.Addon{
			{
				Name:                              "vpc-cni",
				UseDefaultPodIdentityAssociations: true,
			},
			{
				Name: "eks-pod-identity-agent",
			},
		}
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--config-file=-",
				"--verbose", "4",
				"--kubeconfig="+params.KubeconfigPath,
			).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.Reader(clusterConfig))

		Expect(cmd).To(RunSuccessfully())
	})

	It("should have only VPC CNI installed", func() {
		config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
		_, err = clientset.AppsV1().Deployments(metav1.NamespaceSystem).Get(context.Background(), "coredns", metav1.GetOptions{})
		Expect(apierrors.IsNotFound(err)).To(BeTrue(), "expected coredns to not exist")
		daemonSets := clientset.AppsV1().DaemonSets(metav1.NamespaceSystem)
		_, err = daemonSets.Get(context.Background(), "kube-proxy", metav1.GetOptions{})
		Expect(apierrors.IsNotFound(err)).To(BeTrue(), "expected kube-proxy to not exist")
		_, err = daemonSets.Get(context.Background(), "aws-node", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred(), "expected aws-node to exist")
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
