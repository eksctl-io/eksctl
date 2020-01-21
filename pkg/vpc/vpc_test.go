package vpc

import (
	"errors"
	"fmt"
	"net"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
	"github.com/weaveworks/eksctl/pkg/utils/strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type setSubnetsCase struct {
	cfg   *api.ClusterConfig
	error error
}

type importVPCCase struct {
	cfg               *api.ClusterConfig
	id                string
	describeVPCOutput *ec2.DescribeVpcsOutput
	describeVPCError  error
	error             error
}

type useFromClusterCase struct {
	cfg   *api.ClusterConfig
	stack *cfn.Stack
	error error
}

type endpointAccessCase struct {
	cfg                   *api.ClusterConfig
	clusterName           string
	private               bool
	public                bool
	describeClusterOutput *eks.DescribeClusterOutput
	error                 error
}

type cleanupSubnetsCase struct {
	cfg  *api.ClusterConfig
	want *api.ClusterConfig
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
			if err := SetSubnets(subnetsCase.cfg); err != nil {
				Expect(err).To(Equal(subnetsCase.error))
			} else {
				// make sure that expected error is nil as well
				Expect(subnetsCase.error).Should(BeNil())
			}
		},
		Entry("VPC with valid details", setSubnetsCase{
			cfg: api.NewClusterConfig(),
		}),
		Entry("VPC with nil CIDR", setSubnetsCase{
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
		}),
		Entry("VPC with invalid CIDR prefix", setSubnetsCase{
			cfg: &api.ClusterConfig{
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
			error: fmt.Errorf("VPC CIDR prefix must be between /16 and /24"),
		}),
		Entry("VPC with invalid CIDR", setSubnetsCase{
			cfg: &api.ClusterConfig{
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
			error: fmt.Errorf("Unexpected IP address type: <nil>"),
		}),

		Entry("VPC with invalid number of subnets", setSubnetsCase{
			cfg: &api.ClusterConfig{
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
			error: fmt.Errorf("insufficient number of subnets (have 8, but need 10) for 5 availability zones"),
		}),
		Entry("VPC with multiple AZs", setSubnetsCase{
			cfg: &api.ClusterConfig{
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

			if err := UseFromCluster(p, clusterCase.stack, clusterCase.cfg); err != nil {
				Expect(err.Error()).To(ContainSubstring(clusterCase.error.Error()))
			} else {
				// make sure that expected error is nil as well
				Expect(clusterCase.error).Should(BeNil())
			}
		},
		Entry("No output", useFromClusterCase{
			cfg:   api.NewClusterConfig(),
			stack: &cfn.Stack{},
			error: fmt.Errorf("no output"),
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
				return vpcCase.describeVPCOutput
			}

			p.MockEC2().On("DescribeVpcs", MatchedBy(func(input *ec2.DescribeVpcsInput) bool {
				return input != nil
			})).Return(mockResultFn, vpcCase.describeVPCError)

			if err := importVPC(p, vpcCase.cfg, vpcCase.id); err != nil {
				Expect(err.Error()).To(Equal(vpcCase.error.Error()))
			} else {
				Expect(vpcCase.cfg.VPC.ID == vpcCase.id)
			}
		},
		Entry("VPC with valid details", importVPCCase{
			cfg: api.NewClusterConfig(),
			id:  "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("192.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            nil,
		}),
		Entry("VPC with invalid id", importVPCCase{
			cfg:              api.NewClusterConfig(),
			id:               "invalidID",
			describeVPCError: errors.New("unable to describe vpc"),
			error:            errors.New("unable to describe vpc"),
		}),
		Entry("VPC with id mismatch", importVPCCase{
			cfg: &api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "anotherID",
					},
				},
			},
			id: "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						VpcId: strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            fmt.Errorf("VPC ID %q is not the same as %q", "anotherID", "validID"),
		}),
		Entry("VPC with CIDR mismatch", importVPCCase{
			cfg: api.NewClusterConfig(),
			id:  "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("10.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            fmt.Errorf("VPC CIDR block %q is not the same as %q", "192.168.0.0/16", "10.168.0.0/16"),
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
			id: "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("10.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            fmt.Errorf("VPC CIDR block %q is not the same as %q", "192.168.0.0/16", "10.168.0.0/16"),
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
			id: "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("*"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            fmt.Errorf("invalid CIDR address: *"),
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
			cluster = newFakeClusterWithEndpoints(e.private, e.public, e.clusterName)
			mockResultFn := func(_ *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
				return &eks.DescribeClusterOutput{Cluster: cluster}
			}

			p.MockEKS().On("DescribeCluster", MatchedBy(func(input *eks.DescribeClusterInput) bool {
				return input != nil
			})).Return(mockResultFn, e.error)

			if err := UseEndpointAccessFromCluster(p, e.cfg); err != nil {
				Expect(err.Error()).To(Equal(eks.ErrCodeResourceNotFoundException))
			} else {
				Expect(e.cfg.VPC.ClusterEndpoints.PrivateAccess).To(Equal(&e.private))
				Expect(e.cfg.VPC.ClusterEndpoints.PublicAccess).To(Equal(&e.public))
			}
		},
		Entry("Private=false, Public=true", endpointAccessCase{
			cfg:                   api.NewClusterConfig(),
			clusterName:           "false-true-cluster",
			private:               false,
			public:                true,
			describeClusterOutput: &eks.DescribeClusterOutput{Cluster: cluster},
			error:                 nil,
		}),
		Entry("Private=true, Public=false", endpointAccessCase{
			cfg:                   api.NewClusterConfig(),
			clusterName:           "true-false-cluster",
			private:               true,
			public:                false,
			describeClusterOutput: &eks.DescribeClusterOutput{Cluster: cluster},
			error:                 nil,
		}),
		Entry("Private=true, Public=true", endpointAccessCase{
			cfg:                   api.NewClusterConfig(),
			clusterName:           "true-true-cluster",
			private:               true,
			public:                true,
			describeClusterOutput: &eks.DescribeClusterOutput{Cluster: cluster},
			error:                 nil,
		}),
		Entry("Private=false, Public=false", endpointAccessCase{
			cfg:                   api.NewClusterConfig(),
			clusterName:           "notFoundCluster",
			private:               false,
			public:                false,
			describeClusterOutput: nil,
			error:                 errors.New(eks.ErrCodeResourceNotFoundException),
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
			clusterName:           "notFoundCluster",
			private:               false,
			public:                false,
			describeClusterOutput: nil,
			error:                 nil,
		}),
	)
})

