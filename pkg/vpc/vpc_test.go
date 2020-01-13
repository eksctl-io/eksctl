package vpc

import (
	"errors"
	"fmt"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type setSubnetsCase struct {
	Cfg   *api.ClusterConfig
	Error error
}

type importVPCCase struct {
	cfg               *api.ClusterConfig
	ID                string
	DescribeVPCOutput *ec2.DescribeVpcsOutput
	DescribeVPCError  error
	Error             error
}

type useFromClusterCase struct {
	cfg   *api.ClusterConfig
	Stack *cfn.Stack
	Error error
}

type endpointAccessCase struct {
	cfg         *api.ClusterConfig
	ClusterName string
	Private     bool
	Public      bool
	DCOutput    *eks.DescribeClusterOutput
	Error       error
}

var (
	cluster *eks.Cluster
	p       *mockprovider.MockProvider
)

func newFakeClusterWithEndpoints(private, public bool, name string) *eks.Cluster {
	cluster := NewFakeCluster(name, eks.ClusterStatusActive)
	vpcCfgReq := eks.VpcConfigResponse{}
	vpcCfgReq.SetEndpointPrivateAccess(private).SetEndpointPublicAccess(public)
	cluster.SetResourcesVpcConfig(&vpcCfgReq)
	return cluster
}

var _ = Describe("VPC - Set Subnets", func() {
	DescribeTable("Set subnets",
		func(subnetsCase setSubnetsCase) {
			if err := SetSubnets(subnetsCase.Cfg); err != nil {
				Expect(err).To(Equal(subnetsCase.Error))
			} else {
				// make sure that expected error is nil as well
				Expect(subnetsCase.Error).Should(BeNil())
			}
		},
		Entry("VPC with valid details", setSubnetsCase{
			Cfg: api.NewClusterConfig(),
		}),
		Entry("VPC with nil CIDR", setSubnetsCase{
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
		}),
		Entry("VPC with invalid CIDR prefix", setSubnetsCase{
			Cfg: &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				IAM: &api.ClusterIAM{},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						CIDR: &ipnet.IPNet{
							IPNet: net.IPNet{
								IP:   []byte{192, 168, 0, 0},
								Mask: []byte{255, 255, 255, 255}, // invalid mask
							},
						},
					},
				},
				CloudWatch: &api.ClusterCloudWatch{
					ClusterLogging: &api.ClusterCloudWatchLogging{},
				},
			},
			Error: fmt.Errorf("VPC CIDR prefix must be between /16 and /24"),
		}),
		Entry("VPC with invalid CIDR", setSubnetsCase{
			Cfg: &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				IAM: &api.ClusterIAM{},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						CIDR: &ipnet.IPNet{
							IPNet: net.IPNet{
								IP:   []byte{}, // invalid IP
								Mask: []byte{255, 255, 0, 0},
							},
						},
					},
				},
				CloudWatch: &api.ClusterCloudWatch{
					ClusterLogging: &api.ClusterCloudWatchLogging{},
				},
			},
			Error: fmt.Errorf("Unexpected IP address type: <nil>"),
		}),

		Entry("VPC with invalid number of subnets", setSubnetsCase{
			Cfg: &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				IAM: &api.ClusterIAM{},
				VPC: api.NewClusterVPC(),
				CloudWatch: &api.ClusterCloudWatch{
					ClusterLogging: &api.ClusterCloudWatchLogging{},
				},
				AvailabilityZones: []string{"1", "2", "3", "4", "5"}, // more AZ than required
			},
			Error: fmt.Errorf("insufficient number of subnets (have 8, but need 10) for 5 availability zones"),
		}),
		Entry("VPC with multiple AZs", setSubnetsCase{
			Cfg: &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				IAM: &api.ClusterIAM{},
				VPC: api.NewClusterVPC(),
				CloudWatch: &api.ClusterCloudWatch{
					ClusterLogging: &api.ClusterCloudWatchLogging{},
				},
				AvailabilityZones: []string{"1", "2", "3"},
			},
		}),
	)
})

var _ = Describe("VPC - Use From Cluster", func() {
	BeforeEach(func() {
		p = mockprovider.NewMockProvider()
	})

	DescribeTable("Use from Cluster",
		func(clusterCase useFromClusterCase) {
			cluster = newFakeClusterWithEndpoints(true, true, "dummy cluster")
			mockResultFn := func(_ *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
				return &eks.DescribeClusterOutput{Cluster: cluster}
			}

			p.MockEKS().On("DescribeCluster", MatchedBy(func(input *eks.DescribeClusterInput) bool {
				return input != nil
			})).Return(mockResultFn, nil)

			if err := UseFromCluster(p, clusterCase.Stack, clusterCase.cfg); err != nil {
				Expect(err).To(Equal(clusterCase.Error))
			} else {
				// make sure that expected error is nil as well
				Expect(clusterCase.Error).Should(BeNil())
			}
		},
		Entry("No output", useFromClusterCase{
			cfg:   api.NewClusterConfig(),
			Stack: &cfn.Stack{},
			Error: fmt.Errorf("no output \"VPC\""),
		}),
	)
})

