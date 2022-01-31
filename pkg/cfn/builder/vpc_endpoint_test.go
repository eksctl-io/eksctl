package builder

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/vpc"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
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
			switch api.Partition(vc.clusterConfig.Metadata.Region) {
			case api.PartitionAWS:
				vc.clusterConfig.AvailabilityZones = []string{"us-west-2a", "us-west-2b", "us-west-2c", "us-west-2d"}
			case api.PartitionChina:
				vc.clusterConfig.AvailabilityZones = []string{"cn-north-1a", "cn-north-1b"}
			default:
				panic("not supported in tests")
			}
		}
		if vc.clusterConfig.VPC.ID == "" {
			Expect(vpc.SetSubnets(vc.clusterConfig.VPC, vc.clusterConfig.AvailabilityZones)).To(Succeed())
		}

		var provider api.ClusterProvider
		if vc.createProvider != nil {
			provider = vc.createProvider()
		} else {
			provider = mockprovider.NewMockProvider()
		}

		rs := newResourceSet()
		var vpcResourceSet VPCResourceSet = NewIPv4VPCResourceSet(rs, vc.clusterConfig, provider.EC2())
		if vc.clusterConfig.VPC.ID != "" {
			vpcResourceSet = NewExistingVPCResourceSet(rs, vc.clusterConfig, provider.EC2())
		}
		vpcID, subnetDetails, err := vpcResourceSet.CreateTemplate()
		if vc.err != "" {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("subnets must be associated with a non-main route table"))
			return
		}

		Expect(err).NotTo(HaveOccurred())
		if vc.clusterConfig.PrivateCluster.Enabled {
			vpcEndpointResourceSet := NewVPCEndpointResourceSet(provider.EC2(), provider.Region(), rs, vc.clusterConfig, vpcID, subnetDetails.Private, gfnt.NewString("sg-test"))
			Expect(vpcEndpointResourceSet.AddResources()).To(Succeed())
			s3Endpoint := rs.template.Resources["VPCEndpointS3"].(*gfnec2.VPCEndpoint)
			routeIdsSlice, ok := s3Endpoint.RouteTableIds.Raw().(gfnt.Slice)
			Expect(ok).To(BeTrue())
			sort.Slice(routeIdsSlice, func(i, j int) bool {
				return routeIdsSlice[i].String() < routeIdsSlice[j].String()
			})
		} else if vc.clusterConfig.VPC.ID != "" {
			Expect(rs.template.Resources).To(BeEmpty())
			return
		}

		rs.template.Outputs = nil
		resourceJSON, err := rs.template.JSON()
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
				mockDescribeVPCEndpoints(provider, false)
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
				mockDescribeVPCEndpoints(provider, true)
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
				mockDescribeVPCEndpoints(provider, false)
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
				mockDescribeVPCEndpoints(provider, false)
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
					RouteTables: []*ec2.RouteTable{
						{
							VpcId:        aws.String("vpc-custom"),
							RouteTableId: aws.String("rt-main"),
							Associations: []*ec2.RouteTableAssociation{
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
				provider.MockEC2().On("DescribeRouteTables", mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
					return len(input.Filters) > 0
				})).Return(output, nil)
				return provider
			},
			err: "subnets must be associated with a non-main route table",
		}),
	)
})

