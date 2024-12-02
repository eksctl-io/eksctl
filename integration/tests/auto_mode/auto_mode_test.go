//go:build integration
// +build integration

//revive:disable Not changing package name
package auto_mode

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	testing.Init()
	params = tests.NewParams("auto-mode")
}

func TestAutoMode(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("Auto Mode", Ordered, func() {
	var clusterConfig *api.ClusterConfig
	var eksAPI awsapi.EKS
	describeComputeConfig := func() *ekstypes.ComputeConfigResponse {
		cluster, err := eksAPI.DescribeCluster(context.Background(), &awseks.DescribeClusterInput{
			Name: aws.String(clusterConfig.Metadata.Name),
		})
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		return cluster.Cluster.ComputeConfig
	}

	assertAutoMode := func(enabled bool) {
		cc := describeComputeConfig()
		Expect(*cc.Enabled).To(Equal(enabled), "expected computeConfig.enabled to be %v", enabled)
		if enabled {
			Expect(cc.NodePools).To(ConsistOf("general-purpose", "system"))
			Expect(*cc.NodeRoleArn).NotTo(BeEmpty(), "expected cc.nodeRoleArn to be non-empty")
		} else {
			Expect(cc.NodePools).To(BeEmpty())
			Expect(cc.NodeRoleArn).To(BeNil(), "expected cc.nodeRoleArn to be nil")
		}
	}

	BeforeAll(func() {
		By("creating a cluster with Auto Mode enabled")
		clusterConfig = api.NewClusterConfig()
		clusterConfig.Metadata.Name = params.ClusterName
		clusterConfig.Metadata.Region = params.Region
		clusterConfig.AutoModeConfig = &api.AutoModeConfig{
			Enabled: api.Enabled(),
		}
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"cluster",
				"--config-file=-",
				"--verbose", "4",
				"--kubeconfig="+params.KubeconfigPath,
			).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.Reader(clusterConfig))

		Expect(cmd).To(RunSuccessfully())
		ctl, err := eks.New(context.Background(), &api.ProviderConfig{Region: params.Region}, clusterConfig)
		Expect(err).NotTo(HaveOccurred())
		eksAPI = ctl.AWSProvider.EKS()
	})

	It("should have Auto Mode enabled", func() {
		assertAutoMode(true)
	})

	It("should schedule workloads on nodes launched by Auto Mode", func() {
		test, err := kube.NewTest(params.KubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		d := test.CreateDeploymentFromFile(test.Namespace, "../../data/podinfo.yaml")
		test.WaitForDeploymentReady(d, 30*time.Minute)
		deployment, err := test.GetDeployment(test.Namespace, "podinfo")
		Expect(err).NotTo(HaveOccurred())
		nodeList := test.ListNodes(metav1.ListOptions{})
		Expect(nodeList.Items).To(HaveLen(1))
		node := nodeList.Items[0]
		const labelName = "eks.amazonaws.com/compute-type"
		computeType, ok := node.Labels[labelName]
		Expect(ok).To(BeTrue(), "expected to find label %s on node %s", labelName, node.Name)
		Expect(computeType).To(Equal("auto"))
		podList := test.ListPodsFromDeployment(deployment)
		Expect(podList.Items).To(HaveLen(2))
		for _, pod := range podList.Items {
			Expect(node.Name).To(Equal(pod.Spec.NodeName))
		}
	})

	It("should disable and re-enable Auto Mode", func() {
		updateAutoMode := func(enabled bool) {
			clusterConfig.AutoModeConfig.Enabled = aws.Bool(enabled)
			cmd := params.EksctlUpdateCmd.
				WithArgs(
					"auto-mode-config",
					"--config-file=-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))

			Expect(cmd).To(RunSuccessfully())
		}
		By("disabling Auto Mode")
		updateAutoMode(false)
		assertAutoMode(false)
		By("enabling Auto Mode")
		updateAutoMode(true)
		assertAutoMode(true)
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
})
