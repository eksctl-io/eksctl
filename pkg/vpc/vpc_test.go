package vpc

import (
	"context"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/pkg/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	. "github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
	. "github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

type setSubnetsCase struct {
	vpc               *api.ClusterVPC
	availabilityZones []string
	localZones        []string

	error error
}

type importVPCCase struct {
	cfg               *api.ClusterConfig
	id                string
	describeVPCOutput *ec2.DescribeVpcsOutput
	describeVPCError  error
	error             error
}

func describeImportVPCCase(desc string) func(importVPCCase) string {
	return func(c importVPCCase) string {
		var expect = "works"
		if c.error != nil {
			expect = "returns error"
		}
		return fmt.Sprintf("%s %s", desc, expect)
	}
}

type useFromClusterCase struct {
	cfg     *api.ClusterConfig
	stack   *cfntypes.Stack
	mockEC2 func(*mocksv2.EC2)

	expectedVPC  *api.ClusterVPC
	errorMatcher types.GomegaMatcher
}

type endpointAccessCase struct {
	cfg         *api.ClusterConfig
	clusterName string
	private     bool
	public      bool
	error       error
}

type importAllSubnetsCase struct {
	cfg                     api.ClusterConfig
	privateSubnetOutpostARN *string

	expectedSubnets api.ClusterSubnets
	expectedErr     string
}

type cleanupSubnetsCase struct {
	cfg  *api.ClusterConfig
	want *api.ClusterConfig
}

type selectSubnetsCase struct {
	nodegroupAZs []string
	subnets      api.AZSubnetMapping
	expectIDs    []string
}

type selectSubnetsByIDCase struct {
	ng                   *api.NodeGroupBase
	publicSubnetMapping  api.AZSubnetMapping
	privateSubnetMapping api.AZSubnetMapping
	outputSubnetIDs      []string
	expectedErr          error
}

func newFakeClusterWithEndpoints(private, public bool, name string) *ekstypes.Cluster {
	cluster := NewFakeCluster(name, ekstypes.ClusterStatusActive)
	vpcCfgReq := &ekstypes.VpcConfigResponse{
		EndpointPrivateAccess: private,
		EndpointPublicAccess:  public,
	}
	cluster.ResourcesVpcConfig = vpcCfgReq
	return cluster
}

var _ = Describe("VPC", func() {
	Describe("SplitInto16", func() {
		It("splits the block into 16", func() {
			expected := []string{
				"192.168.0.0/20",
				"192.168.16.0/20",
				"192.168.32.0/20",
				"192.168.48.0/20",
				"192.168.64.0/20",
				"192.168.80.0/20",
				"192.168.96.0/20",
				"192.168.112.0/20",
				"192.168.128.0/20",
				"192.168.144.0/20",
				"192.168.160.0/20",
				"192.168.176.0/20",
				"192.168.192.0/20",
				"192.168.208.0/20",
				"192.168.224.0/20",
				"192.168.240.0/20",
			}

			//192.168.0.0/16
			input := net.IPNet{
				IP:   []byte{192, 168, 0, 0},
				Mask: []byte{255, 255, 0, 0},
			}

			subnets, err := SplitInto(&input, 16, 20)
			Expect(err).NotTo(HaveOccurred())
			Expect(subnets).To(HaveLen(16))
			for i, subnet := range subnets {
				Expect(subnet.String()).To(Equal(expected[i]))
			}

		})
	})

	Describe("SplitInto8", func() {
		It("splits the block into 8", func() {
			expected := []string{
				"192.168.0.0/19",
				"192.168.32.0/19",
				"192.168.64.0/19",
				"192.168.96.0/19",
				"192.168.128.0/19",
				"192.168.160.0/19",
				"192.168.192.0/19",
				"192.168.224.0/19",
			}

			//192.168.0.0/16
			input := net.IPNet{
				IP:   []byte{192, 168, 0, 0},
				Mask: []byte{255, 255, 0, 0},
			}

			subnets, err := SplitInto(&input, 8, 19)
			Expect(err).NotTo(HaveOccurred())
			Expect(subnets).To(HaveLen(8))
			for i, subnet := range subnets {
				Expect(subnet.String()).To(Equal(expected[i]))
			}

		})
	})

	Describe("SplitInto with an invalid CIDR block size", func() {
		It("should return an error if the block size is invalid", func() {
			blockSizes := []int{29, 15}
			for _, s := range blockSizes {
				_, err := SplitInto(&net.IPNet{
					IP:   []byte{192, 168, 0, 0},
					Mask: []byte{255, 255, 0, 0},
				}, 8, s)
				Expect(err).To(MatchError(ContainSubstring("CIDR block size must be between a /16 netmask and /28 netmask")))
			}

		})
	})

	DescribeTable("Set subnets",
		func(subnetsCase setSubnetsCase) {
			err := SetSubnets(subnetsCase.vpc, subnetsCase.availabilityZones, subnetsCase.localZones)
			if subnetsCase.error != nil {
				Expect(err).To(MatchError(subnetsCase.error.Error()))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("VPC with valid details", setSubnetsCase{
			vpc: api.NewClusterVPC(false),
		}),
		Entry("VPC with nil CIDR", setSubnetsCase{
			vpc: &api.ClusterVPC{
				Network: api.Network{
					CIDR: nil,
				},
			},
		}),
		Entry("VPC with invalid CIDR prefix", setSubnetsCase{
			vpc: &api.ClusterVPC{
				Network: api.Network{
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{192, 168, 0, 0},
							Mask: []byte{255, 255, 255, 255}, // invalid mask
						},
					},
				},
			},
			error: fmt.Errorf("VPC CIDR prefix must be between /16 and /24"),
		}),
		Entry("VPC with invalid CIDR", setSubnetsCase{
			vpc: &api.ClusterVPC{
				Network: api.Network{
					CIDR: &ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   []byte{}, // invalid IP
							Mask: []byte{255, 255, 0, 0},
						},
					},
				},
			},
			error: fmt.Errorf("unexpected IP address type: <nil>"),
		}),
		Entry("VPC with valid number of subnets", setSubnetsCase{
			vpc:               api.NewClusterVPC(false),
			availabilityZones: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		}),
		Entry("VPC with invalid number of subnets", setSubnetsCase{
			vpc:               api.NewClusterVPC(false),
			availabilityZones: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, // more AZ than required
			error:             fmt.Errorf("cannot create more than 16 subnets, 18 requested"),
		}),
		Entry("VPC with multiple AZs", setSubnetsCase{
			vpc:               api.NewClusterVPC(false),
			availabilityZones: []string{"1", "2", "3"},
		}),
		Entry("VPC with AZs and local zones", setSubnetsCase{
			vpc:               api.NewClusterVPC(false),
			availabilityZones: []string{"us-west-2a", "us-west-2b"},
			localZones:        []string{"us-west-2-lax-1a", "us-west-lax-1b"},
		}),
		Entry("VPC with a single AZ (Outposts)", setSubnetsCase{
			vpc:               api.NewClusterVPC(false),
			availabilityZones: []string{"us-west-2a"},
		}),
	)

	type setSubnetsEntry struct {
		availabilityZones []string
		localZones        []string

		expectedSubnets          *api.ClusterSubnets
		expectedLocalZoneSubnets *api.ClusterSubnets
	}

	DescribeTable("SetSubnets CIDR assignment", func(e setSubnetsEntry) {
		vpc := api.NewClusterVPC(false)
		err := SetSubnets(vpc, e.availabilityZones, e.localZones)
		Expect(err).NotTo(HaveOccurred())
		Expect(vpc.Subnets).To(Equal(e.expectedSubnets))
		Expect(vpc.LocalZoneSubnets).To(Equal(e.expectedLocalZoneSubnets))

	},
		Entry("both availabilityZones and localZones are set", setSubnetsEntry{
			availabilityZones: []string{"us-west-2a", "us-west-2b", "us-west-2c"},
			localZones:        []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},

			expectedSubnets: &api.ClusterSubnets{
				Public: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						AZ:        "us-west-2a",
						CIDR:      ipnet.MustParseCIDR("192.168.0.0/20"),
						CIDRIndex: 0,
					},
					"us-west-2b": api.AZSubnetSpec{
						AZ:        "us-west-2b",
						CIDR:      ipnet.MustParseCIDR("192.168.16.0/20"),
						CIDRIndex: 1,
					},
					"us-west-2c": api.AZSubnetSpec{
						AZ:        "us-west-2c",
						CIDR:      ipnet.MustParseCIDR("192.168.32.0/20"),
						CIDRIndex: 2,
					},
				},
				Private: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						AZ:        "us-west-2a",
						CIDR:      ipnet.MustParseCIDR("192.168.80.0/20"),
						CIDRIndex: 5,
					},
					"us-west-2b": api.AZSubnetSpec{
						AZ:        "us-west-2b",
						CIDR:      ipnet.MustParseCIDR("192.168.96.0/20"),
						CIDRIndex: 6,
					},
					"us-west-2c": api.AZSubnetSpec{
						AZ:        "us-west-2c",
						CIDR:      ipnet.MustParseCIDR("192.168.112.0/20"),
						CIDRIndex: 7,
					},
				},
			},

			expectedLocalZoneSubnets: &api.ClusterSubnets{
				Public: api.AZSubnetMapping{
					"us-west-2-lax-1a": api.AZSubnetSpec{
						AZ:        "us-west-2-lax-1a",
						CIDR:      ipnet.MustParseCIDR("192.168.48.0/20"),
						CIDRIndex: 3,
					},
					"us-west-2-lax-1b": api.AZSubnetSpec{
						AZ:        "us-west-2-lax-1b",
						CIDR:      ipnet.MustParseCIDR("192.168.64.0/20"),
						CIDRIndex: 4,
					},
				},
				Private: api.AZSubnetMapping{
					"us-west-2-lax-1a": api.AZSubnetSpec{
						AZ:        "us-west-2-lax-1a",
						CIDR:      ipnet.MustParseCIDR("192.168.128.0/20"),
						CIDRIndex: 8,
					},
					"us-west-2-lax-1b": api.AZSubnetSpec{
						AZ:        "us-west-2-lax-1b",
						CIDR:      ipnet.MustParseCIDR("192.168.144.0/20"),
						CIDRIndex: 9,
					},
				},
			},
		}),

		Entry("only availabilityZones is set", setSubnetsEntry{
			availabilityZones: []string{"us-west-2a", "us-west-2b", "us-west-2c"},

			expectedSubnets: &api.ClusterSubnets{
				Public: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						AZ:        "us-west-2a",
						CIDR:      ipnet.MustParseCIDR("192.168.0.0/19"),
						CIDRIndex: 0,
					},
					"us-west-2b": api.AZSubnetSpec{
						AZ:        "us-west-2b",
						CIDR:      ipnet.MustParseCIDR("192.168.32.0/19"),
						CIDRIndex: 1,
					},
					"us-west-2c": api.AZSubnetSpec{
						AZ:        "us-west-2c",
						CIDR:      ipnet.MustParseCIDR("192.168.64.0/19"),
						CIDRIndex: 2,
					},
				},
				Private: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						AZ:        "us-west-2a",
						CIDR:      ipnet.MustParseCIDR("192.168.96.0/19"),
						CIDRIndex: 3,
					},
					"us-west-2b": api.AZSubnetSpec{
						AZ:        "us-west-2b",
						CIDR:      ipnet.MustParseCIDR("192.168.128.0/19"),
						CIDRIndex: 4,
					},
					"us-west-2c": api.AZSubnetSpec{
						AZ:        "us-west-2c",
						CIDR:      ipnet.MustParseCIDR("192.168.160.0/19"),
						CIDRIndex: 5,
					},
				},
			},

			expectedLocalZoneSubnets: &api.ClusterSubnets{
				Public:  api.NewAZSubnetMapping(),
				Private: api.NewAZSubnetMapping(),
			},
		}),

		Entry("control plane on Outposts", setSubnetsEntry{
			availabilityZones: []string{"us-west-2a"},

			expectedSubnets: &api.ClusterSubnets{
				Public: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						AZ:        "us-west-2a",
						CIDR:      ipnet.MustParseCIDR("192.168.0.0/19"),
						CIDRIndex: 0,
					},
				},
				Private: api.AZSubnetMapping{
					"us-west-2a": api.AZSubnetSpec{
						AZ:        "us-west-2a",
						CIDR:      ipnet.MustParseCIDR("192.168.32.0/19"),
						CIDRIndex: 1,
					},
				},
			},

			expectedLocalZoneSubnets: &api.ClusterSubnets{
				Public:  api.NewAZSubnetMapping(),
				Private: api.NewAZSubnetMapping(),
			},
		}),
	)

	DescribeTable("Use from Cluster",
		func(clusterCase useFromClusterCase) {
			p := mockprovider.NewMockProvider()
			cluster := newFakeClusterWithEndpoints(true, true, "dummy cluster")

			p.MockEKS().On("DescribeCluster", Anything, MatchedBy(func(input *eks.DescribeClusterInput) bool {
				return input != nil
			})).Return(&eks.DescribeClusterOutput{Cluster: cluster}, nil)

			if clusterCase.mockEC2 != nil {
				clusterCase.mockEC2(p.MockEC2())
			}
			err := UseFromClusterStack(context.Background(), p, clusterCase.stack, clusterCase.cfg)
			if clusterCase.errorMatcher != nil {
				Expect(err.Error()).To(clusterCase.errorMatcher)
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(clusterCase.cfg.VPC).To(Equal(clusterCase.expectedVPC))
			}
		},
		Entry("no output", useFromClusterCase{
			cfg:          api.NewClusterConfig(),
			stack:        &cfntypes.Stack{},
			errorMatcher: MatchRegexp(`no output "(?:VPC|SecurityGroup)"`),
		}),

		Entry("outputs for subnets in availability zones", useFromClusterCase{
			cfg: api.NewClusterConfig(),
			stack: &cfntypes.Stack{
				Outputs: []cfntypes.Output{
					{
						OutputKey:   aws.String("VPC"),
						OutputValue: aws.String("vpc-123"),
					},
					{
						OutputKey:   aws.String("SecurityGroup"),
						OutputValue: aws.String("sg-123"),
					},
					{
						OutputKey:   aws.String("SubnetsPublic"),
						OutputValue: aws.String("subnet-1"),
					},
				},
			},
			expectedVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-123",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
				},
				SecurityGroup: "sg-123",
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
						},
					},
					Private: api.NewAZSubnetMapping(),
				},
				LocalZoneSubnets: &api.ClusterSubnets{
					Private: api.NewAZSubnetMapping(),
					Public:  api.NewAZSubnetMapping(),
				},
				ManageSharedNodeSecurityGroupRules: aws.Bool(true),
				AutoAllocateIPv6:                   aws.Bool(false),
				NAT: &api.ClusterNAT{
					Gateway: aws.String("Single"),
				},
				ClusterEndpoints: &api.ClusterEndpoints{
					PublicAccess:  aws.Bool(true),
					PrivateAccess: aws.Bool(true),
				},
			},

			mockEC2: func(ec2Mock *mocksv2.EC2) {
				ec2Mock.On("DescribeSubnets", Anything, Anything).Return(func(_ context.Context, input *ec2.DescribeSubnetsInput, _ ...func(options *ec2.Options)) *ec2.DescribeSubnetsOutput {
					return &ec2.DescribeSubnetsOutput{
						Subnets: []ec2types.Subnet{
							{
								SubnetId:         aws.String(input.SubnetIds[0]),
								AvailabilityZone: aws.String("us-west-2a"),
								VpcId:            aws.String("vpc-123"),
								CidrBlock:        aws.String("192.168.0.0/20"),
							},
						},
					}
				}, nil).On("DescribeVpcs", Anything, Anything).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []ec2types.Vpc{
						{
							VpcId:     aws.String("vpc-123"),
							CidrBlock: aws.String("192.168.0.0/20"),
						},
					},
				}, nil)
			},
		}),

		Entry("outputs for subnets in availability zones and local zones", useFromClusterCase{
			cfg: api.NewClusterConfig(),
			stack: &cfntypes.Stack{
				Outputs: []cfntypes.Output{
					{
						OutputKey:   aws.String("VPC"),
						OutputValue: aws.String("vpc-123"),
					},
					{
						OutputKey:   aws.String("SecurityGroup"),
						OutputValue: aws.String("sg-123"),
					},
					{
						OutputKey:   aws.String("SubnetsPublic"),
						OutputValue: aws.String("subnet-1"),
					},
					{
						OutputKey:   aws.String("SubnetsLocalZonePrivate"),
						OutputValue: aws.String("subnet-lz1"),
					},
				},
			},
			expectedVPC: &api.ClusterVPC{
				Network: api.Network{
					ID:   "vpc-123",
					CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
				},
				SecurityGroup: "sg-123",
				Subnets: &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID:   "subnet-1",
							AZ:   "us-west-2a",
							CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
						},
					},
					Private: api.NewAZSubnetMapping(),
				},
				LocalZoneSubnets: &api.ClusterSubnets{
					Public: api.NewAZSubnetMapping(),
					Private: api.AZSubnetMapping{
						"us-west-2-lax-1a": api.AZSubnetSpec{
							ID:   "subnet-lz1",
							AZ:   "us-west-2-lax-1a",
							CIDR: ipnet.MustParseCIDR("192.168.0.16/20"),
						},
					},
				},
				ManageSharedNodeSecurityGroupRules: aws.Bool(true),
				AutoAllocateIPv6:                   aws.Bool(false),
				NAT: &api.ClusterNAT{
					Gateway: aws.String("Single"),
				},
				ClusterEndpoints: &api.ClusterEndpoints{
					PublicAccess:  aws.Bool(true),
					PrivateAccess: aws.Bool(true),
				},
			},

			mockEC2: func(ec2Mock *mocksv2.EC2) {
				ec2Mock.On("DescribeSubnets", Anything, Anything).Return(func(_ context.Context, input *ec2.DescribeSubnetsInput, _ ...func(options *ec2.Options)) *ec2.DescribeSubnetsOutput {
					subnet := ec2types.Subnet{
						SubnetId: aws.String(input.SubnetIds[0]),
						VpcId:    aws.String("vpc-123"),
					}
					if input.SubnetIds[0] == "subnet-lz1" {
						subnet.AvailabilityZone = aws.String("us-west-2-lax-1a")
						subnet.CidrBlock = aws.String("192.168.0.16/20")
					} else {
						subnet.AvailabilityZone = aws.String("us-west-2a")
						subnet.CidrBlock = aws.String("192.168.0.0/20")
					}
					return &ec2.DescribeSubnetsOutput{
						Subnets: []ec2types.Subnet{subnet},
					}
				}, nil).On("DescribeVpcs", Anything, Anything).Return(&ec2.DescribeVpcsOutput{
					Vpcs: []ec2types.Vpc{
						{
							VpcId:     aws.String("vpc-123"),
							CidrBlock: aws.String("192.168.0.0/20"),
						},
					},
				}, nil)
			},
		}),
	)

	DescribeTable("importVPC",
		func(vpcCase importVPCCase) {
			p := mockprovider.NewMockProvider()

			mockResultFn := func(context.Context, *ec2.DescribeVpcsInput, ...func(*ec2.Options)) *ec2.DescribeVpcsOutput {
				return vpcCase.describeVPCOutput
			}

			p.MockEC2().On("DescribeVpcs", Anything, MatchedBy(func(input *ec2.DescribeVpcsInput) bool {
				return input != nil
			})).Return(mockResultFn, func(context.Context, *ec2.DescribeVpcsInput, ...func(*ec2.Options)) error {
				return vpcCase.describeVPCError
			})

			err := importVPC(context.Background(), p.EC2(), vpcCase.cfg, vpcCase.id)
			if vpcCase.error != nil {
				Expect(err).To(MatchError(vpcCase.error.Error()))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(vpcCase.cfg.VPC.ID == vpcCase.id)
			}
		},
		Entry(describeImportVPCCase("VPC with valid details"), importVPCCase{
			cfg: api.NewClusterConfig(),
			id:  "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						CidrBlock: strings.Pointer("192.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            nil,
		}),
		Entry(describeImportVPCCase("VPC with invalid id"), importVPCCase{
			cfg:              api.NewClusterConfig(),
			id:               "invalidID",
			describeVPCError: errors.New("unable to describe vpc"),
			error:            errors.New("unable to describe vpc"),
		}),
		Entry(describeImportVPCCase("VPC with id mismatch"), importVPCCase{
			cfg: &api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "anotherID",
					},
				},
			},
			id: "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						VpcId: strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            fmt.Errorf("VPC ID %q is not the same as %q", "anotherID", "validID"),
		}),
		Entry(describeImportVPCCase("VPC with CIDR mismatch"), importVPCCase{
			cfg: api.NewClusterConfig(),
			id:  "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						CidrBlock: strings.Pointer("10.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            fmt.Errorf("VPC CIDR block %q not found in VPC", "192.168.0.0/16"),
		}),
		Entry(describeImportVPCCase("VPC with nil CIDR"), importVPCCase{
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
				Vpcs: []ec2types.Vpc{
					{
						CidrBlock: strings.Pointer("10.168.0.0/16"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            nil,
		}),
		Entry(describeImportVPCCase("VPC with nil CIDR and invalid CIDR"), importVPCCase{
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
				Vpcs: []ec2types.Vpc{
					{
						CidrBlock: strings.Pointer("*"),
						VpcId:     strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            fmt.Errorf("invalid CIDR address: *"),
		}),
		Entry(describeImportVPCCase("VPC with secondary CIDR"), importVPCCase{
			cfg: func() *api.ClusterConfig {
				cfg := api.NewClusterConfig()
				cfg.VPC.CIDR = ipnet.MustParseCIDR("10.1.0.0/16")
				return cfg
			}(),
			id: "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						CidrBlock: strings.Pointer("10.0.0.0/16"),
						CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
							{
								CidrBlock: strings.Pointer("10.1.0.0/16"),
							},
						},
						VpcId: strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            nil,
		}),
		Entry(describeImportVPCCase("VPC with mismatching secondary CIDR"), importVPCCase{
			cfg: func() *api.ClusterConfig {
				cfg := api.NewClusterConfig()
				cfg.VPC.CIDR = ipnet.MustParseCIDR("10.2.0.0/16")
				return cfg
			}(),
			id: "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []ec2types.Vpc{
					{
						CidrBlock: strings.Pointer("10.0.0.0/16"),
						CidrBlockAssociationSet: []ec2types.VpcCidrBlockAssociation{
							{
								CidrBlock: strings.Pointer("10.1.0.0/16"),
							},
						},
						VpcId: strings.Pointer("validID"),
					},
				},
			},
			describeVPCError: nil,
			error:            fmt.Errorf(`VPC CIDR block "10.2.0.0/16" not found in VPC`),
		}),
	)

	DescribeTable("can set cluster endpoint configuration on VPC from running Cluster",
		func(e endpointAccessCase) {
			p := mockprovider.NewMockProvider()
			cluster := newFakeClusterWithEndpoints(e.private, e.public, e.clusterName)

			p.MockEKS().On("DescribeCluster", Anything, MatchedBy(func(input *eks.DescribeClusterInput) bool {
				return input != nil
			})).Return(&eks.DescribeClusterOutput{Cluster: cluster}, e.error)

			err := UseEndpointAccessFromCluster(context.Background(), p, e.cfg)
			if e.error != nil {
				Expect(err).To(MatchError(e.error.Error()))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(e.cfg.VPC.ClusterEndpoints.PrivateAccess).To(Equal(&e.private))
				Expect(e.cfg.VPC.ClusterEndpoints.PublicAccess).To(Equal(&e.public))
			}
		},
		Entry("Private=false, Public=true", endpointAccessCase{
			cfg:         api.NewClusterConfig(),
			clusterName: "false-true-cluster",
			private:     false,
			public:      true,
			error:       nil,
		}),
		Entry("Private=true, Public=false", endpointAccessCase{
			cfg:         api.NewClusterConfig(),
			clusterName: "true-false-cluster",
			private:     true,
			public:      false,
			error:       nil,
		}),
		Entry("Private=true, Public=true", endpointAccessCase{
			cfg:         api.NewClusterConfig(),
			clusterName: "true-true-cluster",
			private:     true,
			public:      true,
			error:       nil,
		}),
		Entry("Private=false, Public=false", endpointAccessCase{
			cfg:         api.NewClusterConfig(),
			clusterName: "notFoundCluster",
			private:     false,
			public:      false,
			error:       &ekstypes.ResourceNotFoundException{},
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
			clusterName: "notFoundCluster",
			private:     false,
			public:      false,
			error:       nil,
		}),
	)

	cfgWithAllAZ := &api.ClusterConfig{
		VPC: &api.ClusterVPC{
			Subnets: &api.ClusterSubnets{
				Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"az1": {
						ID: "private1",
					},
					"az2": {
						ID: "private2",
					},
					"az3": {
						ID: "private3",
					},
				}),
				Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"az1": {
						ID: "public1",
					},
					"az2": {
						ID: "public2",
					},
					"az3": {
						ID: "public3",
					},
				}),
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
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "private1",
							},
							"az2": {
								ID: "private2",
							},
							"az3": {
								ID: "private3",
							},
						}),
						Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "public1",
							},
							"az2": {
								ID: "public2",
							},
							"az3": {
								ID: "public3",
							},
						}),
					},
				},
				AvailabilityZones: []string{"az1", "az2", "az3"},
			},
		}),

		Entry("Private subnet with invalid AZ", cleanupSubnetsCase{
			cfg: &api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
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
						}),
						Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "public1",
							},
							"az2": {
								ID: "public2",
							},
							"az3": {
								ID: "public3",
							},
						}),
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
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "private1",
							},
							"az2": {
								ID: "private2",
							},
							"az3": {
								ID: "private3",
							},
						}),
						Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
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
						}),
					},
				},
				AvailabilityZones: []string{"az1", "az2", "az3"},
			},
			want: cfgWithAllAZ,
		}),
	)

	DescribeTable("Can import all subnets",
		func(e importAllSubnetsCase) {
			p := mockprovider.NewMockProvider()

			p.MockEC2().On("DescribeSubnets",
				Anything,
				&ec2.DescribeSubnetsInput{Filters: []ec2types.Filter{{
					Name: strings.Pointer("vpc-id"), Values: []string{"vpc1"},
				}, {
					Name: strings.Pointer("cidr-block"), Values: []string{"192.168.64.0/18"},
				}}},
			).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []ec2types.Subnet{
					{
						AvailabilityZone: strings.Pointer("az2"),
						CidrBlock:        strings.Pointer("192.168.64.0/18"),
						SubnetId:         strings.Pointer("private2"),
						VpcId:            strings.Pointer("vpc1"),
					},
				},
			}, nil)
			p.MockEC2().On("DescribeSubnets",
				Anything,
				&ec2.DescribeSubnetsInput{Filters: []ec2types.Filter{{
					Name: strings.Pointer("vpc-id"), Values: []string{"vpc1"},
				}, {
					Name: strings.Pointer("availability-zone"), Values: []string{"az3"},
				}}},
			).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []ec2types.Subnet{
					{
						AvailabilityZone: strings.Pointer("az3"),
						CidrBlock:        strings.Pointer("192.168.128.0/18"),
						SubnetId:         strings.Pointer("private3"),
						VpcId:            strings.Pointer("vpc1"),
					},
				},
			}, nil)
			p.MockEC2().On("DescribeSubnets",
				Anything,
				&ec2.DescribeSubnetsInput{SubnetIds: []string{"private1"}},
			).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []ec2types.Subnet{
					{
						AvailabilityZone: strings.Pointer("az1"),
						CidrBlock:        strings.Pointer("192.168.0.0/20"),
						SubnetId:         strings.Pointer("private1"),
						VpcId:            strings.Pointer("vpc1"),
						OutpostArn:       e.privateSubnetOutpostARN,
					},
				},
			}, nil)
			p.MockEC2().On("DescribeSubnets",
				Anything,
				&ec2.DescribeSubnetsInput{SubnetIds: []string{"public1"}},
			).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []ec2types.Subnet{
					{
						AvailabilityZone: strings.Pointer("az1"),
						CidrBlock:        strings.Pointer("192.168.1.0/20"),
						SubnetId:         strings.Pointer("public1"),
						VpcId:            strings.Pointer("vpc1"),
					},
				},
			}, nil)

			p.MockEC2().On("DescribeVpcs",
				Anything,
				&ec2.DescribeVpcsInput{VpcIds: []string{"vpc1"}},
			).Return(&ec2.DescribeVpcsOutput{
				NextToken: nil,
				Vpcs: []ec2types.Vpc{
					{
						CidrBlock: strings.Pointer("192.168.0.0/16"),
						VpcId:     strings.Pointer("vpc1"),
					},
				},
			}, nil)

			err := ImportSubnetsFromSpec(context.Background(), p, &e.cfg)
			if e.expectedErr != "" {
				Expect(err).To(MatchError(e.expectedErr))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(*e.cfg.VPC.Subnets).To(Equal(e.expectedSubnets))
			}
		},

		Entry("Subnet are matching with AZs", importAllSubnetsCase{
			cfg: api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc1",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "private1",
							},
						}),
						Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "public1",
							},
						}),
					},
				},
			},
			expectedSubnets: api.ClusterSubnets{
				Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"az1": {
						ID:   "private1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
					},
				}),
				Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"az1": {
						ID:   "public1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.1.0/20"),
					},
				}),
			},
		}),

		Entry("Private subnet is not matching with AZ", importAllSubnetsCase{
			cfg: api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc1",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"invalidAZ": {
								ID: "private1",
							},
						}),
						Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "public1",
							},
						}),
					},
				},
			},
			expectedSubnets: api.ClusterSubnets{
				Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"invalidAZ": {
						ID:   "private1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
					},
				}),
				Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"az1": {
						ID:   "public1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.1.0/20"),
					},
				}),
			},
		}),
		Entry("Public subnet is not matching with AZ", importAllSubnetsCase{
			cfg: api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc1",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "private1",
							},
						}),
						Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"invalidAZ": {
								ID: "public1",
							},
						}),
					},
				},
			},
			expectedSubnets: api.ClusterSubnets{
				Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"az1": {
						ID:   "private1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
					},
				}),
				Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"invalidAZ": {
						ID:   "public1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.1.0/20"),
					},
				}),
			},
		}),
		Entry("Subnets given names and identified by various params", importAllSubnetsCase{
			cfg: api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc1",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"subone": {
								ID: "private1",
							},
							"subtwo": {
								AZ:   "az2",
								CIDR: ipnet.MustParseCIDR("192.168.64.0/18"),
							},
							"subthree": {
								AZ: "az3",
							},
						}),
					},
				},
			},
			expectedSubnets: api.ClusterSubnets{
				Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"subone": {
						ID:   "private1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
					},
					"subtwo": {
						ID:   "private2",
						AZ:   "az2",
						CIDR: ipnet.MustParseCIDR("192.168.64.0/18"),
					},
					"subthree": {
						ID:   "private3",
						AZ:   "az3",
						CIDR: ipnet.MustParseCIDR("192.168.128.0/18"),
					},
				}),
			},
		}),
		Entry("Subnets identified by CIDR", importAllSubnetsCase{
			cfg: api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc1",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "private1",
							},
							"az2": {
								CIDR: ipnet.MustParseCIDR("192.168.64.0/18"),
							},
						}),
						Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"invalidAZ": {
								ID: "public1",
							},
						}),
					},
				},
			},
			expectedSubnets: api.ClusterSubnets{
				Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"az1": {
						ID:   "private1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.0.0/20"),
					},
					"az2": {
						ID:   "private2",
						AZ:   "az2",
						CIDR: ipnet.MustParseCIDR("192.168.64.0/18"),
					},
				}),
				Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"invalidAZ": {
						ID:   "public1",
						AZ:   "az1",
						CIDR: ipnet.MustParseCIDR("192.168.1.0/20"),
					},
				}),
			},
		}),

		Entry("[Outposts] Subnets not on Outposts", importAllSubnetsCase{
			cfg: api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc1",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "private1",
							},
						}),
					},
				},
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				},
			},

			expectedErr: "all subnets must be on the control plane Outpost when specifying pre-existing subnets for a cluster on Outposts; found invalid private subnet(s): private1",
		}),

		Entry("[Outposts] Subnets not on the same Outpost as the control plane", importAllSubnetsCase{
			cfg: api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc1",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "private1",
							},
						}),
					},
				},
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				},
			},
			privateSubnetOutpostARN: aws.String("arn:aws:outposts:us-west-2:1234:outpost/op-5678"),

			expectedErr: "all subnets must be on the control plane Outpost when specifying pre-existing subnets for a cluster on Outposts; found invalid private subnet(s): private1",
		}),

		Entry("[Outposts] Subnets on the same Outpost as the control plane", importAllSubnetsCase{
			cfg: api.ClusterConfig{
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc1",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"az1": {
								ID: "private1",
							},
						}),
					},
				},
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				},
			},
			privateSubnetOutpostARN: aws.String("arn:aws:outposts:us-west-2:1234:outpost/op-1234"),

			expectedSubnets: api.ClusterSubnets{
				Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
					"az1": {
						ID:         "private1",
						AZ:         "az1",
						CIDR:       ipnet.MustParseCIDR("192.168.0.0/20"),
						OutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
					},
				}),
			},
		}),
	)

	DescribeTable("select subnets",
		func(e selectSubnetsCase) {
			ids, err := selectNodeGroupZoneSubnets(e.nodegroupAZs, e.subnets, func(zone string) error { return nil })
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(ConsistOf(e.expectIDs))
		},
		Entry("one subnet", selectSubnetsCase{
			nodegroupAZs: []string{"a"},
			subnets: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				"a": {
					ID: "id-1",
					AZ: "a",
				},
			}),
			expectIDs: []string{"id-1"},
		}),

		Entry("one AZ", selectSubnetsCase{
			nodegroupAZs: []string{"us-east-1a"},
			subnets: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				"a": {
					ID: "id-1",
					AZ: "us-east-1a",
				},
				"b": {
					ID: "id-2",
					AZ: "us-east-1a",
				},
			}),
			expectIDs: []string{"id-1", "id-2"},
		}),
	)

	DescribeTable("select subnets by id",
		func(e selectSubnetsByIDCase) {
			subnetIDs, err := selectNodeGroupSubnetsFromIDs(context.Background(), e.ng, e.publicSubnetMapping, e.privateSubnetMapping,
				&api.ClusterConfig{}, mockprovider.NewMockProvider().EC2(), func(zone string) error { return nil }, func(zone string) error { return nil })

			if e.expectedErr != nil {
				Expect(err.Error()).To(Equal(e.expectedErr.Error()))
			} else {
				Expect(err).To(BeNil())
				Expect(subnetIDs).To(Equal(e.outputSubnetIDs))
			}
		},
		Entry("set private subnet, by name, with privateNetworking disabled", selectSubnetsByIDCase{
			ng: &api.NodeGroupBase{
				Subnets:           []string{"subnet-name"},
				PrivateNetworking: false,
			},
			privateSubnetMapping: api.AZSubnetMapping{"subnet-name": api.AZSubnetSpec{}},
			expectedErr:          fmt.Errorf("subnet subnet-name is specified as private in ClusterConfig, thus must only be used when `privateNetworking` is enabled"),
		}),
		Entry("set private subnet, by ID, with privateNetworking disabled", selectSubnetsByIDCase{
			ng: &api.NodeGroupBase{
				Subnets:           []string{"subnet-id"},
				PrivateNetworking: false,
			},
			privateSubnetMapping: api.AZSubnetMapping{"subnet-name": api.AZSubnetSpec{ID: "subnet-id"}},
			expectedErr:          fmt.Errorf("subnet subnet-id is specified as private in ClusterConfig, thus must only be used when `privateNetworking` is enabled"),
		}),
		Entry("set public subnet, by name, with privateNetworking enabled", selectSubnetsByIDCase{
			ng: &api.NodeGroupBase{
				Subnets:           []string{"subnet-name"},
				PrivateNetworking: true,
			},
			publicSubnetMapping: api.AZSubnetMapping{"subnet-name": api.AZSubnetSpec{ID: "subnet-id"}},
			outputSubnetIDs:     []string{"subnet-id"},
		}),
		Entry("set public subnet, by ID, with privateNetworking enabled", selectSubnetsByIDCase{
			ng: &api.NodeGroupBase{
				Subnets:           []string{"subnet-id"},
				PrivateNetworking: true,
			},
			publicSubnetMapping: api.AZSubnetMapping{"subnet-name": api.AZSubnetSpec{ID: "subnet-id"}},
			outputSubnetIDs:     []string{"subnet-id"},
		}),
	)

})