var serviceDetailsJSON = `
{
  "ServiceNames": [
    "com.amazonaws.us-west-2.ec2",
    "com.amazonaws.us-west-2.ecr.api",
    "com.amazonaws.us-west-2.ecr.dkr",
    "com.amazonaws.us-west-2.s3",
    "com.amazonaws.us-west-2.sts"
  ],
  "ServiceDetails": [
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "PrivateDnsName": "ec2.us-west-2.amazonaws.com",
      "ServiceName": "com.amazonaws.us-west-2.ec2",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-0ee6723c76642b3d8",
      "Owner": "amazon",
      "AvailabilityZones": [
        "us-west-2a",
        "us-west-2b",
        "us-west-2c"
      ],
      "AcceptanceRequired": false,
      "BaseEndpointDnsNames": [
        "ec2.us-west-2.vpce.amazonaws.com"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "PrivateDnsName": "api.ecr.us-west-2.amazonaws.com",
      "ServiceName": "com.amazonaws.us-west-2.ecr.api",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-07d1f428f072fd172",
      "Owner": "amazon",
      "AvailabilityZones": [
        "us-west-2a",
        "us-west-2b",
        "us-west-2c",
        "us-west-2d"
      ],
      "AcceptanceRequired": false,
      "BaseEndpointDnsNames": [
        "api.ecr.us-west-2.vpce.amazonaws.com"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "PrivateDnsName": "*.dkr.ecr.us-west-2.amazonaws.com",
      "ServiceName": "com.amazonaws.us-west-2.ecr.dkr",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-09d74a28015a69002",
      "Owner": "amazon",
      "AvailabilityZones": [
        "us-west-2a",
        "us-west-2b",
        "us-west-2c",
        "us-west-2d"
      ],
      "AcceptanceRequired": false,
      "BaseEndpointDnsNames": [
        "dkr.ecr.us-west-2.vpce.amazonaws.com"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Gateway"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "AcceptanceRequired": false,
      "ServiceName": "com.amazonaws.us-west-2.s3",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-0001be97e1865c74e",
      "Owner": "amazon",
      "AvailabilityZones": [
        "us-west-2a",
        "us-west-2b",
        "us-west-2c",
        "us-west-2d"
      ],
      "BaseEndpointDnsNames": [
        "s3.us-west-2.amazonaws.com"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "AcceptanceRequired": false,
      "ServiceName": "com.amazonaws.us-west-2.s3",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-0b5d83f29260cde0d",
      "Owner": "amazon",
      "AvailabilityZones": [
        "us-west-2a",
        "us-west-2b",
        "us-west-2c",
        "us-west-2d"
      ],
      "BaseEndpointDnsNames": [
        "s3.us-west-2.amazonaws.com"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "PrivateDnsName": "sts.us-west-2.amazonaws.com",
      "ServiceName": "com.amazonaws.us-west-2.sts",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-06681ce20e9a3e8c4",
      "Owner": "amazon",
      "AvailabilityZones": [
        "us-west-2a",
        "us-west-2b",
        "us-west-2c"
      ],
      "AcceptanceRequired": false,
      "BaseEndpointDnsNames": [
        "sts.us-west-2.vpce.amazonaws.com"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Gateway"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "AcceptanceRequired": false,
      "ServiceName": "com.amazonaws.us-west-2.ec2",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-non-existing-endpoint-type",
      "Owner": "amazon",
      "AvailabilityZones": [
        "us-west-2a",
        "us-west-2b",
        "us-west-2c",
        "us-west-2d"
      ],
      "BaseEndpointDnsNames": [
        "ec2.us-west-2.amazonaws.com"
      ]
    }
  ]
}
`
var serviceDetailsJSONChina = `
{
  "ServiceNames": [
    "cn.com.amazonaws.cn-north-1.ec2",
    "cn.com.amazonaws.cn-north-1.ecr.api",
    "cn.com.amazonaws.cn-north-1.ecr.dkr",
    "com.amazonaws.cn-north-1.s3",
    "cn.com.amazonaws.cn-north-1.sts"
  ],
  "ServiceDetails": [
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "PrivateDnsName": "ec2.cn-north-1.amazonaws.com.cn",
      "ServiceName": "cn.com.amazonaws.cn-north-1.ec2",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-0ee6723c76642b3d8",
      "Owner": "amazon",
      "AvailabilityZones": [
        "cn-north-1a",
        "cn-north-1b"
      ],
      "AcceptanceRequired": false,
      "BaseEndpointDnsNames": [
        "ec2.cn-north-1.vpce.amazonaws.com.cn"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "PrivateDnsName": "api.ecr.cn-north-1.amazonaws.com.cn",
      "ServiceName": "cn.com.amazonaws.cn-north-1.ecr.api",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-07d1f428f072fd172",
      "Owner": "amazon",
      "AvailabilityZones": [
        "cn-north-1a",
        "cn-north-1b"
      ],
      "AcceptanceRequired": false,
      "BaseEndpointDnsNames": [
        "api.ecr.cn-north-1.vpce.amazonaws.com.cn"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "PrivateDnsName": "*.dkr.ecr.cn-north-1.amazonaws.com.cn",
      "ServiceName": "cn.com.amazonaws.cn-north-1.ecr.dkr",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-09d74a28015a69002",
      "Owner": "amazon",
      "AvailabilityZones": [
        "cn-north-1a",
        "cn-north-1b"
      ],
      "AcceptanceRequired": false,
      "BaseEndpointDnsNames": [
        "dkr.ecr.cn-north-1.vpce.amazonaws.com.cn"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Gateway"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "AcceptanceRequired": false,
      "ServiceName": "com.amazonaws.cn-north-1.s3",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-0001be97e1865c74e",
      "Owner": "amazon",
      "AvailabilityZones": [
        "cn-north-1a",
        "cn-north-1b"
      ],
      "BaseEndpointDnsNames": [
        "s3.cn-north-1.amazonaws.com"
      ]
    },
    {
      "ServiceType": [
        {
          "ServiceType": "Interface"
        }
      ],
      "Tags": [],
      "ManagesVpcEndpoints": false,
      "PrivateDnsName": "sts.cn-north-1.amazonaws.com.cn",
      "ServiceName": "cn.com.amazonaws.cn-north-1.sts",
      "VpcEndpointPolicySupported": true,
      "ServiceId": "vpce-svc-06681ce20e9a3e8c4",
      "Owner": "amazon",
      "AvailabilityZones": [
        "cn-north-1a",
        "cn-north-1b"
      ],
      "AcceptanceRequired": false,
      "BaseEndpointDnsNames": [
        "sts.cn-north-1.vpce.amazonaws.com.cn"
      ]
    }
  ]
}
`

