package karpenter

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	eventsv1 "k8s.io/api/events/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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
				dumpEventsAndLogs([]string{"karpenter=webhook", "karpenter=controller"})
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
func dumpEventsAndLogs(selectors []string) {
	config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
	Expect(err).NotTo(HaveOccurred())
	clientset, err := kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())
	// list all events in the Karpenter namespace
	events := clientset.EventsV1().Events(karpenter.DefaultNamespace)
	wes, err := events.List(context.Background(), metav1.ListOptions{})
	Expect(err).NotTo(HaveOccurred())
	var webhookEvents []eventsv1.Event
	webhookEvents = append(webhookEvents, wes.Items...)
	for wes.Continue != "" {
		wes, err = events.List(context.Background(), metav1.ListOptions{
			Continue: wes.Continue,
		})
		Expect(err).NotTo(HaveOccurred())
		webhookEvents = append(webhookEvents, wes.Items...)
	}
	fmt.Fprintf(GinkgoWriter, "all events:\n")
	for _, event := range webhookEvents {
		fmt.Fprintf(GinkgoWriter, "name: %s, message: %s\n", event.Name, event.Note)
	}

	for _, selector := range selectors {
		pods := clientset.CoreV1().Pods(karpenter.DefaultNamespace)
		ps, err := pods.List(context.Background(), metav1.ListOptions{
			LabelSelector: selector,
		})
		Expect(err).NotTo(HaveOccurred())
		for _, pod := range ps.Items {
			containerName := pod.Spec.Containers[0].Name
			logs, err := clientset.CoreV1().RESTClient().Get().
				Resource("pods").
				Namespace(pod.Namespace).
				Name(pod.Name).SubResource("log").
				Param("container", containerName).
				Stream(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			content, err := io.ReadAll(logs)
			Expect(err).NotTo(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "container logs for container %s for pod %s: %s\n", containerName, pod.Name, string(content))
		}
	}
}
