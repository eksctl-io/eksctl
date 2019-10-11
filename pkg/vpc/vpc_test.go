package vpc

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/eks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type endpointAccessCase struct {
	ClusterName string
	Private     bool
	Public      bool
	DCOutput    *eks.DescribeClusterOutput
	Error       error
}

var (
	cluster *eks.Cluster
	cfg     *api.ClusterConfig
	p       *mockprovider.MockProvider
)

func NewFakeClusterWithEndpoints(private, public bool, name string) *eks.Cluster {
	cluster := NewFakeCluster(name, eks.ClusterStatusActive)
	vpcCfgReq := eks.VpcConfigResponse{}
	vpcCfgReq.SetEndpointPrivateAccess(private).SetEndpointPublicAccess(public)
	cluster.SetResourcesVpcConfig(&vpcCfgReq)
	return cluster
}

var _ = Describe("VPC Endpoints", func() {

	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		p = mockprovider.NewMockProvider()
	})

	DescribeTable("can set cluster endpoint configuration on VPC from running Cluster",
		func(e endpointAccessCase) {
			p.MockEKS()
			cluster = NewFakeClusterWithEndpoints(e.Private, e.Public, e.ClusterName)
			mockResultFn := func(_ *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
				return &eks.DescribeClusterOutput{Cluster: cluster}
			}

			p.MockEKS().On("DescribeCluster", MatchedBy(func(input *eks.DescribeClusterInput) bool {
				return input != nil
			})).Return(mockResultFn, e.Error)

			err := UseEndpointAccessFromCluster(p, cfg)
			if err != nil {
				Expect(err.Error() == eks.ErrCodeResourceNotFoundException)
			}

			Expect(err == nil)
			Expect(cfg.VPC.ClusterEndpoints.PrivateAccess == &e.Private)
			Expect(cfg.VPC.ClusterEndpoints.PublicAccess == &e.Public)
		},
		Entry("Private=false, Public=true", endpointAccessCase{
			ClusterName: "false-true-cluster",
			Private:     false,
			Public:      true,
			DCOutput:    &eks.DescribeClusterOutput{Cluster: cluster},
			Error:       nil,
		}),
		Entry("Private=true, Public=false", endpointAccessCase{
			ClusterName: "true-false-cluster",
			Private:     true,
			Public:      false,
			DCOutput:    &eks.DescribeClusterOutput{Cluster: cluster},
			Error:       nil,
		}),
		Entry("Private=true, Public=true", endpointAccessCase{
			ClusterName: "true-true-cluster",
			Private:     true,
			Public:      true,
			DCOutput:    &eks.DescribeClusterOutput{Cluster: cluster},
			Error:       nil,
		}),
		Entry("Private=false, Public=false", endpointAccessCase{
			ClusterName: "notFouncCluster",
			Private:     false,
			Public:      false,
			DCOutput:    nil,
			Error:       errors.New(eks.ErrCodeResourceNotFoundException),
		}),
	)
})
