// +build integration

package addons

import (
	"bytes"
	"encoding/json"
	"testing"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("addons")
}

func TestEKSAddons(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [EKS Addons test]", func() {

	Context("Creating a cluster with addons", func() {
		clusterName := params.NewClusterName("addons")

		BeforeSuite(func() {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = clusterName
			clusterConfig.Metadata.Version = api.LatestVersion
			clusterConfig.Metadata.Region = params.Region
			clusterConfig.IAM.WithOIDC = api.Enabled()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:             "vpc-cni",
					AttachPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
				},
				{
					Name:    "coredns",
					Version: "latest",
				},
			}

			ng := &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "ng",
				},
			}
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{ng}

			data, err := json.Marshal(clusterConfig)
			Expect(err).ToNot(HaveOccurred())

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

		})

		AfterSuite(func() {
			cmd := params.EksctlDeleteCmd.WithArgs(
				"cluster", clusterName,
				"--verbose", "2",
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should support addons", func() {
			By("Asserting the addon is listed in `get addons`")
			cmd := params.EksctlGetCmd.
				WithArgs(
					"addons",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("vpc-cni")),
				ContainElement(ContainSubstring("coredns")),
			))

			By("Asserting the addons are healthy")
			cmd = params.EksctlGetCmd.
				WithArgs(
					"addon",
					"--name", "vpc-cni",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			cmd = params.EksctlGetCmd.
				WithArgs(
					"addon",
					"--name", "coredns",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			By("successfully creating the kube-proxy addon")

			cmd = params.EksctlCreateCmd.
				WithArgs(
					"addon",
					"--name", "kube-proxy",
					"--cluster", clusterName,
					"--force",
					"--wait",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			cmd = params.EksctlGetCmd.
				WithArgs(
					"addon",
					"--name", "kube-proxy",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			By("Deleting the kube-proxy addon")
			cmd = params.EksctlDeleteCmd.
				WithArgs(
					"addon",
					"--name", "kube-proxy",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())
		})
	})

	It("should describe addons", func() {
		cmd := params.EksctlUtilsCmd.
			WithArgs(
				"describe-addon-versions",
				"--kubernetes-version", api.LatestVersion,
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("vpc-cni")),
		))
	})

})
