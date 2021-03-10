// +build integration

package addons

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/weaveworks/eksctl/integration/utilities/unowned"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	params  *tests.Params
	cluster *unowned.Cluster
)

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
		BeforeSuite(func() {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = params.ClusterName
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

			if params.UnownedCluster {
				cluster = unowned.NewCluster(clusterConfig)
				cluster.CreateNodegroups("ng")
			} else {
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
			}

		})

		AfterSuite(func() {
			cmd := params.EksctlDeleteCmd.WithArgs(
				"cluster", params.ClusterName,
				"--wait",
				"--verbose", "2",
			)
			Expect(cmd).To(RunSuccessfully())

			if params.UnownedCluster {
				cluster.DeleteStack()
			}
		})

		It("should support addons", func() {
			By("Asserting the addon is listed in `get addons`")
			//its created as part of create cluster for owned clusters
			if params.UnownedCluster {
				cmd := params.EksctlCreateCmd.
					WithArgs(
						"addon",
						"--cluster", params.ClusterName,
						"--name", "vpc-cni",
						"--verbose", "2",
					)
				Expect(cmd).To(RunSuccessfully())
			}
			cmd := params.EksctlGetCmd.
				WithArgs(
					"addons",
					"--cluster", params.ClusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("vpc-cni")),
			))

			Eventually(func() string {
				cmd = params.EksctlGetCmd.
					WithArgs(
						"addons",
						"--cluster", params.ClusterName,
						"--verbose", "2",
					)
				return string(cmd.Run().Out.Contents())
			}, time.Minute*5, time.Second*30).Should(ContainSubstring("ACTIVE"))

			By("Updating the addon")
			cmd = params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--name", "vpc-cni",
					"--cluster", params.ClusterName,
					"--wait",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			By("Deleting the addon")
			cmd = params.EksctlDeleteCmd.
				WithArgs(
					"addon",
					"--name", "vpc-cni",
					"--cluster", params.ClusterName,
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
