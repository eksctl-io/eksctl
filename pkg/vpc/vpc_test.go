package vpc

import (
	"errors"
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/onsi/gomega/types"

	"github.com/weaveworks/eksctl/pkg/eks/mocks"
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
	vpc               *api.ClusterVPC
	availabilityZones []string
	error             error
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
	cfg          *api.ClusterConfig
	stack        *cfn.Stack
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
	cfg      api.ClusterConfig
	expected api.ClusterSubnets
	error    error
}

type cleanupSubnetsCase struct {
	cfg  *api.ClusterConfig
	want *api.ClusterConfig
}

type selectSubnetsCase struct {
	nodegroupAZs     []string
	nodegroupSubnets []string
	subnets          api.AZSubnetMapping
	expectIDs        []string
}

func newFakeClusterWithEndpoints(private, public bool, name string) *eks.Cluster {
	cluster := NewFakeCluster(name, eks.ClusterStatusActive)
	vpcCfgReq := eks.VpcConfigResponse{}
	vpcCfgReq.SetEndpointPrivateAccess(private).SetEndpointPublicAccess(public)
	cluster.SetResourcesVpcConfig(&vpcCfgReq)
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

			subnets, err := SplitInto16(&input)
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

			subnets, err := SplitInto8(&input)
			Expect(err).NotTo(HaveOccurred())
			Expect(subnets).To(HaveLen(8))
			for i, subnet := range subnets {
				Expect(subnet.String()).To(Equal(expected[i]))
			}

		})
	})

	DescribeTable("Set subnets",
		func(subnetsCase setSubnetsCase) {
			err := SetSubnets(subnetsCase.vpc, subnetsCase.availabilityZones)
			if subnetsCase.error != nil {
				Expect(err).To(MatchError(subnetsCase.error.Error()))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("VPC with valid details", setSubnetsCase{
			vpc: api.NewClusterVPC(),
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
			error: fmt.Errorf("Unexpected IP address type: <nil>"),
		}),
		Entry("VPC with valid number of subnets", setSubnetsCase{
			vpc:               api.NewClusterVPC(),
			availabilityZones: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
			error:             nil,
		}),
		Entry("VPC with invalid number of subnets", setSubnetsCase{
			vpc:               api.NewClusterVPC(),
			availabilityZones: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, // more AZ than required
			error:             fmt.Errorf("cannot create more than 16 subnets, 18 requested"),
		}),
		Entry("VPC with multiple AZs", setSubnetsCase{
			vpc:               api.NewClusterVPC(),
			availabilityZones: []string{"1", "2", "3"},
		}),
	)

	DescribeTable("Use from Cluster",
		func(clusterCase useFromClusterCase) {
			p := mockprovider.NewMockProvider()
			cluster := newFakeClusterWithEndpoints(true, true, "dummy cluster")
			mockResultFn := func(_ *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
				return &eks.DescribeClusterOutput{Cluster: cluster}
			}

			p.MockEKS().On("DescribeCluster", MatchedBy(func(input *eks.DescribeClusterInput) bool {
				return input != nil
			})).Return(mockResultFn, nil)

			err := UseFromClusterStack(p, clusterCase.stack, clusterCase.cfg)
			if clusterCase.errorMatcher != nil {
				Expect(err.Error()).To(clusterCase.errorMatcher)
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("No output", useFromClusterCase{
			cfg:          api.NewClusterConfig(),
			stack:        &cfn.Stack{},
			errorMatcher: MatchRegexp(`no output "(?:VPC|SecurityGroup)"`),
		}),
	)

	DescribeTable("importVPC",
		func(vpcCase importVPCCase) {
			p := mockprovider.NewMockProvider()
			p.MockEC2()

			mockResultFn := func(_ *ec2.DescribeVpcsInput) *ec2.DescribeVpcsOutput {
				return vpcCase.describeVPCOutput
			}

			p.MockEC2().On("DescribeVpcs", MatchedBy(func(input *ec2.DescribeVpcsInput) bool {
				return input != nil
			})).Return(mockResultFn, vpcCase.describeVPCError)

			err := importVPC(p.EC2(), vpcCase.cfg, vpcCase.id)
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
				Vpcs: []*ec2.Vpc{
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
				Vpcs: []*ec2.Vpc{
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
				Vpcs: []*ec2.Vpc{
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
		Entry(describeImportVPCCase("VPC with secondary CIDR"), importVPCCase{
			cfg: func() *api.ClusterConfig {
				cfg := api.NewClusterConfig()
				cfg.VPC.CIDR = ipnet.MustParseCIDR("10.1.0.0/16")
				return cfg
			}(),
			id: "validID",
			describeVPCOutput: &ec2.DescribeVpcsOutput{
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("10.0.0.0/16"),
						CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
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
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("10.0.0.0/16"),
						CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
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
			p.MockEKS()
			cluster := newFakeClusterWithEndpoints(e.private, e.public, e.clusterName)
			mockResultFn := func(_ *eks.DescribeClusterInput) *eks.DescribeClusterOutput {
				return &eks.DescribeClusterOutput{Cluster: cluster}
			}

			p.MockEKS().On("DescribeCluster", MatchedBy(func(input *eks.DescribeClusterInput) bool {
				return input != nil
			})).Return(mockResultFn, e.error)

			err := UseEndpointAccessFromCluster(p, e.cfg)
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
			error:       errors.New(eks.ErrCodeResourceNotFoundException),
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
				&ec2.DescribeSubnetsInput{Filters: []*ec2.Filter{{
					Name: strings.Pointer("vpc-id"), Values: aws.StringSlice([]string{"vpc1"}),
				}, {
					Name: strings.Pointer("cidr-block"), Values: aws.StringSlice([]string{"192.168.64.0/18"}),
				}}},
			).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						AvailabilityZone: strings.Pointer("az2"),
						CidrBlock:        strings.Pointer("192.168.64.0/18"),
						SubnetId:         strings.Pointer("private2"),
						VpcId:            strings.Pointer("vpc1"),
					},
				},
			}, nil)
			p.MockEC2().On("DescribeSubnets",
				&ec2.DescribeSubnetsInput{Filters: []*ec2.Filter{{
					Name: strings.Pointer("vpc-id"), Values: aws.StringSlice([]string{"vpc1"}),
				}, {
					Name: strings.Pointer("availability-zone"), Values: aws.StringSlice([]string{"az3"}),
				}}},
			).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						AvailabilityZone: strings.Pointer("az3"),
						CidrBlock:        strings.Pointer("192.168.128.0/18"),
						SubnetId:         strings.Pointer("private3"),
						VpcId:            strings.Pointer("vpc1"),
					},
				},
			}, nil)
			p.MockEC2().On("DescribeSubnets",
				&ec2.DescribeSubnetsInput{SubnetIds: aws.StringSlice([]string{"private1"})},
			).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						AvailabilityZone: strings.Pointer("az1"),
						CidrBlock:        strings.Pointer("192.168.0.0/20"),
						SubnetId:         strings.Pointer("private1"),
						VpcId:            strings.Pointer("vpc1"),
					},
				},
			}, nil)
			p.MockEC2().On("DescribeSubnets",
				&ec2.DescribeSubnetsInput{SubnetIds: aws.StringSlice([]string{"public1"})},
			).Return(&ec2.DescribeSubnetsOutput{
				Subnets: []*ec2.Subnet{
					{
						AvailabilityZone: strings.Pointer("az1"),
						CidrBlock:        strings.Pointer("192.168.1.0/20"),
						SubnetId:         strings.Pointer("public1"),
						VpcId:            strings.Pointer("vpc1"),
					},
				},
			}, nil)

			p.MockEC2().On("DescribeVpcs",
				&ec2.DescribeVpcsInput{VpcIds: aws.StringSlice([]string{"vpc1"})},
			).Return(&ec2.DescribeVpcsOutput{
				NextToken: nil,
				Vpcs: []*ec2.Vpc{
					{
						CidrBlock: strings.Pointer("192.168.0.0/16"),
						VpcId:     strings.Pointer("vpc1"),
					},
				},
			}, nil)

			err := ImportSubnetsFromSpec(p, &e.cfg)
			if e.error != nil {
				Expect(err).To(MatchError(e.error.Error()))
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(*e.cfg.VPC.Subnets).To(Equal(e.expected))
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
			expected: api.ClusterSubnets{
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
			error: nil,
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
			expected: api.ClusterSubnets{
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
			error: nil,
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
			expected: api.ClusterSubnets{
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
			error: nil,
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
			expected: api.ClusterSubnets{
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
			error: nil,
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
			expected: api.ClusterSubnets{
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
			error: nil,
		}),
	)

	DescribeTable("select subnets",
		func(e selectSubnetsCase) {
			ids, err := SelectNodeGroupSubnets(e.nodegroupAZs, e.nodegroupSubnets, e.subnets, nil, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(ConsistOf(e.expectIDs))
		},
		Entry("one subnet", selectSubnetsCase{
			nodegroupSubnets: []string{"a"},
			subnets: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				"a": {
					ID: "id-1",
					AZ: "us-east-1a",
				},
			}),
			expectIDs: []string{"id-1"},
		}),
		Entry("one subnet by id", selectSubnetsCase{
			nodegroupSubnets: []string{"id-1"},
			subnets: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
				"a": {
					ID: "id-1",
					AZ: "us-east-1a",
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

	Context("the user provides an optional subnet id", func() {
		var (
			subnetID string
			mockEC2  *mocks.EC2API
			vpcID    string
			az       string
			azMap    map[string]api.AZSubnetSpec
		)
		BeforeEach(func() {
			subnetID = "user-defined-id"
			vpcID = "vpc-id"
			mockEC2 = &mocks.EC2API{}
			az = "us-east-1a"
			azMap = map[string]api.AZSubnetSpec{
				"a": {
					ID: "id-1",
					AZ: az,
				},
				"b": {
					ID: "id-2",
					AZ: az,
				},
			}
		})
		When("the provided subnet exists", func() {
			It("gets information about the subnet and returns it if it exists", func() {
				mockEC2.On("DescribeSubnets", &ec2.DescribeSubnetsInput{
					SubnetIds: aws.StringSlice([]string{subnetID}),
				}).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []*ec2.Subnet{
						{
							SubnetId: &subnetID,
							VpcId:    &vpcID,
						},
					},
				}, nil)
				ids, err := SelectNodeGroupSubnets([]string{az}, []string{subnetID}, api.AZSubnetMappingFromMap(azMap), mockEC2, vpcID)
				Expect(err).NotTo(HaveOccurred())
				Expect(ids).To(ConsistOf("id-1", "id-2", subnetID))
			})
		})

		When("the provided subnet doesn't exist", func() {
			It("returns a proper error", func() {
				mockEC2.On("DescribeSubnets", &ec2.DescribeSubnetsInput{
					SubnetIds: aws.StringSlice([]string{subnetID}),
				}).Return(nil, errors.New("nope"))
				_, err := SelectNodeGroupSubnets([]string{az}, []string{subnetID}, api.AZSubnetMappingFromMap(azMap), mockEC2, vpcID)
				Expect(err).To(MatchError(ContainSubstring("nope")))
			})
		})

		When("the provided subnet is not part of the cluster's VPC", func() {
			It("returns a proper error", func() {
				mockEC2.On("DescribeSubnets", &ec2.DescribeSubnetsInput{
					SubnetIds: aws.StringSlice([]string{subnetID}),
				}).Return(&ec2.DescribeSubnetsOutput{
					Subnets: []*ec2.Subnet{
						{
							SubnetId: &subnetID,
							VpcId:    aws.String("different-vpc-id"),
						},
					},
				}, nil)
				_, err := SelectNodeGroupSubnets([]string{az}, []string{subnetID}, api.AZSubnetMappingFromMap(azMap), mockEC2, vpcID)
				Expect(err).To(MatchError(ContainSubstring("subnet with id \"user-defined-id\" is not in the attached vpc with id \"vpc-id\"")))
			})
		})
	})
})
