//go:build integration
// +build integration

package windows

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

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
}

func TestWindowsCluster(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var params *tests.Params

var _ = BeforeSuite(func() {
	params = tests.NewParams("windows")
})

var _ = Describe("(Integration) [Windows Nodegroups]", func() {

	createCluster := func(withOIDC bool) {
		By("creating a new cluster with Windows nodegroups")
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = params.NewClusterName("windows")
		clusterConfig.Metadata.Version = api.DefaultVersion
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

		data, err := json.Marshal(clusterConfig)
		Expect(err).NotTo(HaveOccurred())

		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--config-file", "-",
				"--verbose", "4",
				"--kubeconfig", params.KubeconfigPath,
			).
			WithoutArg("--region", params.Region).
			WithStdin(bytes.NewReader(data))
		Expect(cmd).To(RunSuccessfully())
	}

	runWindowsPod := func() {
		By("scheduling a Windows pod")
		kubeTest, err := kube.NewTest(params.KubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

		d := kubeTest.CreateDeploymentFromFile("default", "../../data/windows-server-iis.yaml")
		kubeTest.WaitForDeploymentReady(d, 12*time.Minute)
	}

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

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
