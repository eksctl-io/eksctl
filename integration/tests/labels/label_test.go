//go:build integration
// +build integration

package labels

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("labels")
	if err := api.Register(); err != nil {
		panic("unexpected error registering API scheme")
	}
}

func TestLabels(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("Labels", func() {
	var (
		mng1 string
	)
	BeforeSuite(func() {
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--config-file=-",
				"--verbose=4",
			).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/managed-nodegroup-with-labels.yaml"))
		Expect(cmd).To(RunSuccessfully())
		// corresponds to the label name in the cluster config file
		mng1 = "mng-1"
	})

	AfterSuite(func() {
		params.DeleteClusters()
	})

	It("supports labels", func() {
		By("getting existing labels")
		cmd := params.EksctlGetCmd.
			WithArgs(
				"labels",
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--verbose", "2",
			)
		// It sometimes takes forever for the above set to take effect
		Eventually(func() *gbytes.Buffer { return cmd.Run().Out }, time.Minute*4).Should(gbytes.Say("preset=value"))

		By("setting labels on a managed nodegroup")
		cmd = params.EksctlSetLabelsCmd.
			WithArgs(
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--labels", "fantastic=zombieman",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("getting the newly set labels for a managed nodegroup")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"labels",
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--verbose", "2",
			)
		// It sometimes takes forever for the above set to take effect
		Eventually(func() *gbytes.Buffer { return cmd.Run().Out }, time.Minute*4).Should(gbytes.Say("fantastic=zombieman"))

		By("unsetting labels on a managed nodegroup")
		cmd = params.EksctlUnsetLabelsCmd.
			WithArgs(
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--labels", "fantastic",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})
})
