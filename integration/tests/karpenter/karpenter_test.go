//go:build integration
// +build integration

package karpenter

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	testing.Init()
	if err := api.Register(); err != nil {
		panic(fmt.Errorf("unexpected error registering API scheme: %w", err))
	}
	params = tests.NewParams("karp")
}

func TestKarpenter(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Karpenter", func() {
	Context("Creating a cluster with Karpenter", func() {
		It("should support karpenter", func() {
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file=-",
					"--verbose=4",
					"--kubeconfig", params.KubeconfigPath,
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/cluster-config.yaml"))
			// For dumping information, we need the kubeconfig. We just log the error here to
			// know that it failed for debugging then carry on.
			session := cmd.Run()
			if session.ExitCode() != 0 {
				fmt.Fprintf(GinkgoWriter, "cluster create command failed:")
				fmt.Fprint(GinkgoWriter, string(session.Out.Contents()))
				fmt.Fprint(GinkgoWriter, string(session.Err.Contents()))
			}
			if session.ExitCode() != 0 {
				describeKarpenterResources([]string{"karpenter-webhook", "karpenter-controller"})
			}

			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			// Check webhook pod
			Expect(kubeTest.WaitForPodsReady(karpenter.DefaultNamespace, metav1.ListOptions{
				LabelSelector: "karpenter=webhook",
			}, 1, 10*time.Minute)).To(Succeed())
			// Check controller pod
			Expect(kubeTest.WaitForPodsReady(karpenter.DefaultNamespace, metav1.ListOptions{
				LabelSelector: "karpenter=controller",
			}, 1, 10*time.Minute)).To(Succeed())
		})
	})
})

// not using kubeTest since kubeTest fatals on error, and we don't want that.
func describeKarpenterResources(names []string) {
	for _, name := range names {
		cmd := exec.Command("kubectl", "--kubeconfig", params.KubeconfigPath, "describe", "replicaset", name, "-n", karpenter.DefaultNamespace)
		output, err := cmd.Output()
		Expect(err).NotTo(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "describe replicaset %s", name)
		fmt.Fprint(GinkgoWriter, string(output))
		cmd = exec.Command("kubectl", "--kubeconfig", params.KubeconfigPath, "describe", "deployment", name, "-n", karpenter.DefaultNamespace)
		output, err = cmd.Output()
		Expect(err).NotTo(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "describe deployment %s", name)
		fmt.Fprint(GinkgoWriter, string(output))
	}
}