var _ = Describe("VPC - Clean up subnets", func() {
	cfgWithAllAZ := &api.ClusterConfig{
		VPC: &api.ClusterVPC{
			Subnets: &api.ClusterSubnets{
				Private: map[string]api.Network{
					"az1": {
						ID: "private1",
					},
					"az2": {
						ID: "private2",
					},
					"az3": {
						ID: "private3",
					},
				},
				Public: map[string]api.Network{
					"az1": {
						ID: "public1",
					},
					"az2": {
						ID: "public2",
					},
					"az3": {
						ID: "public3",
					},
				},
			},
		},
		AvailabilityZones: []string{"az1", "az2", "az3"},
	}

	DescribeTable("clean up the subnets details in spec if given AZ is invalid",
		func(e cleanupSubnetsCase) {
			cleanupSubnets(e.cfg)
			Expect(e.cfg).To(Equal(cfgWithAllAZ))
		},

		Entry("All AZs are valid", cleanupSubnetsCase{
			cfg: &api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Subnets: &api.ClusterSubnets{
						Private: map[string]api.Network{
							"az1": {
								ID: "private1",
							},
							"az2": {
								ID: "private2",
							},
							"az3": {
								ID: "private3",
							},
						},
						Public: map[string]api.Network{
							"az1": {
								ID: "public1",
							},
							"az2": {
								ID: "public2",
							},
							"az3": {
								ID: "public3",
							},
						},
					},
				},
				AvailabilityZones: []string{"az1", "az2", "az3"},
			},
		}),

		Entry("Private subnet with invalid AZ", cleanupSubnetsCase{
			cfg: &api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Subnets: &api.ClusterSubnets{
						Private: map[string]api.Network{
							"az1": {
								ID: "private1",
							},
							"az2": {
								ID: "private2",
							},
							"az3": {
								ID: "private3",
							},
							"invalid AZ": {
								ID: "invalid private id",
							},
						},
						Public: map[string]api.Network{
							"az1": {
								ID: "public1",
							},
							"az2": {
								ID: "public2",
							},
							"az3": {
								ID: "public3",
							},
						},
					},
				},
				AvailabilityZones: []string{"az1", "az2", "az3"},
			},
			want: cfgWithAllAZ,
		}),
		Entry("Public subnet with invalid AZ", cleanupSubnetsCase{
			cfg: &api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Subnets: &api.ClusterSubnets{
						Private: map[string]api.Network{
							"az1": {
								ID: "private1",
							},
							"az2": {
								ID: "private2",
							},
							"az3": {
								ID: "private3",
							},
						},
						Public: map[string]api.Network{
							"az1": {
								ID: "public1",
							},
							"az2": {
								ID: "public2",
							},
							"az3": {
								ID: "public3",
							},
							"invalid AZ": {
								ID: "invalid public id",
							},
						},
					},
				},
				AvailabilityZones: []string{"az1", "az2", "az3"},
			},
			want: cfgWithAllAZ,
		}),
	)
})
