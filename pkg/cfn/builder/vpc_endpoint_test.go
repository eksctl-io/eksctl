package builder_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/outposts"
	outpoststypes "github.com/aws/aws-sdk-go-v2/service/outposts/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	_ "embed"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

type vpcResourceSetCase struct {
	clusterConfig  *api.ClusterConfig
	expectedFile   string
	createProvider func() api.ClusterProvider
	err            string
}

var _ = Describe("VPC Endpoint Builder", func() {
	DescribeTable("Adds resources to template", func(vc vpcResourceSetCase) {
		api.SetClusterConfigDefaults(vc.clusterConfig)

		if len(vc.clusterConfig.AvailabilityZones) == 0 {
			switch api.Partitions.ForRegion(vc.clusterConfig.Metadata.Region) {
			case api.PartitionAWS:
				vc.clusterConfig.AvailabilityZones = makeZones("us-west-2", 4)
			case api.PartitionChina:
				vc.clusterConfig.AvailabilityZones = makeZones("cn-north-1", 2)
			case api.PartitionISO:
				vc.clusterConfig.AvailabilityZones = makeZones("us-iso-east-1", 2)
			default:
				panic("not supported in tests")
			}
		}
		if vc.clusterConfig.VPC.ID == "" {
			Expect(vpc.SetSubnets(vc.clusterConfig.VPC, vc.clusterConfig.AvailabilityZones, nil)).To(Succeed())
		}

		var provider api.ClusterProvider
		if vc.createProvider != nil {
			provider = vc.createProvider()
		} else {
			provider = mockprovider.NewMockProvider()
		}

		rs := builder.NewRS()
		template := builder.GetTemplate(rs)
		var vpcResourceSet builder.VPCResourceSet = builder.NewIPv4VPCResourceSet(rs, vc.clusterConfig, provider.EC2(), false)
		if vc.clusterConfig.VPC.ID != "" {
			vpcResourceSet = builder.NewExistingVPCResourceSet(rs, vc.clusterConfig, provider.EC2())
		}
		vpcID, subnetDetails, err := vpcResourceSet.CreateTemplate(context.Background())
		if vc.err != "" {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("subnets must be associated with a non-main route table"))
			return
		}

		Expect(err).NotTo(HaveOccurred())
		if vc.clusterConfig.PrivateCluster.Enabled {
			vpcEndpointResourceSet := builder.NewVPCEndpointResourceSet(provider.EC2(), provider.Region(), rs, vc.clusterConfig, vpcID, subnetDetails.Private, gfnt.NewString("sg-test"))
			Expect(vpcEndpointResourceSet.AddResources(context.Background())).To(Succeed())
			s3Endpoint := template.Resources["VPCEndpointS3"].(*gfnec2.VPCEndpoint)
			routeIdsSlice, ok := s3Endpoint.RouteTableIds.Raw().(gfnt.Slice)
			Expect(ok).To(BeTrue())
			sort.Slice(routeIdsSlice, func(i, j int) bool {
				return routeIdsSlice[i].String() < routeIdsSlice[j].String()
			})
		} else if vc.clusterConfig.VPC.ID != "" {
			Expect(template.Resources).To(BeEmpty())
			return
		}

		template.Outputs = nil
		resourceJSON, err := template.JSON()
		Expect(err).NotTo(HaveOccurred())

		expectedJSON, err := os.ReadFile("testdata/" + vc.expectedFile)
		Expect(err).NotTo(HaveOccurred())
		Expect(resourceJSON).To(MatchJSON(expectedJSON))
	},
		Entry("Standard cluster", vpcResourceSetCase{
			clusterConfig: api.NewClusterConfig(),
			expectedFile:  "vpc_public.json",
		}),
		Entry("Private cluster", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "us-west-2",
				},
				VPC: api.NewClusterVPC(false),
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPCEndpoints(provider, serviceDetailsJSON)
				return provider
			},
			expectedFile: "vpc_private.json",
		}),
		Entry("China region cluster", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "cn-north-1",
				},
				VPC: api.NewClusterVPC(false),
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPCEndpoints(provider, serviceDetailsChinaJSON)
				provider.SetRegion("cn-north-1")
				return provider
			},
			expectedFile: "vpc_private_china.json",
		}),
		Entry("Non-private cluster with a user-supplied VPC", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "us-west-2",
				},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc-custom",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"us-west-2a": {
								ID: "subnet-custom1",
							},
							"us-west-2b": {
								ID: "subnet-custom2",
							},
						}),
						Public: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{}),
					},
				},
				AvailabilityZones: []string{"us-west-2a", "us-west-2b"},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPC(provider)
				return provider
			},
		}),
		Entry("Private cluster with a user-supplied VPC", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "us-west-2",
				},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc-custom",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"us-west-2a": {
								ID: "subnet-custom1",
							},
							"us-west-2b": {
								ID: "subnet-custom2",
							},
						}),
					},
				},
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPC(provider)
				mockDescribeVPCEndpoints(provider, serviceDetailsJSON)
				mockDescribeRouteTables(provider, []string{"subnet-custom1", "subnet-custom2"})
				return provider
			},
			expectedFile: "custom_vpc_private_endpoint.json",
		}),
		Entry("Private cluster with a user-supplied VPC has same route table", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "us-west-2",
				},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc-custom",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"us-west-2a": {
								ID: "subnet-custom1",
							},
							"us-west-2b": {
								ID: "subnet-custom2",
							},
						}),
					},
				},
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPC(provider)
				mockDescribeVPCEndpoints(provider, serviceDetailsJSON)
				mockDescribeRouteTablesSame(provider, []string{"subnet-custom1", "subnet-custom2"})
				return provider
			},
			expectedFile: "custom_vpc_private_endpoint_same_route_table.json",
		}),
		Entry("Private cluster with a user-supplied VPC having subnets with an explicit main route table association", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "us-west-2",
				},
				VPC: &api.ClusterVPC{
					Network: api.Network{
						ID: "vpc-custom",
					},
					Subnets: &api.ClusterSubnets{
						Private: api.AZSubnetMappingFromMap(map[string]api.AZSubnetSpec{
							"us-west-2a": {
								ID: "subnet-custom1",
							},
						}),
					},
				},
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPC(provider)
				output := &ec2.DescribeRouteTablesOutput{
					RouteTables: []ec2types.RouteTable{
						{
							VpcId:        aws.String("vpc-custom"),
							RouteTableId: aws.String("rt-main"),
							Associations: []ec2types.RouteTableAssociation{
								{
									RouteTableId:            aws.String("rt-main"),
									RouteTableAssociationId: aws.String("rtbassoc-custom1"),
									Main:                    aws.Bool(true),
								},
								{
									RouteTableId:            aws.String("rt-main"),
									SubnetId:                aws.String("subnet-custom1"),
									RouteTableAssociationId: aws.String("rtbassoc-custom2"),
									Main:                    aws.Bool(false),
								},
							},
						},
					},
				}
				provider.MockEC2().On("DescribeRouteTables", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
					return len(input.Filters) > 0
				})).Return(output, nil)
				return provider
			},
			err: "subnets must be associated with a non-main route table",
		}),

		Entry("Private cluster on Outposts", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "us-west-2",
				},
				VPC: api.NewClusterVPC(false),
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN:   "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
					ControlPlaneInstanceType: "m5.xlarge",
				},
				AvailabilityZones: []string{"us-west-2a"},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPCEndpoints(provider, serviceDetailsOutpostsJSON)
				mockOutposts(provider, "arn:aws:outposts:us-west-2:1234:outpost/op-1234", "us-west-2a")
				return provider
			},
			expectedFile: "vpc_private_outposts.json",
		}),

		Entry("Private cluster on Outposts in a China region", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "cn-north-1",
				},
				VPC: api.NewClusterVPC(false),
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN:   "arn:aws:outposts:cn-north-1:1234:outpost/op-1234",
					ControlPlaneInstanceType: "m5.xlarge",
				},
				AvailabilityZones: []string{"cn-north-1a"},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPCEndpoints(provider, serviceDetailsOutpostsChinaJSON)
				mockOutposts(provider, "arn:aws:outposts:cn-north-1:1234:outpost/op-1234", "cn-north-1a")
				provider.SetRegion("cn-north-1")
				return provider
			},
			expectedFile: "vpc_private_outposts_china.json",
		}),

		Entry("Private cluster in an ISO region", vpcResourceSetCase{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Region: "us-iso-east-1",
				},
				VPC: api.NewClusterVPC(false),
				PrivateCluster: &api.PrivateCluster{
					Enabled: true,
				},
			},
			createProvider: func() api.ClusterProvider {
				provider := mockprovider.NewMockProvider()
				mockDescribeVPCEndpoints(provider, serviceDetailsISOJSON)
				provider.SetRegion("us-iso-east-1")
				return provider
			},
			expectedFile: "vpc_private_iso.json",
		}),
	)
})

