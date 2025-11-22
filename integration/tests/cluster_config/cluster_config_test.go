//go:build integration

//revive:disable Not changing package name
package cluster_config

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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
)

var eksAPI awsapi.EKS

var _ = BeforeSuite(func() {
	if params.SkipCreate {
		return
	}
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
	cmd := params.EksctlCreateCmd.WithArgs(
		"cluster",
		"--config-file", "-",
		"--verbose", "4",
	).
		WithoutArg("--region", params.Region).
		WithStdin(clusterutils.Reader(clusterConfig))

	Expect(cmd).To(RunSuccessfully())

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
