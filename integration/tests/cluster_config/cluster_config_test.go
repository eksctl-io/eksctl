//go:build integration

//revive:disable Not changing package name
package cluster_config

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("cluster-config")
}

func TestClusterConfig(t *testing.T) {
	testutils.RegisterAndRun(t)
}

const (
	expectedControlPlaneTier = "tier-xl"
	expectedSupportType      = api.SupportTypeStandard

	// Values set on create, then changed by `eksctl utils update-control-plane-component-config`
	// below to assert the update path end-to-end.
	expectedEventTTL        = "45m"
	expectedUpdatedEventTTL = "50m"
	expectedMinPort         = 30000
	expectedMaxPort         = 32767
	expectedScoringStrategy = "MostAllocated"
	// Valid range on a tiered control plane is 10s-15s; 15s is the default, so use 10s to
	// assert the value actually took effect.
	expectedHPASyncPeriod = "10s"
)

var eksAPI awsapi.EKS

var _ = BeforeSuite(func() {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = params.ClusterName
	clusterConfig.Metadata.Region = params.Region
	clusterConfig.Metadata.Version = params.Version
	clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{}
	clusterConfig.UpgradePolicy = &api.UpgradePolicy{
		SupportType: expectedSupportType,
	}
	clusterConfig.ControlPlaneScalingConfig = &api.ControlPlaneScalingConfig{
		Tier: aws.String(expectedControlPlaneTier),
	}
	clusterConfig.KubeAPIServerConfig = &api.KubeAPIServerConfig{
		EventTTL: aws.String(expectedEventTTL),
		ServiceNodePortRange: &api.ServiceNodePortRange{
			MinPort: aws.Int(expectedMinPort),
			MaxPort: aws.Int(expectedMaxPort),
		},
	}
	clusterConfig.KubeSchedulerConfig = &api.KubeSchedulerConfig{
		NodeResourcesFit: &api.NodeResourcesFitConfig{
			ScoringStrategy: &api.ScoringStrategy{
				Type: aws.String(expectedScoringStrategy),
			},
		},
	}
	clusterConfig.KubeControllerManagerConfig = &api.KubeControllerManagerConfig{
		HorizontalPodAutoscalerControllerConfig: &api.HorizontalPodAutoscalerControllerConfig{
			HorizontalPodAutoscalerSyncPeriod: aws.String(expectedHPASyncPeriod),
		},
	}
	if !params.SkipCreate {
		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--config-file", "-",
			"--verbose", "4",
		).
			WithoutArg("--region", params.Region).
			WithStdin(clusterutils.Reader(clusterConfig))

		Expect(cmd).To(RunSuccessfully())
	}

	// Initialised even when create is skipped, so that the suite can be re-run against an
	// existing cluster via -eksctl.skip.create.
	clusterProvider, err := eks.New(context.Background(), &api.ProviderConfig{Region: params.Region}, clusterConfig)
	Expect(err).NotTo(HaveOccurred())
	eksAPI = clusterProvider.AWSProvider.EKS()
})