//go:embed testdata/service_details.json
var serviceDetailsJSON []byte

//go:embed testdata/service_details_china.json
var serviceDetailsChinaJSON []byte

//go:embed testdata/service_details_outposts.json
var serviceDetailsOutpostsJSON []byte

//go:embed testdata/service_details_outposts_china.json
var serviceDetailsOutpostsChinaJSON []byte

//go:embed testdata/service_details_iso.json
var serviceDetailsISOJSON []byte

func mockDescribeVPC(provider *mockprovider.MockProvider) {
	provider.MockEC2().On("DescribeVpcs", mock.Anything, &ec2.DescribeVpcsInput{
		VpcIds: []string{"vpc-custom"},
	}).Return(&ec2.DescribeVpcsOutput{
		Vpcs: []ec2types.Vpc{
			{
				VpcId: aws.String("vpc-custom"),
				Ipv6CidrBlockAssociationSet: []ec2types.VpcIpv6CidrBlockAssociation{
					{
						Ipv6CidrBlock: aws.String("foo"),
					},
				},
			},
		},
	}, nil)
}

func mockDescribeVPCEndpoints(provider *mockprovider.MockProvider, serviceDetailsJSON []byte) {
	var output *ec2.DescribeVpcEndpointServicesOutput
	Expect(json.Unmarshal(serviceDetailsJSON, &output)).To(Succeed())

	provider.MockEC2().On("DescribeVpcEndpointServices", mock.Anything, mock.MatchedBy(func(e *ec2.DescribeVpcEndpointServicesInput) bool {
		return reflect.DeepEqual(e.ServiceNames, output.ServiceNames)
	})).Return(output, nil)
}

