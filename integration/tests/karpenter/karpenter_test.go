//go:build integration
// +build integration

package karpenter

import (
	"fmt"
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
	params = tests.NewParams("karp")
}

func TestKarpenter(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Karpenter", func() {
	AfterEach(func() {
		cmd := params.EksctlDeleteCmd.WithArgs(
			"cluster", params.ClusterName,
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
					"--kubeconfig", params.KubeconfigPath,
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/cluster-config.yaml"))
			Expect(cmd).To(RunSuccessfully())

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