var _ = Describe("(Integration) [Cluster Config test]", func() {

	Context("Cluster with config options", func() {

		It("upgradePolicy should be set", func() {
			cluster, err := eksAPI.DescribeCluster(context.Background(), &awseks.DescribeClusterInput{
				Name: aws.String(params.ClusterName),
			})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			Expect(string(cluster.Cluster.UpgradePolicy.SupportType)).To(Equal(expectedSupportType))
		})

		It("control plane policy should be set", func() {
			cluster, err := eksAPI.DescribeCluster(context.Background(), &awseks.DescribeClusterInput{
				Name: aws.String(params.ClusterName),
			})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			Expect(string(cluster.Cluster.ControlPlaneScalingConfig.Tier)).To(Equal(expectedControlPlaneTier))
		})

		It("control plane component config should be set", func() {
			cluster, err := eksAPI.DescribeCluster(context.Background(), &awseks.DescribeClusterInput{
				Name: aws.String(params.ClusterName),
			})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			apiServerConfig := cluster.Cluster.KubeApiServerConfig
			Expect(apiServerConfig).NotTo(BeNil())
			Expect(apiServerConfig.EventTtl).To(HaveValue(Equal(expectedEventTTL)))
			Expect(apiServerConfig.ServiceNodePortRange).NotTo(BeNil())
			Expect(apiServerConfig.ServiceNodePortRange.MinPort).To(BeNumerically("==", expectedMinPort))
			Expect(apiServerConfig.ServiceNodePortRange.MaxPort).To(BeNumerically("==", expectedMaxPort))

			schedulerConfig := cluster.Cluster.KubeSchedulerConfig
			Expect(schedulerConfig).NotTo(BeNil())
			Expect(schedulerConfig.NodeResourcesFit).NotTo(BeNil())
			Expect(schedulerConfig.NodeResourcesFit.ScoringStrategy).NotTo(BeNil())
			Expect(string(schedulerConfig.NodeResourcesFit.ScoringStrategy.Type)).To(Equal(expectedScoringStrategy))

			controllerManagerConfig := cluster.Cluster.KubeControllerManagerConfig
			Expect(controllerManagerConfig).NotTo(BeNil())
			Expect(controllerManagerConfig.HorizontalPodAutoscalerControllerConfig).NotTo(BeNil())
			Expect(controllerManagerConfig.HorizontalPodAutoscalerControllerConfig.HorizontalPodAutoscalerSyncPeriod).
				To(HaveValue(Equal(expectedHPASyncPeriod)))
		})
	})

	Context("eksctl utils update-control-plane-component-config", Serial, func() {

		It("should reject a config file with no component values set", func() {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = params.ClusterName
			clusterConfig.Metadata.Region = params.Region
			// Passes the loader's "at least one component" check, but carries no values.
			clusterConfig.KubeAPIServerConfig = &api.KubeAPIServerConfig{}

			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-control-plane-component-config",
				"--config-file", "-",
				"--verbose", "4",
			).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			session := cmd.Run()
			Expect(session.ExitCode()).NotTo(Equal(0))
			Expect(string(session.Err.Contents()) + string(session.Buffer().Contents())).
				To(ContainSubstring("no control plane component config values are set; nothing to update"))
		})

		It("should update only the components present in the config file", func() {
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = params.ClusterName
			clusterConfig.Metadata.Region = params.Region
			// Only kube-apiserver is specified; the other two components must be left as
			// they were set on create, rather than cleared.
			clusterConfig.KubeAPIServerConfig = &api.KubeAPIServerConfig{
				EventTTL: aws.String(expectedUpdatedEventTTL),
			}

			// The command waits for the ControlPlaneComponentConfigUpdate to reach a terminal
			// state before returning, which takes several minutes. EksctlUtilsCmd's default
			// 5m timeout is not enough.
			cmd := params.EksctlUtilsCmd.WithArgs(
				"update-control-plane-component-config",
				"--config-file", "-",
				"--verbose", "4",
			).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig)).
				WithTimeout(30 * time.Minute)
			Expect(cmd).To(RunSuccessfully())

			cluster, err := eksAPI.DescribeCluster(context.Background(), &awseks.DescribeClusterInput{
				Name: aws.String(params.ClusterName),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(cluster.Cluster.KubeApiServerConfig.EventTtl).To(HaveValue(Equal(expectedUpdatedEventTTL)))

			By("leaving the components absent from the config file unchanged")
			Expect(cluster.Cluster.KubeApiServerConfig.ServiceNodePortRange).NotTo(BeNil())
			Expect(cluster.Cluster.KubeApiServerConfig.ServiceNodePortRange.MinPort).To(BeNumerically("==", expectedMinPort))
			Expect(string(cluster.Cluster.KubeSchedulerConfig.NodeResourcesFit.ScoringStrategy.Type)).To(Equal(expectedScoringStrategy))
			Expect(cluster.Cluster.KubeControllerManagerConfig.HorizontalPodAutoscalerControllerConfig.HorizontalPodAutoscalerSyncPeriod).
				To(HaveValue(Equal(expectedHPASyncPeriod)))
		})
	})

})

var _ = AfterSuite(func() {
	if params.SkipDelete {
		return
	}
	cmd := params.EksctlDeleteCmd.WithArgs(
		"cluster", params.ClusterName,
		"--disable-nodegroup-eviction",
		"--verbose", "2",
	)
	Expect(cmd).To(RunSuccessfully())
})