func mockOutposts(provider *mockprovider.MockProvider, outpostARN, az string) {
	provider.MockOutposts().On("GetOutpost", mock.Anything, &outposts.GetOutpostInput{
		OutpostId: aws.String(outpostARN),
	}).Return(&outposts.GetOutpostOutput{
		Outpost: &outpoststypes.Outpost{
			AvailabilityZone: aws.String(az),
		},
	}, nil)
	provider.MockOutposts().On("GetOutpostInstanceTypes", mock.Anything, &outposts.GetOutpostInstanceTypesInput{
		OutpostId: aws.String(outpostARN),
	}).Return(&outposts.GetOutpostInstanceTypesOutput{
		InstanceTypes: []outpoststypes.InstanceTypeItem{
			{
				InstanceType: aws.String("m5.xlarge"),
			},
		},
	}, nil)
}

func mockDescribeRouteTables(provider *mockprovider.MockProvider, subnetIDs []string) {
	output := &ec2.DescribeRouteTablesOutput{
		RouteTables: make([]ec2types.RouteTable, len(subnetIDs)),
	}

	for i, subnetID := range subnetIDs {
		rtID := aws.String(fmt.Sprintf("rtb-custom-%d", i+1))
		output.RouteTables[i] = ec2types.RouteTable{
			VpcId:        aws.String("vpc-custom"),
			RouteTableId: rtID,
			Associations: []ec2types.RouteTableAssociation{
				{
					RouteTableId:            rtID,
					SubnetId:                aws.String(subnetID),
					RouteTableAssociationId: aws.String("rtbassoc-custom"),
					Main:                    aws.Bool(false),
				},
			},
		}
	}

	provider.MockEC2().On("DescribeRouteTables", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
		return len(input.Filters) > 0
	})).Return(output, nil)
}

func mockDescribeRouteTablesSame(provider *mockprovider.MockProvider, subnetIDs []string) {
	output := &ec2.DescribeRouteTablesOutput{
		RouteTables: make([]ec2types.RouteTable, len(subnetIDs)),
	}

	for i, subnetID := range subnetIDs {
		rtID := aws.String("rtb-custom-1")
		output.RouteTables[i] = ec2types.RouteTable{
			VpcId:        aws.String("vpc-custom"),
			RouteTableId: rtID,
			Associations: []ec2types.RouteTableAssociation{
				{
					RouteTableId:            rtID,
					SubnetId:                aws.String(subnetID),
					RouteTableAssociationId: aws.String("rtbassoc-custom"),
					Main:                    aws.Bool(false),
				},
			},
		}
	}

	provider.MockEC2().On("DescribeRouteTables", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
		return len(input.Filters) > 0
	})).Return(output, nil)
}

func makeZones(region string, count int) []string {
	var ret []string
	for i := 0; i < count; i++ {
		ret = append(ret, fmt.Sprintf("%s%c", region, 'a'+i))
	}
	return ret
}
