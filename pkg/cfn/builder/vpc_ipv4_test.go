package builder_test

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/builder/fakes"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

var _ = Describe("VPC Template Builder", func() {
	var (
		vpcRs   *builder.IPv4VPCResourceSet
		cfg     *api.ClusterConfig
		mockEC2 = &mocksv2.EC2{}
	)

	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		cfg.VPC = vpcConfig()
		cfg.AvailabilityZones = []string{azA, azB}
	})

	JustBeforeEach(func() {
		vpcRs = builder.NewIPv4VPCResourceSet(builder.NewRS(), cfg, mockEC2, false)
	})

	Describe("AddResources", func() {
		var (
			addErr        error
			vpcID         *gfnt.Value
			subnetDetails *builder.SubnetDetails
			vpcTemplate   *fakes.FakeTemplate
		)

		JustBeforeEach(func() {
			vpcID, subnetDetails, addErr = vpcRs.CreateTemplate(context.Background())
			vpcTemplate = &fakes.FakeTemplate{}
			templateBody, err := vpcRs.RenderJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(json.Unmarshal(templateBody, vpcTemplate)).To(Succeed())
		})

		It("should not error", func() {
			Expect(addErr).NotTo(HaveOccurred())
		})

		It("returns the VPC resource", func() {
			Expect(vpcID).To(Equal(gfnt.MakeRef(vpcResourceKey)))
		})

		It("adds the correct gateway resources to the resource set", func() {
			Expect(vpcTemplate.Resources).To(HaveKey(igwKey))
			Expect(vpcTemplate.Resources).To(HaveKey(gaKey))
			Expect(vpcTemplate.Resources[gaKey].Properties.InternetGatewayID).To(Equal(makeRef(igwKey)))
			Expect(vpcTemplate.Resources[gaKey].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
		})

		It("adds the public subnet routes to the resource set", func() {
			Expect(vpcTemplate.Resources).To(HaveKey(pubRouteTable))
			Expect(vpcTemplate.Resources[pubRouteTable].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))

			Expect(vpcTemplate.Resources).To(HaveKey(pubSubnetRoute))
			Expect(vpcTemplate.Resources[pubSubnetRoute].Properties.RouteTableID).To(Equal(makeRef(pubRouteTable)))
			Expect(vpcTemplate.Resources[pubSubnetRoute].Properties.DestinationCidrBlock).To(Equal("0.0.0.0/0"))
			Expect(vpcTemplate.Resources[pubSubnetRoute].Properties.GatewayID).To(Equal(makeRef(igwKey)))
		})

		It("returns public subnet settings", func() {
			Expect(subnetDetails.Public).To(ContainElement(builder.SubnetResource{
				Subnet:           gfnt.MakeRef(publicSubnetRef2),
				RouteTable:       gfnt.MakeRef(pubRouteTable),
				AvailabilityZone: azB,
			}))
			Expect(subnetDetails.Public).To(ContainElement(builder.SubnetResource{
				Subnet:           gfnt.MakeRef(publicSubnetRef1),
				RouteTable:       gfnt.MakeRef(pubRouteTable),
				AvailabilityZone: azA,
			}))
		})

		It("adds the public subnets to the resource set", func() {
			Expect(vpcTemplate.Resources).To(HaveKey(publicSubnetRef1))
			Expect(vpcTemplate.Resources[publicSubnetRef1].Properties.AvailabilityZone).To(Equal(azA))
			Expect(vpcTemplate.Resources[publicSubnetRef1].Properties.CidrBlock).To(Equal("192.168.32.0/19"))
			Expect(vpcTemplate.Resources[publicSubnetRef1].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
			Expect(vpcTemplate.Resources[publicSubnetRef1].Properties.Tags[0].Key).To(Equal("kubernetes.io/role/elb"))
			Expect(vpcTemplate.Resources[publicSubnetRef1].Properties.Tags[0].Value).To(Equal("1"))
			Expect(vpcTemplate.Resources[publicSubnetRef1].Properties.Tags[1].Key).To(Equal("Name"))
			Expect(vpcTemplate.Resources[publicSubnetRef1].Properties.Tags[1].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${AWS::StackName}/SubnetPublicUSWEST2A",
			}))
			Expect(vpcTemplate.Resources[publicSubnetRef1].Properties.MapPublicIPOnLaunch).To(BeTrue())

			Expect(vpcTemplate.Resources).To(HaveKey(publicSubnetRef2))
			Expect(vpcTemplate.Resources[publicSubnetRef2].Properties.AvailabilityZone).To(Equal(azB))
			Expect(vpcTemplate.Resources[publicSubnetRef2].Properties.CidrBlock).To(Equal("192.168.0.0/19"))
			Expect(vpcTemplate.Resources[publicSubnetRef2].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
			Expect(vpcTemplate.Resources[publicSubnetRef2].Properties.Tags[0].Key).To(Equal("kubernetes.io/role/elb"))
			Expect(vpcTemplate.Resources[publicSubnetRef2].Properties.Tags[0].Value).To(Equal("1"))
			Expect(vpcTemplate.Resources[publicSubnetRef2].Properties.Tags[1].Key).To(Equal("Name"))
			Expect(vpcTemplate.Resources[publicSubnetRef2].Properties.Tags[1].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${AWS::StackName}/SubnetPublicUSWEST2B",
			}))
			Expect(vpcTemplate.Resources[publicSubnetRef2].Properties.MapPublicIPOnLaunch).To(BeTrue())

			Expect(vpcTemplate.Resources).To(HaveKey(rtaPublicA))
			Expect(vpcTemplate.Resources[rtaPublicA].Properties.SubnetID).To(Equal(makeRef(publicSubnetRef1)))
			Expect(vpcTemplate.Resources[rtaPublicA].Properties.RouteTableID).To(Equal(makeRef(pubRouteTable)))
			Expect(vpcTemplate.Resources).To(HaveKey(rtaPublicB))
			Expect(vpcTemplate.Resources[rtaPublicB].Properties.SubnetID).To(Equal(makeRef(publicSubnetRef2)))
			Expect(vpcTemplate.Resources[rtaPublicB].Properties.RouteTableID).To(Equal(makeRef(pubRouteTable)))
		})

		It("returns private subnet settings", func() {
			Expect(subnetDetails.Private).To(HaveLen(2))
			Expect(subnetDetails.Private).To(ContainElement(builder.SubnetResource{
				Subnet:           gfnt.MakeRef(privateSubnetRef2),
				RouteTable:       gfnt.MakeRef(privRouteTableB),
				AvailabilityZone: azB,
			}))
			Expect(subnetDetails.Private).To(ContainElement(builder.SubnetResource{
				Subnet:           gfnt.MakeRef(privateSubnetRef1),
				RouteTable:       gfnt.MakeRef(privRouteTableA),
				AvailabilityZone: azA,
			}))
		})

		It("adds the private subnets to the resource set", func() {
			Expect(vpcTemplate.Resources).To(HaveKey(privateSubnetRef1))
			Expect(vpcTemplate.Resources[privateSubnetRef1].Properties.AvailabilityZone).To(Equal(azA))
			Expect(vpcTemplate.Resources[privateSubnetRef1].Properties.CidrBlock).To(Equal("192.168.128.0/19"))
			Expect(vpcTemplate.Resources[privateSubnetRef1].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
			Expect(vpcTemplate.Resources[privateSubnetRef1].Properties.Tags[0].Key).To(Equal("kubernetes.io/role/internal-elb"))
			Expect(vpcTemplate.Resources[privateSubnetRef1].Properties.Tags[0].Value).To(Equal("1"))
			Expect(vpcTemplate.Resources[privateSubnetRef1].Properties.Tags[1].Key).To(Equal("Name"))
			Expect(vpcTemplate.Resources[privateSubnetRef1].Properties.Tags[1].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${AWS::StackName}/SubnetPrivateUSWEST2A",
			}))

			Expect(vpcTemplate.Resources).To(HaveKey(privateSubnetRef2))
			Expect(vpcTemplate.Resources[privateSubnetRef2].Properties.AvailabilityZone).To(Equal(azB))
			Expect(vpcTemplate.Resources[privateSubnetRef2].Properties.CidrBlock).To(Equal("192.168.96.0/19"))
			Expect(vpcTemplate.Resources[privateSubnetRef2].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
			Expect(vpcTemplate.Resources[privateSubnetRef2].Properties.Tags[0].Key).To(Equal("kubernetes.io/role/internal-elb"))
			Expect(vpcTemplate.Resources[privateSubnetRef2].Properties.Tags[0].Value).To(Equal("1"))
			Expect(vpcTemplate.Resources[privateSubnetRef2].Properties.Tags[1].Key).To(Equal("Name"))
			Expect(vpcTemplate.Resources[privateSubnetRef2].Properties.Tags[1].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${AWS::StackName}/SubnetPrivateUSWEST2B",
			}))

			Expect(vpcTemplate.Resources).To(HaveKey(rtaPrivateA))
			Expect(vpcTemplate.Resources[rtaPrivateA].Properties.SubnetID).To(Equal(makeRef(privateSubnetRef1)))
			Expect(vpcTemplate.Resources[rtaPrivateA].Properties.RouteTableID).To(Equal(makeRef(privRouteTableA)))
			Expect(vpcTemplate.Resources).To(HaveKey(rtaPrivateB))
			Expect(vpcTemplate.Resources[rtaPrivateB].Properties.SubnetID).To(Equal(makeRef(privateSubnetRef2)))
			Expect(vpcTemplate.Resources[rtaPrivateB].Properties.RouteTableID).To(Equal(makeRef(privRouteTableB)))
		})

		Context("highly available nat is set", func() {
			BeforeEach(func() {
				*cfg.VPC.NAT.Gateway = api.ClusterHighlyAvailableNAT
			})

			It("adds HA nat gateway resources to the resource set", func() {
				Expect(vpcTemplate.Resources).To(HaveKey("NATIPUSWEST2A"))
				Expect(vpcTemplate.Resources["NATIPUSWEST2A"].Properties.Domain).To(Equal("vpc"))
				Expect(vpcTemplate.Resources).To(HaveKey("NATIPUSWEST2B"))
				Expect(vpcTemplate.Resources["NATIPUSWEST2B"].Properties.Domain).To(Equal("vpc"))

				Expect(vpcTemplate.Resources).To(HaveKey("NATGatewayUSWEST2A"))
				Expect(vpcTemplate.Resources["NATGatewayUSWEST2A"].Properties.AllocationID).To(Equal(makeGetAttr("NATIPUSWEST2A", "AllocationId")))
				Expect(vpcTemplate.Resources["NATGatewayUSWEST2A"].Properties.SubnetID).To(Equal(makeRef(publicSubnetRef1)))

				Expect(vpcTemplate.Resources).To(HaveKey("NATGatewayUSWEST2B"))
				Expect(vpcTemplate.Resources["NATGatewayUSWEST2B"].Properties.AllocationID).To(Equal(makeGetAttr("NATIPUSWEST2B", "AllocationId")))
				Expect(vpcTemplate.Resources["NATGatewayUSWEST2B"].Properties.SubnetID).To(Equal(makeRef(publicSubnetRef2)))

				Expect(vpcTemplate.Resources).To(HaveKey(privRouteTableA))
				Expect(vpcTemplate.Resources[privRouteTableA].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
				Expect(vpcTemplate.Resources).To(HaveKey(privRouteTableB))
				Expect(vpcTemplate.Resources[privRouteTableB].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))

				Expect(vpcTemplate.Resources).To(HaveKey("NATPrivateSubnetRouteUSWEST2A"))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2A"].Properties.RouteTableID).To(Equal(makeRef(privRouteTableA)))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2A"].Properties.DestinationCidrBlock).To(Equal("0.0.0.0/0"))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2A"].Properties.NatGatewayID).To(Equal(makeRef("NATGatewayUSWEST2A")))
				Expect(vpcTemplate.Resources).To(HaveKey("NATPrivateSubnetRouteUSWEST2B"))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2B"].Properties.RouteTableID).To(Equal(makeRef(privRouteTableB)))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2B"].Properties.DestinationCidrBlock).To(Equal("0.0.0.0/0"))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2B"].Properties.NatGatewayID).To(Equal(makeRef("NATGatewayUSWEST2B")))

				Expect(vpcTemplate.Resources).To(HaveKey(rtaPublicA))
				Expect(vpcTemplate.Resources).To(HaveKey(rtaPublicB))
			})
		})

		Context("single nat is set", func() {
			BeforeEach(func() {
				*cfg.VPC.NAT.Gateway = api.ClusterSingleNAT
			})

			It("adds HA nat gateway resources to the resource set", func() {
				Expect(vpcTemplate.Resources).To(HaveKey("NATIP"))
				Expect(vpcTemplate.Resources["NATIP"].Properties.Domain).To(Equal("vpc"))

				Expect(vpcTemplate.Resources).To(HaveKey("NATGateway"))
				Expect(vpcTemplate.Resources["NATGateway"].Properties.AllocationID).To(Equal(makeGetAttr("NATIP", "AllocationId")))
				Expect(vpcTemplate.Resources["NATGateway"].Properties.SubnetID).To(Equal(makeRef(publicSubnetRef1)))

				Expect(vpcTemplate.Resources).To(HaveKey(privRouteTableA))
				Expect(vpcTemplate.Resources[privRouteTableA].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
				Expect(vpcTemplate.Resources).To(HaveKey(privRouteTableB))
				Expect(vpcTemplate.Resources[privRouteTableB].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))

				Expect(vpcTemplate.Resources).To(HaveKey("NATPrivateSubnetRouteUSWEST2A"))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2A"].Properties.RouteTableID).To(Equal(makeRef(privRouteTableA)))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2A"].Properties.DestinationCidrBlock).To(Equal("0.0.0.0/0"))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2A"].Properties.NatGatewayID).To(Equal(makeRef("NATGateway")))
				Expect(vpcTemplate.Resources).To(HaveKey("NATPrivateSubnetRouteUSWEST2B"))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2B"].Properties.RouteTableID).To(Equal(makeRef(privRouteTableB)))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2B"].Properties.DestinationCidrBlock).To(Equal("0.0.0.0/0"))
				Expect(vpcTemplate.Resources["NATPrivateSubnetRouteUSWEST2B"].Properties.NatGatewayID).To(Equal(makeRef("NATGateway")))

				Expect(vpcTemplate.Resources).To(HaveKey(rtaPublicA))
				Expect(vpcTemplate.Resources).To(HaveKey(rtaPublicB))
			})
		})

		Context("nat is disabled", func() {
			BeforeEach(func() {
				*cfg.VPC.NAT.Gateway = api.ClusterDisableNAT
			})

			It("adds HA nat gateway resources to the resource set", func() {
				Expect(vpcTemplate.Resources).NotTo(HaveKey("NATIP"))
				Expect(vpcTemplate.Resources).NotTo(HaveKey("NATGateway"))
				Expect(vpcTemplate.Resources).NotTo(HaveKey("NATPrivateSubnetRouteUSWEST2A"))
				Expect(vpcTemplate.Resources).NotTo(HaveKey("NATPrivateSubnetRouteUSWEST2B"))

				Expect(vpcTemplate.Resources).To(HaveKey(privRouteTableA))
				Expect(vpcTemplate.Resources[privRouteTableA].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
				Expect(vpcTemplate.Resources).To(HaveKey(privRouteTableB))
				Expect(vpcTemplate.Resources[privRouteTableB].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))

				Expect(vpcTemplate.Resources).To(HaveKey(rtaPublicA))
				Expect(vpcTemplate.Resources).To(HaveKey(rtaPublicB))
			})
		})

		Context("an invalid nat option is set", func() {
			BeforeEach(func() {
				*cfg.VPC.NAT.Gateway = "some-trash"
			})

			It("returns an error", func() {
				Expect(addErr).To(MatchError("some-trash is not a valid NAT gateway mode"))
			})
		})

		Context("when AutoAllocateIPv6 is enabled", func() {
			var expectedFnCIDR string
			BeforeEach(func() {
				autoAllocated := true
				cfg.VPC.AutoAllocateIPv6 = &autoAllocated
				expectedFnCIDR = `{ "Fn::Cidr": [{ "Fn::Select": [ 0, { "Fn::GetAtt": ["VPC", "Ipv6CidrBlocks"] }]}, 6, 64 ]}`
			})

			It("adds the AutoAllocatedCIDRv6 vpc resource to the resource set", func() {
				Expect(vpcTemplate.Resources).To(HaveKey("AutoAllocatedCIDRv6"))
				Expect(vpcTemplate.Resources["AutoAllocatedCIDRv6"].Properties.AmazonProvidedIpv6CidrBlock).To(BeTrue())
				Expect(vpcTemplate.Resources["AutoAllocatedCIDRv6"].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
			})

			It("adds the correct subnet resources to the resource set", func() {
				Expect(vpcTemplate.Resources).To(HaveKey("PublicUSWEST2ACIDRv6"))
				Expect(vpcTemplate.Resources["PublicUSWEST2ACIDRv6"].Properties.SubnetID).To(Equal(makeRef(publicSubnetRef1)))
				assertIpv6CidrBlockCreatedWithSelect(vpcTemplate.Resources["PublicUSWEST2BCIDRv6"].Properties.Ipv6CidrBlock, expectedFnCIDR)
				assertIpv6CidrBlockCreatedWithSelect(vpcTemplate.Resources["PublicUSWEST2ACIDRv6"].Properties.Ipv6CidrBlock, expectedFnCIDR)
				assertIpv6CidrBlockCreatedWithSelect(vpcTemplate.Resources["PrivateUSWEST2BCIDRv6"].Properties.Ipv6CidrBlock, expectedFnCIDR)
				assertIpv6CidrBlockCreatedWithSelect(vpcTemplate.Resources["PrivateUSWEST2ACIDRv6"].Properties.Ipv6CidrBlock, expectedFnCIDR)
			})

			It("adds the ipv6 route table entry", func() {
				Expect(vpcTemplate.Resources).To(HaveKey(pubRouteTable))
				Expect(vpcTemplate.Resources[pubRouteTable].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))

				Expect(vpcTemplate.Resources).To(HaveKey(builder.PubSubIPv6RouteKey))
				Expect(vpcTemplate.Resources[builder.PubSubIPv6RouteKey].Properties.RouteTableID).To(Equal(makeRef(pubRouteTable)))
				Expect(vpcTemplate.Resources[builder.PubSubIPv6RouteKey].Properties.DestinationIpv6CidrBlock).To(Equal(builder.InternetIPv6CIDR))
				Expect(vpcTemplate.Resources[builder.PubSubIPv6RouteKey].Properties.GatewayID).To(Equal(makeRef(igwKey)))
			})
		})

		Context("when the vpc is fully private", func() {
			BeforeEach(func() {
				cfg.PrivateCluster.Enabled = true
			})

			It("disables the nat", func() {
				Expect(vpcTemplate.Resources).NotTo(HaveKey("NATIP"))
				Expect(vpcTemplate.Resources).NotTo(HaveKey("NATGateway"))
			})

			It("does not add an internet gateway", func() {
				Expect(vpcTemplate.Resources).NotTo(HaveKey(igwKey))
			})

			It("does not set public subnet resources", func() {
				Expect(subnetDetails.Public).To(HaveLen(0))
				Expect(vpcTemplate.Resources).NotTo(HaveKey(pubSubnetRoute))
				Expect(vpcTemplate.Resources).NotTo(HaveKey(pubSubnetRoute))
				Expect(vpcTemplate.Resources).NotTo(HaveKey(publicSubnetRef1))
				Expect(vpcTemplate.Resources).NotTo(HaveKey(publicSubnetRef1))
				Expect(vpcTemplate.Resources).NotTo(HaveKey(rtaPublicA))
				Expect(vpcTemplate.Resources).NotTo(HaveKey(rtaPublicB))

				Expect(subnetDetails.Private).To(HaveLen(2))
				Expect(vpcTemplate.Resources).To(HaveKey(privRouteTableA))
				Expect(vpcTemplate.Resources).To(HaveKey(privRouteTableB))
				Expect(vpcTemplate.Resources).To(HaveKey(privateSubnetRef1))
				Expect(vpcTemplate.Resources).To(HaveKey(privateSubnetRef2))
				Expect(vpcTemplate.Resources).To(HaveKey(rtaPrivateA))
				Expect(vpcTemplate.Resources).To(HaveKey(rtaPrivateB))

				Expect(vpcTemplate.Resources[privRouteTableA].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
				Expect(vpcTemplate.Resources[rtaPrivateA].Properties.SubnetID).To(Equal(makeRef(privateSubnetRef1)))
				Expect(vpcTemplate.Resources[rtaPrivateA].Properties.RouteTableID).To(Equal(makeRef(privRouteTableA)))
				Expect(vpcTemplate.Resources[privRouteTableB].Properties.VpcID).To(Equal(makeRef(vpcResourceKey)))
				Expect(vpcTemplate.Resources[rtaPrivateB].Properties.SubnetID).To(Equal(makeRef(privateSubnetRef2)))
				Expect(vpcTemplate.Resources[rtaPrivateB].Properties.RouteTableID).To(Equal(makeRef(privRouteTableB)))
			})
		})
	})

	Describe("AddOutputs", func() {
		var (
			vpcTemplate *fakes.FakeTemplate
		)

		JustBeforeEach(func() {
			_, _, err := vpcRs.CreateTemplate(context.Background())
			Expect(err).NotTo(HaveOccurred())
			vpcTemplate = &fakes.FakeTemplate{}
			templateBody, err := vpcRs.RenderJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(json.Unmarshal(templateBody, vpcTemplate)).To(Succeed())
		})

		Context("simple config, no nat no subnets", func() {
			BeforeEach(func() {
				*cfg.VPC.NAT.Gateway = api.ClusterDisableNAT
				cfg.VPC.Subnets.Private = nil
				cfg.VPC.Subnets.Public = nil
			})

			It("adds the cluster vpc reference to the outputs", func() {
				Expect(vpcTemplate.Outputs).To(HaveKey(vpcResourceKey))
			})
		})

		Context("if NAT is not nil", func() {
			It("adds the nat mode and gateway to the outputs", func() {
				Expect(vpcTemplate.Outputs).To(HaveKey("FeatureNATMode"))
			})
		})

		Context("if there are subnets", func() {
			It("adds the subnet refs to the output", func() {
				Expect(vpcTemplate.Outputs).To(HaveKey("SubnetsPublic"))
				Expect(vpcTemplate.Outputs).To(HaveKey("SubnetsPrivate"))
			})
		})

		Context("the cluster is fully private", func() {
			BeforeEach(func() {
				cfg.PrivateCluster.Enabled = true
			})

			It("adds the fully private output", func() {
				Expect(vpcTemplate.Outputs).To(HaveKey("ClusterFullyPrivate"))
			})
		})
	})

	Describe("PublicSubnetRefs", func() {
		It("returns the references of public subnets", func() {
			_, subnetDetails, err := vpcRs.CreateTemplate(context.Background())
			Expect(err).NotTo(HaveOccurred())
			refs := subnetDetails.PublicSubnetRefs()
			Expect(refs).To(HaveLen(2))
			Expect(refs).To(ContainElement(makePrimitive(publicSubnetRef1)))
			Expect(refs).To(ContainElement(makePrimitive(publicSubnetRef2)))
		})
	})

	Describe("PrivateSubnetRefs", func() {
		It("returns the references of private subnets", func() {
			_, subnetDetails, err := vpcRs.CreateTemplate(context.Background())
			Expect(err).NotTo(HaveOccurred())
			refs := subnetDetails.PrivateSubnetRefs()
			Expect(refs).To(HaveLen(2))
			Expect(refs).To(ContainElement(makePrimitive(privateSubnetRef1)))
			Expect(refs).To(ContainElement(makePrimitive(privateSubnetRef2)))
		})
	})
})

