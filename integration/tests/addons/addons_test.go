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
			clusterConfig.Metadata.Version = "1.19"
			clusterConfig.Metadata.Region = params.Region
			clusterConfig.IAM.WithOIDC = api.Enabled()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:             "vpc-cni",
					AttachPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
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
			))

			By("Updating the addon")
			cmd = params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--name", "vpc-cni",
					"--cluster", clusterName,
					"--wait",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			By("Deleting the addon")
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

	It("should describe addons", func() {
		cmd := params.EksctlUtilsCmd.
			WithArgs(
				"describe-addon-versions",
				"--kubernetes-version", "1.19",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("vpc-cni")),
		))
	})

})