var _ = Describe("VPC - Import VPC", func() {
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

			if err := Import(p, vpcCase.cfg, vpcCase.ID); err != nil {
				Expect(err.Error()).To(Equal(vpcCase.Error.Error()))
			} else {
				Expect(vpcCase.cfg.VPC.ID == vpcCase.ID)
			}
		},
		Entry("VPC with valid details", importVPCCase{
			cfg: api.NewClusterConfig(),
			ID:  "validID",
			DescribeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("192.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            nil,
		}),
		Entry("VPC with invalid id", importVPCCase{
			cfg:              api.NewClusterConfig(),
			ID:               "invalidID",
			DescribeVPCError: errors.New("unable to describe vpc"),
			Error:            errors.New("unable to describe vpc"),
		}),
		Entry("VPC with ID mismatch", importVPCCase{
			cfg: &api.ClusterConfig{
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
						VpcId: strings.Pointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            fmt.Errorf("VPC ID %q is not the same as %q", "anotherID", "validID"),
		}),
		Entry("VPC with CIDR mismatch", importVPCCase{
			cfg: api.NewClusterConfig(),
			ID:  "validID",
			DescribeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("10.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            fmt.Errorf("VPC CIDR block %q is not the same as %q", "192.168.0.0/16", "10.168.0.0/16"),
		}),
		Entry("VPC with nil CIDR", importVPCCase{
			cfg: &api.ClusterConfig{
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
						CidrBlock: strings.Pointer("10.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            fmt.Errorf("VPC CIDR block %q is not the same as %q", "192.168.0.0/16", "10.168.0.0/16"),
		}),
		Entry("VPC with nil CIDR and invalid CIDR", importVPCCase{
			cfg: &api.ClusterConfig{
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
						CidrBlock: strings.Pointer("*"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			DescribeVPCError: nil,
			Error:            fmt.Errorf("invalid CIDR address: *"),
		}),
	)
})

var _ = Describe("VPC - Cluster Endpoints", func() {
	BeforeEach(func() {
		p = mockprovider.NewMockProvider()
	})

	DescribeTable("can set cluster endpoint configuration on VPC from running Cluster",
		func(e endpointAccessCase) {
			p.MockEKS()
			cluster = newFakeClusterWithEndpoints(e.Private, e.Public, e.ClusterName)
			mockResultFn := func(_ *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
				return &eks.DescribeClusterOutput{Cluster: cluster}
			}

			p.MockEKS().On("DescribeCluster", MatchedBy(func(input *eks.DescribeClusterInput) bool {
				return input != nil
			})).Return(mockResultFn, e.Error)

			if err := UseEndpointAccessFromCluster(p, e.cfg); err != nil {
				Expect(err.Error()).To(Equal(eks.ErrCodeResourceNotFoundException))
			} else {
				Expect(e.cfg.VPC.ClusterEndpoints.PrivateAccess).To(Equal(&e.Private))
				Expect(e.cfg.VPC.ClusterEndpoints.PublicAccess).To(Equal(&e.Public))
			}
		},
		Entry("Private=false, Public=true", endpointAccessCase{
			cfg:         api.NewClusterConfig(),
			ClusterName: "false-true-cluster",
			Private:     false,
			Public:      true,
			DCOutput:    &eks.DescribeClusterOutput{Cluster: cluster},
			Error:       nil,
		}),
		Entry("Private=true, Public=false", endpointAccessCase{
			cfg:         api.NewClusterConfig(),
			ClusterName: "true-false-cluster",
			Private:     true,
			Public:      false,
			DCOutput:    &eks.DescribeClusterOutput{Cluster: cluster},
			Error:       nil,
		}),
		Entry("Private=true, Public=true", endpointAccessCase{
			cfg:         api.NewClusterConfig(),
			ClusterName: "true-true-cluster",
			Private:     true,
			Public:      true,
			DCOutput:    &eks.DescribeClusterOutput{Cluster: cluster},
			Error:       nil,
		}),
		Entry("Private=false, Public=false", endpointAccessCase{
			cfg:         api.NewClusterConfig(),
			ClusterName: "notFoundCluster",
			Private:     false,
			Public:      false,
			DCOutput:    nil,
			Error:       errors.New(eks.ErrCodeResourceNotFoundException),
		}),
		Entry("Nil Cluster endpoint from config", endpointAccessCase{
			cfg: &api.ClusterConfig{
				TypeMeta: api.ClusterConfigTypeMeta(),
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				IAM: &api.ClusterIAM{},
				VPC: &api.ClusterVPC{
					ClusterEndpoints: nil,
				},
				CloudWatch: &api.ClusterCloudWatch{
					ClusterLogging: &api.ClusterCloudWatchLogging{},
				},
			},
			ClusterName: "notFoundCluster",
			Private:     false,
			Public:      false,
			DCOutput:    nil,
			Error:       nil,
		}),
	)
})
