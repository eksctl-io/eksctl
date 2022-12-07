//go:build integration

package windows_managed

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
	"time"
)

const (
	DefaultTimeOut           = 25 * time.Minute
	InitialAl2Nodegroup      = "al2-1"
	Windows2019FullNodegroup = "WindowsServer2019FullContainer"
	Windows2022FullNodegroup = "WindowsServer2022FullContainer"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("ManagedWindows")
}

func TestWindowsManaged(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

	cmd := params.EksctlCreateCmd.WithArgs(
		"cluster",
		"--verbose", "4",
		"--name", params.ClusterName,
		"--tags", "alpha.eksctl.io/description=eksctl integration test",
		"--managed",
		"--nodegroup-name", InitialAl2Nodegroup,
		"--node-labels", "ng-name="+InitialAl2Nodegroup,
		"--nodes", "2",
		"--instance-types", "t3a.xlarge",
		"--version", params.Version,
		"--kubeconfig", params.KubeconfigPath,
	)
	Expect(cmd).To(RunSuccessfully())
})

var _ = Describe("(Integration) [EKS Windows Managed Nodegroups]", func() {
	makeClusterConfig := func(clusterName string) *api.ClusterConfig {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = clusterName
		clusterConfig.Metadata.Region = params.Region
		clusterConfig.Metadata.Version = params.Version
		return clusterConfig
	}

	createClusterWithWindowsMNG := func(clusterName, amiType string) {
		clusterConfig := makeClusterConfig(clusterName)
		clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "windows",
					AMIFamily:    amiType,
					VolumeSize:   aws.Int(80),
					InstanceType: "t3a.xlarge",
				},
			},
		}

		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--verbose", "4",
				"--config-file", "-",
			).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.Reader(clusterConfig))
		Expect(cmd).To(RunSuccessfully())
	}

	runWindowsPod := func(workload string) {
		By("scheduling a Windows pod")
		kubeTest, err := kube.NewTest(params.KubeconfigPath)
		Expect(err).NotTo(HaveOccurred())

		d := kubeTest.CreateDeploymentFromFile("default", fmt.Sprintf("../../data/%s", workload))
		kubeTest.WaitForDeploymentReady(d, DefaultTimeOut)
	}

	deleteCluster := func(clusterName string) {
		By("deleting the windows cluster")
		cmd := params.EksctlDeleteCmd.WithArgs(
			"cluster",
			"--name", clusterName,
		)
		Expect(cmd).To(RunSuccessfully())
	}

	Context("Create Windows managed nodegroup with taints", func() {
		It("should create Windows managed nodegroup taints applied", func() {
			taints := []api.NodeGroupTaint{
				{
					Key:    "key1",
					Value:  "value1",
					Effect: "NoSchedule",
				},
				{
					Key:    "key2",
					Effect: "NoSchedule",
				},
				{
					Key:    "key3",
					Value:  "value2",
					Effect: "NoExecute",
				},
			}
			clusterConfig := makeClusterConfig(params.ClusterName)
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         "windows-taints",
						AMIFamily:    api.NodeImageFamilyWindowsServer2022CoreContainer,
						VolumeSize:   aws.Int(120),
						InstanceType: "t3a.xlarge",
					},
					Taints: taints,
				},
			}

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"nodegroup",
					"--verbose", "4",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			Expect(cmd).To(RunSuccessfully())

			config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			clientset, err := kubernetes.NewForConfig(config)
			Expect(err).NotTo(HaveOccurred())

			mapTaints := func(taints []api.NodeGroupTaint) []corev1.Taint {
				var ret []corev1.Taint
				for _, t := range taints {
					ret = append(ret, corev1.Taint{
						Key:    t.Key,
						Value:  t.Value,
						Effect: t.Effect,
					})
				}
				return ret
			}
			tests.AssertNodeTaints(tests.ListNodes(clientset, "taints"), mapTaints(taints))
		})
	})
	Context("Create Windows managed nodegroup via CLI", func() {
		It("should not return error when creating Windows managed nodegroup via CLI", func() {
			cmd := params.EksctlCreateCmd.WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--nodes", "4",
				"--managed",
				"--instance-types", "t3a.xlarge",
				"--node-ami-family=WindowsServer2019CoreContainer",
			)
			Expect(cmd).To(RunSuccessfully())
		})
	})
	Context("Create Windows managed nodegroups and run pods then cleanup", func() {
		DescribeTable("it should be able to run Windows pods", func(windowsAmiType, workload, clusterName string) {
			createClusterWithWindowsMNG(clusterName, windowsAmiType)
			runWindowsPod(workload)
			deleteCluster(clusterName)
		},
			Entry("Windows Server 2019 Full", api.NodeImageFamilyWindowsServer2019CoreContainer, "windows-server-iis.yaml", "managed-windows-cluster1"),
			Entry("Windows Server 2022 Full", api.NodeImageFamilyWindowsServer2022FullContainer, "windows-server-iis-2022.yaml", "managed-windows-cluster2"),
		)
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