func makePrimitive(primitive string) *gfnt.Value {
	output, err := gfnt.NewValueFromPrimitive(makeRef(primitive))
	Expect(err).NotTo(HaveOccurred())
	return output
}

func makeRef(value string) map[string]interface{} {
	return map[string]interface{}{"Ref": value}
}

func makeGetAttr(values ...interface{}) map[string]interface{} {
	return map[string]interface{}{"Fn::GetAtt": values}
}

func makeRTOutput(subnetIds []string, main bool) *ec2.DescribeRouteTablesOutput {
	return &ec2.DescribeRouteTablesOutput{
		RouteTables: []ec2types.RouteTable{{
			RouteTableId: aws.String("this-is-a-route-table"),
			Associations: []ec2types.RouteTableAssociation{{
				SubnetId: aws.String(subnetIds[0]),
				Main:     aws.Bool(main),
			}, {
				SubnetId: aws.String(subnetIds[1]),
				Main:     aws.Bool(main),
			}},
		}},
	}
}

func assertIpv6CidrBlockCreatedWithSelect(cidrBlock interface{}, expectedFnCIDR string) {
	ExpectWithOffset(1, cidrBlock.(map[string]interface{})).To(HaveKey("Fn::Select"))
	fnSelectValue := cidrBlock.(map[string]interface{})["Fn::Select"].([]interface{})
	ExpectWithOffset(1, fnSelectValue).To(HaveLen(2))
	ExpectWithOffset(1, fnSelectValue[0].(float64)).To(BeNumerically("~", 0, 8))
	actualFnCIDR, err := json.Marshal(fnSelectValue[1])
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ExpectWithOffset(1, actualFnCIDR).To(MatchJSON([]byte(expectedFnCIDR)))
}
