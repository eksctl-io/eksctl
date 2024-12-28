package builder_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gfnt "goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

var _ = Describe("AssignSubnets", func() {
	type assignSubnetsEntry struct {
		np                    api.NodePool
		updateEC2Mocks        func(*mocksv2.EC2)
		updateClusterConfig   func(config *api.ClusterConfig)
		localZones            []string
		availabilityZones     []string
		instanceTypes         []ec2types.InstanceType
		customInstanceSupport bool

		expectedErr       string
		expectedSubnetIDs []string
	}

	toSubnetIDs := func(subnetRefs *gfnt.Value) []string {
		subnetsSlice, ok := subnetRefs.Raw().(gfnt.Slice)
		Expect(ok).To(BeTrue(), fmt.Sprintf("expected subnet refs to be of type %T; got %T", gfnt.Slice{}, subnetRefs.Raw()))
		var subnetIDs []string
		for _, subnetID := range subnetsSlice {
			subnetIDs = append(subnetIDs, subnetID.String())
		}
		return subnetIDs
	}

	const vpcID = "vpc-1"

	DescribeTable("assigns subnets to a nodegroup", func(e assignSubnetsEntry) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.VPC.ID = vpcID

		provider := mockprovider.NewMockProvider()
		if !e.customInstanceSupport {
			mockSubnetsAndAZInstanceSupport(clusterConfig, provider, e.availabilityZones, e.localZones, e.instanceTypes)
		}

		if e.updateEC2Mocks != nil {
			e.updateEC2Mocks(provider.MockEC2())
		}
		if e.updateClusterConfig != nil {
			e.updateClusterConfig(clusterConfig)
		}

		subnetRefs, err := builder.AssignSubnets(context.Background(), e.np, clusterConfig, provider.EC2())
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
		subnetIDs := toSubnetIDs(subnetRefs)
		Expect(err).NotTo(HaveOccurred())
		Expect(subnetIDs).To(ConsistOf(e.expectedSubnetIDs))

	},

		Entry("self-managed nodegroup with availability zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			availabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
			expectedSubnetIDs: []string{"subnet-public-us-west-1a", "subnet-public-us-west-1b", "subnet-public-us-west-1c"},
		}),

		Entry("managed nodegroup with availability zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			availabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
			expectedSubnetIDs: []string{"subnet-public-us-west-1a", "subnet-public-us-west-1b", "subnet-public-us-west-1c"},
		}),

		Entry("self-managed nodegroup with local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{},
				LocalZones:    []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			localZones:        []string{"us-west-2-lax-1a", "us-west-2-lax-1b", "us-west-2-lax-1c"},
			expectedSubnetIDs: []string{"subnet-public-us-west-2-lax-1a", "subnet-public-us-west-2-lax-1b"},
		}),

		Entry("self-managed nodegroup with privateNetworking and local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					PrivateNetworking: true,
				},
				LocalZones: []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			localZones:        []string{"us-west-2-lax-1a", "us-west-2-lax-1b", "us-west-2-lax-1d"},
			expectedSubnetIDs: []string{"subnet-private-us-west-2-lax-1a", "subnet-private-us-west-2-lax-1b"},
		}),

		Entry("self-managed nodegroup with local zones and subnet IDs", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets: []string{"subnet-z1", "subnet-z2"},
				},
				LocalZones: []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			localZones:        []string{"us-west-2-lax-1a", "us-west-2-lax-1b", "us-west-2-lax-1d", "us-west-2-lax-1e"},
			expectedSubnetIDs: []string{"subnet-z1", "subnet-z2", "subnet-public-us-west-2-lax-1a", "subnet-public-us-west-2-lax-1b"},
			updateEC2Mocks: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2-lax-1e", vpcID)
			},
		}),

		Entry("managed nodegroup with privateNetworking, availability zones and subnet IDs", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					PrivateNetworking: true,
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
					Subnets:           []string{"subnet-z1", "subnet-z2"},
				},
			},
			availabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c", "us-west-1d"},
			expectedSubnetIDs: []string{"subnet-private-us-west-1a", "subnet-private-us-west-1b", "subnet-private-us-west-1c", "subnet-z1", "subnet-z2"},
			updateEC2Mocks: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-1d", vpcID)
			},
		}),

		Entry("managed nodegroup with availability zones and subnet IDs in local zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:           []string{"subnet-z1", "subnet-z2"},
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			availabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
			localZones:        []string{"us-west-2-lax-1e"},
			expectedErr:       "managed nodegroups cannot be launched in local zones",
			updateEC2Mocks: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2-lax-1e", vpcID)
			},
		}),

		Entry("managed nodegroup without subnets, availability zones and local zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{},
			},
			availabilityZones: []string{"us-west-1a", "us-west-1b"},
			expectedSubnetIDs: []string{"subnet-public-us-west-1a", "subnet-public-us-west-1b"},
		}),

		Entry("private self-managed nodegroup without subnets, availability zones and local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					PrivateNetworking: true,
				},
			},
			availabilityZones: []string{"us-west-1a", "us-west-1b"},
			expectedSubnetIDs: []string{"subnet-private-us-west-1a", "subnet-private-us-west-1b"},
		}),

		Entry("supplied subnet ID exists in a different VPC", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets: []string{"subnet-1"},
				},
			},
			expectedErr: `subnet with ID "subnet-1" is not in the attached VPC with ID "vpc-1"`,
			updateEC2Mocks: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2a", "vpc-2")
			},
		}),

		Entry("No subnets in zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AvailabilityZones: []string{"us-west-2z"},
				},
			},
			availabilityZones: []string{"us-west-2a"},
			expectedErr:       "could not find public subnets for zones",
		}),

		Entry("EFA enabled with multiple subnets selects only one subnet", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:    []string{"subnet-1", "subnet-2", "subnet-3"},
					EFAEnabled: aws.Bool(true),
				},
			},
			availabilityZones: []string{"us-west-2d"},
			updateEC2Mocks: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2d", vpcID)
			},
			expectedSubnetIDs: []string{"subnet-1"},
		}),

		Entry("nodegroup with subnet names", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets: []string{"subnet-1", "subnet-2", "subnet-3"},
				},
			},
			availabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
			updateClusterConfig: func(config *api.ClusterConfig) {
				config.VPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"subnet-1": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"subnet-2": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"subnet-3": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-1a", "subnet-1b", "subnet-1c"},
		}),

		Entry("managed nodegroup without AZs, local zones or subnets, but not all AZs support the required instance type", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					InstanceType: "g4dn.xlarge",
				},
			},
			updateClusterConfig: func(config *api.ClusterConfig) {
				config.AvailabilityZones = []string{"us-west-2a", "us-west-2d"}
				config.VPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID: "subnet-1",
							AZ: "us-west-2a",
						},
						"us-west-2d": api.AZSubnetSpec{
							ID: "subnet-2",
							AZ: "us-west-2d",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			updateEC2Mocks: func(e *mocksv2.EC2) {
				e.On("DescribeInstanceTypeOfferings", mock.Anything, mock.Anything, mock.Anything).
					Return(&ec2.DescribeInstanceTypeOfferingsOutput{
						InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
							{
								InstanceType: ec2types.InstanceTypeG4dnXlarge,
								Location:     aws.String("us-west-2a"),
								LocationType: ec2types.LocationTypeAvailabilityZone,
							},
							{
								InstanceType: api.DefaultNodeType,
								Location:     aws.String("us-west-2d"),
								LocationType: ec2types.LocationTypeAvailabilityZone,
							},
						},
					}, nil)
				e.On("DescribeAvailabilityZones", mock.Anything, mock.Anything).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []ec2types.AvailabilityZone{
							{
								ZoneType: aws.String("availability-zone"),
								ZoneName: aws.String("us-west-2a"),
							},
							{
								ZoneType: aws.String("availability-zone"),
								ZoneName: aws.String("us-west-2d"),
							},
						},
					}, nil)
			},
			customInstanceSupport: true,
			expectedSubnetIDs:     []string{"subnet-1"},
		}),

		Entry("managed nodegroup with subnets that are in unsuppoted AZs", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:         "my-nodegroup",
					InstanceType: "g4dn.xlarge",
					Subnets:      []string{"subnet-1"},
				},
			},
			updateClusterConfig: func(config *api.ClusterConfig) {
				config.AvailabilityZones = []string{"us-west-2a", "us-west-2d"}
				config.VPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2d": api.AZSubnetSpec{
							ID: "subnet-1",
							AZ: "us-west-2d",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			updateEC2Mocks: func(e *mocksv2.EC2) {
				e.On("DescribeInstanceTypeOfferings", mock.Anything, mock.Anything, mock.Anything).
					Return(&ec2.DescribeInstanceTypeOfferingsOutput{
						InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
							{
								InstanceType: ec2types.InstanceTypeG4dnXlarge,
								Location:     aws.String("us-west-2a"),
								LocationType: ec2types.LocationTypeAvailabilityZone,
							},
							{
								InstanceType: api.DefaultNodeType,
								Location:     aws.String("us-west-2d"),
								LocationType: ec2types.LocationTypeAvailabilityZone,
							},
						},
					}, nil)
				e.On("DescribeAvailabilityZones", mock.Anything, mock.Anything).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []ec2types.AvailabilityZone{
							{
								ZoneType: aws.String("availability-zone"),
								ZoneName: aws.String("us-west-2d"),
							},
						},
					}, nil)
			},
			customInstanceSupport: true,
			expectedErr:           "failed to select subnet subnet-1: cannot create nodegroup my-nodegroup in availability zone us-west-2d as it does not support all required instance types",
		}),

		Entry("managed nodegroup with AZs that don't support all required instance types", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name:              "my-nodegroup",
					InstanceType:      "g4dn.xlarge",
					AvailabilityZones: []string{"us-west-2a", "us-west-2d"},
				},
			},
			updateClusterConfig: func(config *api.ClusterConfig) {
				config.AvailabilityZones = []string{"us-west-2a", "us-west-2d"}
				config.VPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID: "subnet-1",
							AZ: "us-west-2a",
						},
						"us-west-2d": api.AZSubnetSpec{
							ID: "subnet-2",
							AZ: "us-west-2d",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			updateEC2Mocks: func(e *mocksv2.EC2) {
				e.On("DescribeInstanceTypeOfferings", mock.Anything, mock.Anything, mock.Anything).
					Return(&ec2.DescribeInstanceTypeOfferingsOutput{
						InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
							{
								InstanceType: ec2types.InstanceTypeG4dnXlarge,
								Location:     aws.String("us-west-2a"),
								LocationType: ec2types.LocationTypeAvailabilityZone,
							},
							{
								InstanceType: api.DefaultNodeType,
								Location:     aws.String("us-west-2d"),
								LocationType: ec2types.LocationTypeAvailabilityZone,
							},
						},
					}, nil)
				e.On("DescribeAvailabilityZones", mock.Anything, mock.Anything).
					Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []ec2types.AvailabilityZone{
							{
								ZoneType: aws.String("availability-zone"),
								ZoneName: aws.String("us-west-2a"),
							},
							{
								ZoneType: aws.String("availability-zone"),
								ZoneName: aws.String("us-west-2d"),
							},
						},
					}, nil)
			},
			customInstanceSupport: true,
			expectedErr:           "cannot create nodegroup my-nodegroup in availability zone us-west-2d as it does not support all required instance types",
		}),

		Entry("EKS on Outposts but subnets not on Outposts", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:           []string{"subnet-123"},
					PrivateNetworking: true,
				},
			},
			availabilityZones: []string{"us-west-2a"},
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			},
			updateEC2Mocks: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2a", vpcID)
			},
			expectedErr: `subnet "subnet-123" is not on Outposts`,
		}),

		Entry("EKS on Outposts but subnets in a different Outpost", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:           []string{"subnet-123"},
					PrivateNetworking: true,
				},
			},
			availabilityZones: []string{"us-west-2a"},
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			},
			updateEC2Mocks: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnetsWithOutpost(ec2Mock, "us-west-2a", vpcID, aws.String("arn:aws:outposts:us-west-2:1234:outpost/op-5678"))
			},
			expectedErr: `subnet "subnet-123" is in a different Outpost ARN ("arn:aws:outposts:us-west-2:1234:outpost/op-5678") than the control plane or nodegroup Outpost ("arn:aws:outposts:us-west-2:1234:outpost/op-1234")`,
		}),

		Entry("EKS and subnets in the same Outpost", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:           []string{"subnet-123"},
					PrivateNetworking: true,
				},
			},
			availabilityZones: []string{"us-west-2a"},
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			},
			updateEC2Mocks: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnetsWithOutpost(ec2Mock, "us-west-2a", vpcID, aws.String("arn:aws:outposts:us-west-2:1234:outpost/op-1234"))
			},
			expectedSubnetIDs: []string{"subnet-123"},
		}),
	)
})