func mockDescribeVPC(provider *mockprovider.MockProvider) {
	provider.MockEC2().On("DescribeVpcs", &awsec2.DescribeVpcsInput{
		VpcIds: aws.StringSlice([]string{"vpc-custom"}),
	}).Return(&awsec2.DescribeVpcsOutput{
		Vpcs: []*awsec2.Vpc{
			{
				VpcId: aws.String("vpc-custom"),
				Ipv6CidrBlockAssociationSet: []*awsec2.VpcIpv6CidrBlockAssociation{
					{
						Ipv6CidrBlock: aws.String("foo"),
					},
				},
			},
		},
	}, nil)
}

func mockDescribeVPCEndpoints(provider *mockprovider.MockProvider, china bool) {
	var detailsJSON = serviceDetailsJSON
	if china {
		detailsJSON = serviceDetailsJSONChina
	}

	var output *ec2.DescribeVpcEndpointServicesOutput
	Expect(json.Unmarshal([]byte(detailsJSON), &output)).To(Succeed())

	provider.MockEC2().On("DescribeVpcEndpointServices", mock.MatchedBy(func(e *ec2.DescribeVpcEndpointServicesInput) bool {
		return len(e.ServiceNames) == 5
	})).Return(output, nil)

}

func mockDescribeRouteTables(provider *mockprovider.MockProvider, subnetIDs []string) {
	output := &ec2.DescribeRouteTablesOutput{
		RouteTables: make([]*ec2.RouteTable, len(subnetIDs)),
	}

	for i, subnetID := range subnetIDs {
		rtID := aws.String(fmt.Sprintf("rtb-custom-%d", i+1))
		output.RouteTables[i] = &ec2.RouteTable{
			VpcId:        aws.String("vpc-custom"),
			RouteTableId: rtID,
			Associations: []*ec2.RouteTableAssociation{
				{
					RouteTableId:            rtID,
					SubnetId:                aws.String(subnetID),
					RouteTableAssociationId: aws.String("rtbassoc-custom"),
					Main:                    aws.Bool(false),
				},
			},
		}
	}

	provider.MockEC2().On("DescribeRouteTables", mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
		return len(input.Filters) > 0
	})).Return(output, nil)
}

func mockDescribeRouteTablesSame(provider *mockprovider.MockProvider, subnetIDs []string) {
	output := &ec2.DescribeRouteTablesOutput{
		RouteTables: make([]*ec2.RouteTable, len(subnetIDs)),
	}

	for i, subnetID := range subnetIDs {
		rtID := aws.String("rtb-custom-1")
		output.RouteTables[i] = &ec2.RouteTable{
			VpcId:        aws.String("vpc-custom"),
			RouteTableId: rtID,
			Associations: []*ec2.RouteTableAssociation{
				{
					RouteTableId:            rtID,
					SubnetId:                aws.String(subnetID),
					RouteTableAssociationId: aws.String("rtbassoc-custom"),
					Main:                    aws.Bool(false),
				},
			},
		}
	}

	provider.MockEC2().On("DescribeRouteTables", mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
		return len(input.Filters) > 0
	})).Return(output, nil)
}
