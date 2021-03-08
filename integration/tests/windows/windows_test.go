// +build integration

package windows

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/weaveworks/eksctl/integration/utilities/unowned"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("windows")

}

func TestWindowsCluster(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	params         *tests.Params
	unownedCluster *unowned.Cluster
)

var _ = Describe("(Integration) [Windows Nodegroups]", func() {
	createCluster := func(withOIDC bool) {
		By("creating a new cluster with Windows nodegroups")
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = params.NewClusterName("windows")
		clusterConfig.Metadata.Version = api.LatestVersion
		clusterConfig.Metadata.Region = api.DefaultRegion
		clusterConfig.IAM.WithOIDC = &withOIDC

		clusterConfig.NodeGroups = []*api.NodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name:      "windows",
					AMIFamily: api.NodeImageFamilyWindowsServer2019FullContainer,
				},
			},
		}
		clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "linux",
				},
			},
		}
		if params.UnownedCluster {
			unownedCluster = unowned.NewCluster(clusterConfig)
			clusterConfig.VPC = unownedCluster.VPC
			cmd := params.EksctlUtilsCmd.WithArgs(
				"write-kubeconfig",
				"--verbose", "4",
				"--cluster", clusterConfig.Metadata.Name,
				"--kubeconfig", params.KubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())

			if withOIDC {
				cmd = params.EksctlUtilsCmd.
					WithArgs(
						"associate-iam-oidc-provider",
						"--name", clusterConfig.Metadata.Name,
						"--approve",
						"--verbose", "2",
					)
				Expect(cmd).To(RunSuccessfully())
			}

			cmd = params.EksctlUtilsCmd.WithArgs(
				"install-vpc-controllers",
				"--verbose", "4",
				"--approve",
				"--cluster", clusterConfig.Metadata.Name,
			)
			Expect(cmd).To(RunSuccessfully())

			data, err := json.Marshal(clusterConfig)
			Expect(err).ToNot(HaveOccurred())
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())
		} else {
			data, err := json.Marshal(clusterConfig)
			Expect(err).ToNot(HaveOccurred())
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--verbose", "4",
					"--kubeconfig", params.KubeconfigPath,
					"--install-vpc-controllers",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())
		}

	}

	runWindowsPod := func() {
		By("scheduling a Windows pod")
		kubeTest, err := kube.NewTest(params.KubeconfigPath)
		Expect(err).ToNot(HaveOccurred())

		d := kubeTest.CreateDeploymentFromFile("default", "../../data/windows-server-iis.yaml")
		kubeTest.WaitForDeploymentReady(d, 8*time.Minute)
	}

	AfterSuite(func() {
		params.DeleteClusters()
		if params.UnownedCluster {
			unownedCluster.DeleteStack()
		}
	})

	Context("When creating a cluster with Windows nodegroups", func() {
		DescribeTable("it should be able to run Windows pods", func(withOIDC bool) {
			createCluster(withOIDC)
			runWindowsPod()
		},
			Entry("when withOIDC is disabled", false),
			Entry("when withOIDC is enabled", true),
		)
	})

})
