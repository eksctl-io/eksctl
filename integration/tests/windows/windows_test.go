//go:build integration

package windows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
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
	params = tests.NewParams(fmt.Sprintf("windows-%d", GinkgoParallelProcess()))
})

var _ = Describe("(Integration) [Windows Nodegroups]", func() {

	createCluster := func(withOIDC bool, ami, clusterName string) {
		By("creating a new cluster with Windows nodegroups")
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = clusterName
		clusterConfig.Metadata.Version = api.LatestVersion
		clusterConfig.Metadata.Region = api.DefaultRegion
		clusterConfig.IAM.WithOIDC = &withOIDC

		clusterConfig.NodeGroups = []*api.NodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "windows",
					AMIFamily:    ami,
					VolumeSize:   aws.Int(120),
					InstanceType: "t3a.xlarge",
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

	runWindowsPod := func(workload string) {
		By("scheduling a Windows pod")
		kubeTest, err := kube.NewTest(params.KubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

		d := kubeTest.CreateDeploymentFromFile("default", fmt.Sprintf("../../data/%s", workload))
		kubeTest.WaitForDeploymentReady(d, 45*time.Minute)
	}

	deleteCluster := func(clusterName string) {
		By("deleting the Windows cluster")
		cmd := params.EksctlDeleteCmd.WithArgs(
			"cluster",
			"--name", clusterName,
		)
		Expect(cmd).To(RunSuccessfully())
	}

	Context("When creating a cluster with Windows nodegroups", func() {
		DescribeTable("it should be able to run Windows pods", func(withOIDC bool, ami, workload, clusterName string) {
			createCluster(withOIDC, ami, clusterName)
			runWindowsPod(workload)
			deleteCluster(clusterName)
		},
			Entry("Windows when withOIDC is disabled", false, api.NodeImageFamilyWindowsServer2019FullContainer, "windows-server-iis.yaml", "windows-cluster1"),
			Entry("Windows when withOIDC is enabled", true, api.NodeImageFamilyWindowsServer2019FullContainer, "windows-server-iis.yaml", "windows-cluster2"),
			Entry("Windows Server 2022 when withOIDC is enabled", true, api.NodeImageFamilyWindowsServer2022FullContainer, "windows-server-iis-2022.yaml", "windows-cluster3"),
		)
	})

})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
