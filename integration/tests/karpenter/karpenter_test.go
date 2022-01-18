package karpenter

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/weaveworks/eksctl/integration/runner"
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
	params = tests.NewParams("")
}

func TestKarpenter(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Karpenter", func() {

	var (
		clusterName string
	)

	BeforeEach(func() {
		// the randomly generated name we get usually makes one of the resources have a longer than 64 characters name
		clusterName = fmt.Sprintf("it-karpenter-%d", time.Now().Unix())
	})

	AfterEach(func() {
		cmd := params.EksctlDeleteCmd.WithArgs(
			"cluster", clusterName,
			"--verbose", "4",
		)
		Expect(cmd).To(RunSuccessfully())
	})

	Context("Creating a cluster with Karpenter", func() {
		It("should support karpenter", func() {
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file=-",
					"--verbose=4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.ReaderFromFile(clusterName, params.Region, "testdata/cluster-config.yaml"))
			// For dumping information, we need the kubeconfig. We just log the error here to
			// know that it failed for debugging then carry on.
			session := cmd.Run()
			if session.ExitCode() != 0 {
				fmt.Fprintf(GinkgoWriter, "cluster create command failed:")
				fmt.Fprint(GinkgoWriter, string(session.Out.Contents()))
				fmt.Fprint(GinkgoWriter, string(session.Err.Contents()))
			}
			cmd = params.EksctlUtilsCmd.WithArgs(
				"write-kubeconfig",
				"--verbose", "4",
				"--cluster", clusterName,
				"--kubeconfig", params.KubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())

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
		cmd := exec.Command("kubectl", "describe", "replicaset", name, "-n", karpenter.DefaultNamespace)
		output, err := cmd.Output()
		Expect(err).NotTo(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "describe replicaset %s", name)
		fmt.Fprint(GinkgoWriter, string(output))
		cmd = exec.Command("kubectl", "describe", "deployment", name, "-n", karpenter.DefaultNamespace)
		output, err = cmd.Output()
		Expect(err).NotTo(HaveOccurred())
		fmt.Fprintf(GinkgoWriter, "describe deployment %s", name)
		fmt.Fprint(GinkgoWriter, string(output))
	}
}
