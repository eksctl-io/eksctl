//go:build integration

package windows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
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
	params = tests.NewParams("windows")
})

var _ = Describe("(Integration) [Windows Nodegroups]", func() {

	createCluster := func(withOIDC bool, ami, containerRuntime string) {
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
					AMIFamily: ami,
				},
				ContainerRuntime: &containerRuntime,
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
		kubeTest.WaitForDeploymentReady(d, 12*time.Minute)
	}

	Context("windows with OIDC disabled", func() {
		BeforeEach(func() {
			createCluster(false, api.NodeImageFamilyWindowsServer2019FullContainer, api.ContainerRuntimeDockerForWindows)
		})
		It("should be able to run Windows pods", func() {
			runWindowsPod("windows-server-iis.yaml")
		})
	})

	Context("windows with OIDC enabled", func() {
		BeforeEach(func() {
			createCluster(true, api.NodeImageFamilyWindowsServer2019FullContainer, api.ContainerRuntimeDockerForWindows)
		})
		It("should be able to run Windows pods", func() {
			runWindowsPod("windows-server-iis.yaml")
		})
	})

	Context("windows with 20H2", func() {
		BeforeEach(func() {
			createCluster(true, api.NodeImageFamilyWindowsServer20H2CoreContainer, api.ContainerRuntimeContainerD)
		})
		It("should be able to run Windows pods", func() {
			runWindowsPod("windows-server-iis-20H2.yaml")
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
