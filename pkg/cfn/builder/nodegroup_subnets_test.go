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
	"github.com/weaveworks/eksctl/pkg/vpc"
	vpcfakes "github.com/weaveworks/eksctl/pkg/vpc/fakes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

var _ = Describe("AssignSubnets", func() {
	type assignSubnetsEntry struct {
		np                  api.NodePool
		mockEC2             func(*mocksv2.EC2)
		setSubnetMapping    func(*api.ClusterVPC)
		updateClusterConfig func(config *api.ClusterConfig)
		createVPCImporter   func() vpc.Importer

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
		if e.setSubnetMapping != nil {
			e.setSubnetMapping(clusterConfig.VPC)
		}
		mockProvider := mockprovider.NewMockProvider()
		if e.mockEC2 != nil {
			e.mockEC2(mockProvider.MockEC2())
		}
		if e.updateClusterConfig != nil {
			e.updateClusterConfig(clusterConfig)
		}

		var vpcImporter vpc.Importer
		if e.createVPCImporter != nil {
			vpcImporter = e.createVPCImporter()
		}
		subnetRefs, err := builder.AssignSubnets(context.Background(), e.np, vpcImporter, clusterConfig, mockProvider.EC2())
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
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-1a", "subnet-1b", "subnet-1c"},
		}),

		Entry("managed nodegroup with availability zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-1a", "subnet-1b", "subnet-1c"},
		}),

		Entry("self-managed nodegroup with local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{},
				LocalZones:    []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.LocalZoneSubnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2-lax-1a": api.AZSubnetSpec{
							ID: "subnet-lax-1a",
							AZ: "us-west-2-lax-1a",
						},
						"us-west-2-lax-1b": api.AZSubnetSpec{
							ID: "subnet-lax-1b",
							AZ: "us-west-2-lax-1b",
						},
						"us-west-2-lax-1d": api.AZSubnetSpec{
							ID: "subnet-lax-1d",
							AZ: "us-west-2-lax-1d",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},

			expectedSubnetIDs: []string{"subnet-lax-1a", "subnet-lax-1b"},
		}),

		Entry("self-managed nodegroup with privateNetworking and local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					PrivateNetworking: true,
				},
				LocalZones: []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.LocalZoneSubnets = &api.ClusterSubnets{
					Public: api.NewAZSubnetMapping(),
					Private: api.AZSubnetMapping{
						"us-west-2-lax-1a": api.AZSubnetSpec{
							ID: "subnet-lax-1a",
							AZ: "us-west-2-lax-1a",
						},
						"us-west-2-lax-1b": api.AZSubnetSpec{
							ID: "subnet-lax-1b",
							AZ: "us-west-2-lax-1b",
						},
						"us-west-2-lax-1d": api.AZSubnetSpec{
							ID: "subnet-lax-1d",
							AZ: "us-west-2-lax-1d",
						},
					},
				}
			},

			expectedSubnetIDs: []string{"subnet-lax-1a", "subnet-lax-1b"},
		}),

		Entry("self-managed nodegroup with local zones and subnet IDs", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets: []string{"subnet-z1", "subnet-z2"},
				},
				LocalZones: []string{"us-west-2-lax-1a", "us-west-2-lax-1b"},
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.LocalZoneSubnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-2-lax-1a": api.AZSubnetSpec{
							ID: "subnet-lax-1a",
							AZ: "us-west-2-lax-1a",
						},
						"us-west-2-lax-1b": api.AZSubnetSpec{
							ID: "subnet-lax-1b",
							AZ: "us-west-2-lax-1b",
						},
						"us-west-2-lax-1d": api.AZSubnetSpec{
							ID: "subnet-lax-1d",
							AZ: "us-west-2-lax-1d",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-z1", "subnet-z2", "subnet-lax-1a", "subnet-lax-1b"},

			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2-lax-1e", vpcID)
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1e"),
					},
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2d"),
					},
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1f"),
					},
				})

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
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Private: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Public: api.NewAZSubnetMapping(),
				}
			},
			expectedSubnetIDs: []string{"subnet-1a", "subnet-1b", "subnet-1c", "subnet-z1", "subnet-z2"},

			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2g", vpcID)
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1e"),
					},
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2g"),
					},
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2h"),
					},
				})
			},
		}),

		Entry("managed nodegroup with availability zones and subnet IDs in local zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:           []string{"subnet-z1", "subnet-z2"},
					AvailabilityZones: []string{"us-west-1a", "us-west-1b", "us-west-1c"},
				},
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},

			expectedErr: "managed nodegroups cannot be launched in local zones",

			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2-lax-1e", vpcID)
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1e"),
					},
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2d"),
					},
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1f"),
					},
				})
			},
		}),

		Entry("managed nodegroup without subnets, availability zones and local zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{},
			},
			createVPCImporter: func() vpc.Importer {
				vpcImporter := new(vpcfakes.FakeImporter)
				vpcImporter.SubnetsPublicReturns(gfnt.NewStringSlice("subnet-ref1", "subnet-ref2"))
				return vpcImporter
			},
			expectedSubnetIDs: []string{"subnet-ref1", "subnet-ref2"},
		}),

		Entry("private self-managed nodegroup without subnets, availability zones and local zones", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					PrivateNetworking: true,
				},
			},
			createVPCImporter: func() vpc.Importer {
				vpcImporter := new(vpcfakes.FakeImporter)
				vpcImporter.SubnetsPrivateReturns(gfnt.NewStringSlice("subnet-pref1", "subnet-pref2"))
				return vpcImporter
			},

			expectedSubnetIDs: []string{"subnet-pref1", "subnet-pref2"},
		}),

		Entry("supplied subnet ID exists in a different VPC", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets: []string{"subnet-1"},
				},
			},

			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Public:  api.NewAZSubnetMapping(),
					Private: api.NewAZSubnetMapping(),
				}
			},

			expectedErr: `subnet with ID "subnet-1" is not in the attached VPC with ID "vpc-1"`,

			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2a", "vpc-2")
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1e"),
					},
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2d"),
					},
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1f"),
					},
				})
			},
		}),

		Entry("No subnets in zones", assignSubnetsEntry{
			np: &api.ManagedNodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					AvailabilityZones: []string{"us-west-2z"},
				},
			},

			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Public: api.AZSubnetMapping{
						"us-west-1a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-1a",
						},
						"us-west-1b": api.AZSubnetSpec{
							ID: "subnet-1b",
							AZ: "us-west-1b",
						},
						"us-west-1c": api.AZSubnetSpec{
							ID: "subnet-1c",
							AZ: "us-west-1c",
						},
					},
					Private: api.NewAZSubnetMapping(),
				}
			},

			expectedErr: "could not find public subnets for zones",
		}),

		Entry("EFA enabled with multiple subnets selects only one subnet", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:    []string{"subnet-1", "subnet-2", "subnet-3"},
					EFAEnabled: aws.Bool(true),
				},
			},

			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Public:  api.NewAZSubnetMapping(),
					Private: api.NewAZSubnetMapping(),
				}
			},

			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2d", vpcID)
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1e"),
					},
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2d"),
					},
					{
						ZoneType: aws.String("local-zone"),
						ZoneName: aws.String("us-west-2-lax-1f"),
					},
				})
			},

			expectedSubnetIDs: []string{"subnet-1"},
		}),

		Entry("nodegroup with subnet names", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets: []string{"subnet-1", "subnet-2", "subnet-3"},
				},
			},
			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2a"),
					},
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2d"),
					},
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2b"),
					},
				})
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
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

		Entry("EKS on Outposts but subnets not on Outposts", assignSubnetsEntry{
			np: &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Subnets:           []string{"subnet-123"},
					PrivateNetworking: true,
				},
			},
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-2a",
						},
					},
					Public: api.NewAZSubnetMapping(),
				}
			},
			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnets(ec2Mock, "us-west-2a", vpcID)
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2a"),
					},
				})
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
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-2a",
						},
					},
					Public: api.NewAZSubnetMapping(),
				}
			},
			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnetsWithOutpost(ec2Mock, "us-west-2a", vpcID, aws.String("arn:aws:outposts:us-west-2:1234:outpost/op-5678"))
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2a"),
					},
				})
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
			updateClusterConfig: func(c *api.ClusterConfig) {
				c.Outpost = &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				}
			},
			setSubnetMapping: func(clusterVPC *api.ClusterVPC) {
				clusterVPC.Subnets = &api.ClusterSubnets{
					Private: api.AZSubnetMapping{
						"us-west-2a": api.AZSubnetSpec{
							ID: "subnet-1a",
							AZ: "us-west-2a",
						},
					},
					Public: api.NewAZSubnetMapping(),
				}
			},
			mockEC2: func(ec2Mock *mocksv2.EC2) {
				mockDescribeSubnetsWithOutpost(ec2Mock, "us-west-2a", vpcID, aws.String("arn:aws:outposts:us-west-2:1234:outpost/op-1234"))
				mockDescribeAZs(ec2Mock, []ec2types.AvailabilityZone{
					{
						ZoneType: aws.String("availability-zone"),
						ZoneName: aws.String("us-west-2a"),
					},
				})
			},
			expectedSubnetIDs: []string{"subnet-123"},
		}),
	)
})

func mockDescribeSubnets(ec2Mock *mocksv2.EC2, zoneName, vpcID string) {
	mockDescribeSubnetsWithOutpost(ec2Mock, zoneName, vpcID, nil)
}

func mockDescribeSubnetsWithOutpost(ec2Mock *mocksv2.EC2, zoneName, vpcID string, outpostARN *string) {
	ec2Mock.On("DescribeSubnets", mock.Anything, mock.Anything).Return(func(_ context.Context, input *ec2.DescribeSubnetsInput, _ ...func(options *ec2.Options)) *ec2.DescribeSubnetsOutput {
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

func mockDescribeAZs(ec2Mock *mocksv2.EC2, zones []ec2types.AvailabilityZone) {
	ec2Mock.
		On("DescribeAvailabilityZones", mock.Anything, mock.Anything).Return(&ec2.DescribeAvailabilityZonesOutput{
		AvailabilityZones: zones,
	}, nil)
}
