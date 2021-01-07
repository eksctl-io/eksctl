// +build integration

package unowned_clusters

import (
	"testing"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("e2e")
}

func TestE2E(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [non-eksctl created cluster & nodegroup support]", func() {
	Context("Get, upgrade & delete cluster", func() {
		var (
			clusterName, ng1, ng2 string
		)

		BeforeEach(func() {
			//create the cluster and nodegroup, v1.17
			//clusterName := params.NewClusterName("unowned")
			clusterName = "jk-console"
			ng1 = "ng-1"
			ng2 = "ng-2"
		})

		It("should work", func() {
			By("Getting clusters")
			cmd := params.EksctlGetCmd.
				WithArgs(
					"clusters",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(clusterName)),
			))

			By("Getting nodegroups")
			cmd = params.EksctlGetCmd.
				WithArgs(
					"nodegroups",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("ng-1")),
			))
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("ng-2")),
			))

			By("Enabling OIDC")
			cmd = params.EksctlUtilsCmd.
				WithArgs(
					"associate-iam-oidc-provider",
					"--name", clusterName,
					"--approve",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("vpc-cni")),
			))
			By("Creating an addon")
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"addons",
					"--cluster", clusterName,
					"--name", "vpc-cni",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("vpc-cni")),
			))

			By("Getting an addon")
			cmd = params.EksctlGetCmd.
				WithArgs(
					"addons",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("vpc-cni")),
			))

			By("Creating an IAMServiceAccount")
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"iamserviceaccount",
					"--cluster", clusterName,
					"--name", "test-sa",
					"--namespace", "default",
					"--attach-policy-arn",
					"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
					"--approve",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())
			By("Getting an IAMServiceAccount")
			cmd = params.EksctlGetCmd.
				WithArgs(
					"iamserviceaccount",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("test-sa")),
			))

			By("Upgrading the cluster")
			cmd = params.EksctlUpgradeCmd.
				WithArgs(
					"cluster",
					"--name", "vpc-cni",
					"--cluster", clusterName,
					"--version", "1.18",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			By("Upgrading one of the nodegroups")
			cmd = params.EksctlUpgradeCmd.
				WithArgs(
					"addon",
					"--name", ng1,
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			By("Deleting a nodegroup")
			cmd = params.EksctlDeleteCmd.
				WithArgs(
					"nodegroup",
					"--name", ng2,
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			By("Deleting the cluster")
			cmd = params.EksctlDeleteCmd.
				WithArgs(
					"addon",
					"--name", "vpc-cni",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())
		})
	})
})
