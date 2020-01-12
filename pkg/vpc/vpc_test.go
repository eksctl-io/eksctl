package vpc

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"

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

type importVPCCase struct {
	Cfg               *api.ClusterConfig
	ID                string
	DescribeVPCOutput *ec2.DescribeVpcsOutput
	DescribeVPCError  error
	Error             error
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

var _ = Describe("VPC Endpoints - Cluster Endpoints", func() {

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

			if err := UseEndpointAccessFromCluster(p, cfg); err != nil {
				Expect(err.Error()).To(Equal(eks.ErrCodeResourceNotFoundException))
			} else {
				Expect(cfg.VPC.ClusterEndpoints.PrivateAccess).To(Equal(&e.Private))
				Expect(cfg.VPC.ClusterEndpoints.PublicAccess).To(Equal(&e.Public))
			}
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
			ClusterName: "notFoundCluster",
			Private:     false,
			Public:      false,
			DCOutput:    nil,
			Error:       errors.New(eks.ErrCodeResourceNotFoundException),
		}),
	)
})

var _ = Describe("VPC Endpoints - Import VPC", func() {
	BeforeEach(func() {
		p = mockprovider.NewMockProvider()
	})

	DescribeTable("can import VPC",
		func(vpcCase importVPCCase) {
			p.MockEC2()

			mockResultFn := func(_ *ec2.DescribeVpcsInput) *ec2.DescribeVpcsOutput {
				return vpcCase.DescribeVPCOutput
			}

			p.MockEC2().On("DescribeVpcs", MatchedBy(func(input *ec2.DescribeVpcsInput) bool {
				return input != nil
			})).Return(mockResultFn, vpcCase.DescribeVPCError)

			if err := Import(p, vpcCase.Cfg, vpcCase.ID); err != nil {
				Expect(err.Error()).To(Equal(vpcCase.Error.Error()))
			} else {
				Expect(vpcCase.Cfg.VPC.ID == vpcCase.ID)
			}
		},
		Entry("VPC with valid details", importVPCCase{
			Cfg: api.NewClusterConfig(),
			ID:  "validID",
			DescribeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: stringPointer("192.168.0.0/16"),
						VpcId:     stringPointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            nil,
		}),
		Entry("VPC with invalid id", importVPCCase{
			Cfg:              api.NewClusterConfig(),
			ID:               "invalidID",
			DescribeVPCError: errors.New("unable to describe vpc"),
			Error:            errors.New("unable to describe vpc"),
		}),
		Entry("VPC with ID mismatch", importVPCCase{
			Cfg: &api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "anotherID",
					},
				},
			},
			ID: "validID",
			DescribeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						VpcId: stringPointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            fmt.Errorf("VPC ID %q is not the same as %q", "anotherID", "validID"),
		}),
		Entry("VPC with CIDR mismatch", importVPCCase{
			Cfg: api.NewClusterConfig(),
			ID:  "validID",
			DescribeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: stringPointer("10.168.0.0/16"),
						VpcId:     stringPointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            fmt.Errorf("VPC CIDR block %q is not the same as %q", "192.168.0.0/16", "10.168.0.0/16"),
		}),
		Entry("VPC with nil CIDR", importVPCCase{
			Cfg: &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				IAM: &api.ClusterIAM{},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						CIDR: nil,
					},
				},
				CloudWatch: &api.ClusterCloudWatch{
					ClusterLogging: &api.ClusterCloudWatchLogging{},
				},
			},
			ID: "validID",
			DescribeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: stringPointer("10.168.0.0/16"),
						VpcId:     stringPointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            fmt.Errorf("VPC CIDR block %q is not the same as %q", "192.168.0.0/16", "10.168.0.0/16"),
		}),
		Entry("VPC with nil CIDR and invalid CIDR", importVPCCase{
			Cfg: &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				IAM: &api.ClusterIAM{},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						CIDR: nil,
					},
				},
				CloudWatch: &api.ClusterCloudWatch{
					ClusterLogging: &api.ClusterCloudWatchLogging{},
				},
			},
			ID: "validID",
			DescribeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: stringPointer("*"),
						VpcId:     stringPointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            fmt.Errorf("invalid CIDR address: *"),
		}),
	)
})

func stringPointer(s string) *string {
	return &s
}
