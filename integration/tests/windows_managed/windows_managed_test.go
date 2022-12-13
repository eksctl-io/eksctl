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
	DefaultTimeOut = 45 * time.Minute
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
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = params.ClusterName
	clusterConfig.Metadata.Version = api.DefaultVersion
	clusterConfig.Metadata.Region = api.DefaultRegion

	clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name: "linux",
			},
		},
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         "windows-mng",
				AMIFamily:    api.NodeImageFamilyWindowsServer2022FullContainer,
				VolumeSize:   aws.Int(120),
				InstanceType: "t3a.xlarge",
			},
		},
	}

	cmd := params.EksctlCreateCmd.
		WithArgs(
			"cluster",
			"--config-file", "-",
			"--verbose", "4",
			"--kubeconfig", params.KubeconfigPath,
		).
		WithoutArg("--region", params.Region).
		WithStdin(clusterutils.Reader(clusterConfig))
	Expect(cmd).To(RunSuccessfully())
})

var _ = Describe("(Integration) [EKS Windows Managed Nodegroups]", func() {
	Context("Create Windows managed nodegroup and launch pods", func() {
		It("should launch a Windows pod", func() {
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			d := kubeTest.CreateDeploymentFromFile("default", fmt.Sprintf("../../data/%s", "windows-server-iis-2022.yaml"))
			kubeTest.WaitForDeploymentReady(d, DefaultTimeOut)
		})
	})

	Context("should create Windows managed nodegroup with taints applied", func() {
		It("should not throw error when create Windows managed nodegroup with taints applied", func() {
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
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = params.ClusterName
			clusterConfig.Metadata.Region = params.Region
			clusterConfig.Metadata.Version = params.Version
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

	Context("should create Windows managed nodegroup via CLI", func() {
		It("should not throw error when create Windows managed nodegroup via CLI", func() {
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
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