func mockDescribeSubnets(ec2Mock *mocksv2.EC2, zoneName, vpcID string) {
	mockDescribeSubnetsWithOutpost(ec2Mock, zoneName, vpcID, nil)
}

func mockDescribeSubnetsWithOutpost(ec2Mock *mocksv2.EC2, zoneName, vpcID string, outpostARN *string) {
	ec2Mock.On("DescribeSubnets", mock.Anything, mock.Anything, mock.Anything).Return(func(_ context.Context, input *ec2.DescribeSubnetsInput, _ ...func(options *ec2.Options)) *ec2.DescribeSubnetsOutput {
		return &ec2.DescribeSubnetsOutput{
			Subnets: []ec2types.Subnet{
				{
					SubnetId:         aws.String(input.SubnetIds[0]),
					AvailabilityZone: aws.String(zoneName),
					VpcId:            aws.String(vpcID),
					OutpostArn:       outpostARN,
				},
			},
		}
	}, nil)
}

func mockSubnetsAndAZInstanceSupport(
	cfg *api.ClusterConfig,
	provider *mockprovider.MockProvider,
	availabilityZones []string,
	localZones []string,
	instanceTypes []ec2types.InstanceType,
) {
	azs := []ec2types.AvailabilityZone{}
	offerings := []ec2types.InstanceTypeOffering{}

	publicSubnetMapping := api.AZSubnetMapping{}
	privateSubnetMapping := api.AZSubnetMapping{}
	for _, azName := range availabilityZones {
		publicSubnetMapping[azName] = api.AZSubnetSpec{
			ID: fmt.Sprintf("subnet-public-%s", azName),
			AZ: azName,
		}
		privateSubnetMapping[azName] = api.AZSubnetSpec{
			ID: fmt.Sprintf("subnet-private-%s", azName),
			AZ: azName,
		}
		azs = append(azs, ec2types.AvailabilityZone{
			ZoneType: aws.String("availability-zone"),
			ZoneName: aws.String(azName),
		})
		for _, instance := range instanceTypes {
			offerings = append(offerings, ec2types.InstanceTypeOffering{
				InstanceType: instance,
				Location:     aws.String(azName),
				LocationType: ec2types.LocationTypeAvailabilityZone,
			})
		}
	}
	cfg.AvailabilityZones = availabilityZones
	cfg.VPC.Subnets = &api.ClusterSubnets{
		Public:  publicSubnetMapping,
		Private: privateSubnetMapping,
	}

	publicSubnetMapping = api.AZSubnetMapping{}
	privateSubnetMapping = api.AZSubnetMapping{}
	for _, lzName := range localZones {
		publicSubnetMapping[lzName] = api.AZSubnetSpec{
			ID: fmt.Sprintf("subnet-public-%s", lzName),
			AZ: lzName,
		}
		privateSubnetMapping[lzName] = api.AZSubnetSpec{
			ID: fmt.Sprintf("subnet-private-%s", lzName),
			AZ: lzName,
		}
		azs = append(azs, ec2types.AvailabilityZone{
			ZoneType: aws.String("local-zone"),
			ZoneName: aws.String(lzName),
		})
		for _, instance := range instanceTypes {
			offerings = append(offerings, ec2types.InstanceTypeOffering{
				InstanceType: instance,
				Location:     aws.String(lzName),
				LocationType: ec2types.LocationTypeAvailabilityZone,
			})
		}
	}
	cfg.VPC.LocalZoneSubnets = &api.ClusterSubnets{
		Public:  publicSubnetMapping,
		Private: privateSubnetMapping,
	}

	provider.MockEC2().
		On("DescribeAvailabilityZones", mock.Anything, mock.Anything).
		Return(&ec2.DescribeAvailabilityZonesOutput{
			AvailabilityZones: azs,
		}, nil)
	provider.MockEC2().
		On("DescribeInstanceTypeOfferings", mock.Anything, mock.Anything, mock.Anything).
		Return(&ec2.DescribeInstanceTypeOfferingsOutput{
			InstanceTypeOfferings: offerings,
		}, nil)
}
