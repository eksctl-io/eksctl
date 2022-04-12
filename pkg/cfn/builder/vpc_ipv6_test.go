package builder_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/builder/fakes"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

var _ = Describe("IPv6 VPC builder", func() {
	var (
		cfg *api.ClusterConfig
	)

	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
		cfg.AvailabilityZones = []string{azA, azB}
	})

	It("creates the ipv6 VPC and its resources", func() {
		vpcRs := builder.NewIPv6VPCResourceSet(builder.NewRS(), cfg, nil)
		_, subnetDetails, err := vpcRs.CreateTemplate(context.Background())
		Expect(err).NotTo(HaveOccurred())

		By("returning the references of public subnets")
		pubRefs := subnetDetails.PublicSubnetRefs()
		Expect(pubRefs).To(HaveLen(2))
		Expect(pubRefs).To(ContainElement(makePrimitive(builder.PublicSubnetKey + azAFormatted)))
		Expect(pubRefs).To(ContainElement(makePrimitive(builder.PublicSubnetKey + azBFormatted)))

		By("returning the references of private subnets")
		privRef := subnetDetails.PrivateSubnetRefs()
		Expect(privRef).To(HaveLen(2))
		Expect(privRef).To(ContainElement(makePrimitive(builder.PrivateSubnetKey + azBFormatted)))
		Expect(privRef).To(ContainElement(makePrimitive(builder.PrivateSubnetKey + azBFormatted)))

		vpcTemplate, err := renderTemplate(vpcRs)
		Expect(err).NotTo(HaveOccurred())

		By("creating the VPC resource")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.VPCResourceKey))
		Expect(vpcTemplate.Resources[builder.VPCResourceKey].Type).To(Equal("AWS::EC2::VPC"))
		defaultCidr := api.DefaultCIDR()
		cidr := &defaultCidr
		Expect(vpcTemplate.Resources[builder.VPCResourceKey].Properties).To(Equal(fakes.Properties{
			CidrBlock:          cidr.String(),
			EnableDNSHostnames: true,
			EnableDNSSupport:   true,
			Tags: []fakes.Tag{
				{
					Key:   "Name",
					Value: map[string]interface{}{"Fn::Sub": "${AWS::StackName}/VPC"},
				},
			},
		}))

		By("creating the IPv6 CIDR")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.IPv6CIDRBlockKey))
		Expect(vpcTemplate.Resources[builder.IPv6CIDRBlockKey].Type).To(Equal("AWS::EC2::VPCCidrBlock"))
		Expect(vpcTemplate.Resources[builder.IPv6CIDRBlockKey].Properties).To(Equal(fakes.Properties{
			AmazonProvidedIpv6CidrBlock: true,
			VpcID:                       map[string]interface{}{"Ref": "VPC"},
		}))

		By("creating the internet gateway")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.IGWKey))
		Expect(vpcTemplate.Resources[builder.IGWKey].Type).To(Equal("AWS::EC2::InternetGateway"))
		Expect(vpcTemplate.Resources[builder.IGWKey].Properties).To(Equal(fakes.Properties{
			Tags: []fakes.Tag{
				{
					Key:   "Name",
					Value: map[string]interface{}{"Fn::Sub": "${AWS::StackName}/InternetGateway"},
				},
			},
		}))

		By("creating a VPC gateway attachment to associate the IGW with the VPC")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.GAKey))
		Expect(vpcTemplate.Resources[builder.GAKey].Type).To(Equal("AWS::EC2::VPCGatewayAttachment"))
		Expect(vpcTemplate.Resources[builder.GAKey].Properties).To(Equal(fakes.Properties{
			InternetGatewayID: map[string]interface{}{"Ref": "InternetGateway"},
			VpcID:             map[string]interface{}{"Ref": "VPC"},
		}))

		By("creating a VPC gateway attachment to associate the IGW with the VPC")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.EgressOnlyInternetGatewayKey))
		Expect(vpcTemplate.Resources[builder.EgressOnlyInternetGatewayKey].Type).To(Equal("AWS::EC2::EgressOnlyInternetGateway"))
		Expect(vpcTemplate.Resources[builder.EgressOnlyInternetGatewayKey].Properties).To(Equal(fakes.Properties{
			VpcID: map[string]interface{}{"Ref": "VPC"},
		}))

		By("creating the NAT gateway")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.NATGatewayKey))
		Expect(vpcTemplate.Resources[builder.NATGatewayKey].Type).To(Equal("AWS::EC2::NatGateway"))
		Expect(vpcTemplate.Resources[builder.NATGatewayKey].DependsOn).To(ConsistOf(builder.ElasticIPKey, builder.PublicSubnetKey+azAFormatted, builder.GAKey))
		Expect(vpcTemplate.Resources[builder.NATGatewayKey].Properties).To(Equal(fakes.Properties{
			AllocationID: map[string]interface{}{
				"Fn::GetAtt": []interface{}{
					builder.ElasticIPKey,
					"AllocationId",
				},
			},
			SubnetID: map[string]interface{}{"Ref": builder.PublicSubnetKey + azAFormatted},
			Tags: []fakes.Tag{
				{
					Key:   "Name",
					Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", builder.NATGatewayKey)},
				},
			},
		}))

		By("creating an Elastic IP for the Nat Gateway")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.ElasticIPKey))
		Expect(vpcTemplate.Resources[builder.ElasticIPKey].Type).To(Equal("AWS::EC2::EIP"))
		Expect(vpcTemplate.Resources[builder.ElasticIPKey].DependsOn).To(ConsistOf(gaKey))
		Expect(vpcTemplate.Resources[builder.ElasticIPKey].Properties).To(Equal(fakes.Properties{
			Domain: "vpc",
			Tags: []fakes.Tag{
				{
					Key:   "Name",
					Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", builder.ElasticIPKey)},
				},
			},
		}))

		By("creating a public Route Table")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.PubRouteTableKey))
		Expect(vpcTemplate.Resources[builder.PubRouteTableKey].Type).To(Equal("AWS::EC2::RouteTable"))
		Expect(vpcTemplate.Resources[builder.PubRouteTableKey].Properties).To(Equal(fakes.Properties{
			VpcID: map[string]interface{}{"Ref": "VPC"},
			Tags: []fakes.Tag{
				{
					Key:   "Name",
					Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", builder.PubRouteTableKey)},
				},
			},
		}))

		By("creating public subnet route for IPv4 traffic to IPv4 CIDR")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.PubSubRouteKey))
		Expect(vpcTemplate.Resources[builder.PubSubRouteKey].Type).To(Equal("AWS::EC2::Route"))
		Expect(vpcTemplate.Resources[builder.PubSubRouteKey].DependsOn).To(ConsistOf(builder.GAKey))
		Expect(vpcTemplate.Resources[builder.PubSubRouteKey].Properties).To(Equal(fakes.Properties{
			DestinationCidrBlock: builder.InternetCIDR,
			GatewayID:            map[string]interface{}{"Ref": builder.IGWKey},
			RouteTableID:         map[string]interface{}{"Ref": builder.PubRouteTableKey},
		}))

		By("creating public subnet route for IPv6 traffic to IPv6 CIDR")
		Expect(vpcTemplate.Resources).To(HaveKey(builder.PubSubIPv6RouteKey))
		Expect(vpcTemplate.Resources[builder.PubSubIPv6RouteKey].Type).To(Equal("AWS::EC2::Route"))
		//TODO: we added this, wasn't in the example template. We think its correct?
		Expect(vpcTemplate.Resources[builder.PubSubIPv6RouteKey].DependsOn).To(ConsistOf(builder.GAKey))
		Expect(vpcTemplate.Resources[builder.PubSubIPv6RouteKey].Properties).To(Equal(fakes.Properties{
			DestinationIpv6CidrBlock: builder.InternetIPv6CIDR,
			GatewayID:                map[string]interface{}{"Ref": builder.IGWKey},
			RouteTableID:             map[string]interface{}{"Ref": builder.PubRouteTableKey},
		}))

		By("creating a private route table for each AZ")
		privateRouteTableA := builder.PrivateRouteTableKey + azAFormatted
		Expect(vpcTemplate.Resources).To(HaveKey(privateRouteTableA))
		Expect(vpcTemplate.Resources[privateRouteTableA].Type).To(Equal("AWS::EC2::RouteTable"))
		Expect(vpcTemplate.Resources[privateRouteTableA].Properties).To(Equal(fakes.Properties{
			VpcID: map[string]interface{}{"Ref": "VPC"},
			Tags: []fakes.Tag{
				{
					Key:   "Name",
					Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", privateRouteTableA)},
				},
			},
		}))
		privateRouteTableB := builder.PrivateRouteTableKey + azBFormatted
		Expect(vpcTemplate.Resources).To(HaveKey(privateRouteTableB))
		Expect(vpcTemplate.Resources[privateRouteTableB].Type).To(Equal("AWS::EC2::RouteTable"))
		Expect(vpcTemplate.Resources[privateRouteTableB].Properties).To(Equal(fakes.Properties{
			VpcID: map[string]interface{}{"Ref": "VPC"},
			Tags: []fakes.Tag{
				{
					Key:   "Name",
					Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", privateRouteTableB)},
				},
			},
		}))

		By("creating a route to the NAT gateway for each private subnet in the AZs")
		privateRouteA := builder.PrivateSubnetRouteKey + azAFormatted
		Expect(vpcTemplate.Resources).To(HaveKey(privateRouteA))
		Expect(vpcTemplate.Resources[privateRouteA].Type).To(Equal("AWS::EC2::Route"))
		Expect(vpcTemplate.Resources[privateRouteA].DependsOn).To(ConsistOf(builder.NATGatewayKey, builder.GAKey))
		Expect(vpcTemplate.Resources[privateRouteA].Properties).To(Equal(fakes.Properties{
			DestinationCidrBlock: builder.InternetCIDR,
			NatGatewayID:         map[string]interface{}{"Ref": builder.NATGatewayKey},
			RouteTableID:         map[string]interface{}{"Ref": privateRouteTableA},
		}))

		privateRouteB := builder.PrivateSubnetRouteKey + azBFormatted
		Expect(vpcTemplate.Resources).To(HaveKey(privateRouteB))
		Expect(vpcTemplate.Resources[privateRouteB].Type).To(Equal("AWS::EC2::Route"))
		Expect(vpcTemplate.Resources[privateRouteB].DependsOn).To(ConsistOf(builder.NATGatewayKey, builder.GAKey))
		Expect(vpcTemplate.Resources[privateRouteB].Properties).To(Equal(fakes.Properties{
			DestinationCidrBlock: builder.InternetCIDR,
			NatGatewayID:         map[string]interface{}{"Ref": builder.NATGatewayKey},
			RouteTableID:         map[string]interface{}{"Ref": privateRouteTableB},
		}))

		By("creating a ipv6 route to the ingress only internet gateway for each private subnet in the AZs")
		privateRouteA = builder.PrivateSubnetIpv6RouteKey + azAFormatted
		Expect(vpcTemplate.Resources).To(HaveKey(privateRouteA))
		Expect(vpcTemplate.Resources[privateRouteA].Type).To(Equal("AWS::EC2::Route"))
		Expect(vpcTemplate.Resources[privateRouteA].Properties).To(Equal(fakes.Properties{
			DestinationIpv6CidrBlock:    builder.InternetIPv6CIDR,
			EgressOnlyInternetGatewayID: map[string]interface{}{"Ref": builder.EgressOnlyInternetGatewayKey},
			RouteTableID:                map[string]interface{}{"Ref": privateRouteTableA},
		}))
		privateRouteB = builder.PrivateSubnetIpv6RouteKey + azBFormatted
		Expect(vpcTemplate.Resources).To(HaveKey(privateRouteB))
		Expect(vpcTemplate.Resources[privateRouteB].Type).To(Equal("AWS::EC2::Route"))
		Expect(vpcTemplate.Resources[privateRouteB].Properties).To(Equal(fakes.Properties{
			DestinationIpv6CidrBlock:    builder.InternetIPv6CIDR,
			EgressOnlyInternetGatewayID: map[string]interface{}{"Ref": builder.EgressOnlyInternetGatewayKey},
			RouteTableID:                map[string]interface{}{"Ref": privateRouteTableB},
		}))

		By("creating a public and private subnet for each AZ")
		assertSubnetSet := func(az, subnetKey, kubernetesTag string, cidrBlockIndex float64, mapPublicIpOnLaunch bool) {
			Expect(vpcTemplate.Resources).To(HaveKey(subnetKey))
			Expect(vpcTemplate.Resources[subnetKey].Type).To(Equal("AWS::EC2::Subnet"))
			Expect(vpcTemplate.Resources[subnetKey].DependsOn).To(ConsistOf(builder.IPv6CIDRBlockKey))
			Expect(vpcTemplate.Resources[subnetKey].Properties.AvailabilityZone).To(Equal(az))
			Expect(vpcTemplate.Resources[subnetKey].Properties.MapPublicIPOnLaunch).To(Equal(mapPublicIpOnLaunch))

			Expect(vpcTemplate.Resources[subnetKey].Properties.VpcID).To(Equal(map[string]interface{}{"Ref": "VPC"}))
			Expect(vpcTemplate.Resources[subnetKey].Properties.Tags).To(ConsistOf(
				fakes.Tag{
					Key:   kubernetesTag,
					Value: "1",
				},
				fakes.Tag{
					Key:   "Name",
					Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", subnetKey)},
				},
			))

			expectedFnIPv4CIDR := `{ "Fn::Cidr": [{ "Fn::GetAtt": ["VPC", "CidrBlock"]}, 6, 13 ]}`
			assertCidrBlockCreatedWithSelect(vpcTemplate.Resources[subnetKey].Properties.CidrBlock, expectedFnIPv4CIDR, cidrBlockIndex)

			expectedFnIPv6CIDR := `{ "Fn::Cidr": [{ "Fn::Select": [ 0, { "Fn::GetAtt": ["VPC", "Ipv6CidrBlocks"] }]}, 6, 64 ]}`
			assertCidrBlockCreatedWithSelect(vpcTemplate.Resources[subnetKey].Properties.Ipv6CidrBlock, expectedFnIPv6CIDR, cidrBlockIndex)
		}
		assertSubnetSet(azA, builder.PublicSubnetKey+azAFormatted, "kubernetes.io/role/elb", float64(0), true)
		Expect(vpcTemplate.Resources[builder.PublicSubnetKey+azAFormatted].Properties.AssignIpv6AddressOnCreation).To(BeNil())
		assertSubnetSet(azB, builder.PublicSubnetKey+azBFormatted, "kubernetes.io/role/elb", float64(1), true)
		Expect(vpcTemplate.Resources[builder.PublicSubnetKey+azBFormatted].Properties.AssignIpv6AddressOnCreation).To(BeNil())

		assertSubnetSet(azA, builder.PrivateSubnetKey+azAFormatted, "kubernetes.io/role/internal-elb", float64(2), false)
		Expect(*vpcTemplate.Resources[builder.PrivateSubnetKey+azAFormatted].Properties.AssignIpv6AddressOnCreation).To(Equal(true))
		assertSubnetSet(azB, builder.PrivateSubnetKey+azBFormatted, "kubernetes.io/role/internal-elb", float64(3), false)
		Expect(*vpcTemplate.Resources[builder.PrivateSubnetKey+azAFormatted].Properties.AssignIpv6AddressOnCreation).To(Equal(true))

		By("creating route table associations", func() {
			assertSubnetRouteTableAssociation := func(routeTableAssociationKey, subnetKey, routeTableKey string) {
				Expect(vpcTemplate.Resources).To(HaveKey(routeTableAssociationKey))
				Expect(vpcTemplate.Resources[routeTableAssociationKey].Type).To(Equal("AWS::EC2::SubnetRouteTableAssociation"))
				Expect(vpcTemplate.Resources[routeTableAssociationKey].Properties).To(Equal(fakes.Properties{
					RouteTableID: map[string]interface{}{"Ref": routeTableKey},
					SubnetID:     map[string]interface{}{"Ref": subnetKey},
				}))
			}

			By("associating all public subnets with the public route table", func() {
				assertSubnetRouteTableAssociation(builder.PubRouteTableAssociation+azAFormatted, builder.PublicSubnetKey+azAFormatted, builder.PubRouteTableKey)
				assertSubnetRouteTableAssociation(builder.PubRouteTableAssociation+azBFormatted, builder.PublicSubnetKey+azBFormatted, builder.PubRouteTableKey)
			})

			By("associating each private subnet with its private route table", func() {
				assertSubnetRouteTableAssociation(builder.PrivateRouteTableAssociation+azAFormatted, builder.PrivateSubnetKey+azAFormatted, builder.PrivateRouteTableKey+azAFormatted)
				assertSubnetRouteTableAssociation(builder.PrivateRouteTableAssociation+azBFormatted, builder.PrivateSubnetKey+azBFormatted, builder.PrivateRouteTableKey+azBFormatted)
			})
		})

		By("outputting the VPC on the stack")
		Expect(vpcTemplate.Outputs).To(HaveKey(builder.VPCResourceKey))
		Expect(vpcTemplate.Outputs.(map[string]interface{})[builder.VPCResourceKey].(map[string]interface{})["Value"]).To(Equal(map[string]interface{}{"Ref": builder.VPCResourceKey}))
		Expect(vpcTemplate.Outputs.(map[string]interface{})[builder.VPCResourceKey].(map[string]interface{})["Export"]).To(Equal(map[string]interface{}{
			"Name": map[string]interface{}{
				"Fn::Sub": fmt.Sprintf("${AWS::StackName}::%s", builder.VPCResourceKey),
			},
		}))

		By("outputting the public subnets on the stack")
		Expect(vpcTemplate.Outputs).To(HaveKey(outputs.ClusterSubnetsPublic))
		Expect(vpcTemplate.Outputs.(map[string]interface{})[outputs.ClusterSubnetsPublic].(map[string]interface{})["Value"]).To(Equal(map[string]interface{}{
			"Fn::Join": []interface{}{
				",",
				[]interface{}{
					map[string]interface{}{"Ref": builder.PublicSubnetKey + azAFormatted},
					map[string]interface{}{"Ref": builder.PublicSubnetKey + azBFormatted},
				},
			},
		}))
		Expect(vpcTemplate.Outputs.(map[string]interface{})[outputs.ClusterSubnetsPublic].(map[string]interface{})["Export"]).To(Equal(map[string]interface{}{
			"Name": map[string]interface{}{
				"Fn::Sub": fmt.Sprintf("${AWS::StackName}::%s", outputs.ClusterSubnetsPublic),
			},
		}))

		By("outputting the private subnets on the stack")
		Expect(vpcTemplate.Outputs).To(HaveKey(outputs.ClusterSubnetsPrivate))
		Expect(vpcTemplate.Outputs.(map[string]interface{})[outputs.ClusterSubnetsPrivate].(map[string]interface{})["Value"]).To(Equal(map[string]interface{}{
			"Fn::Join": []interface{}{
				",",
				[]interface{}{
					map[string]interface{}{"Ref": builder.PrivateSubnetKey + azAFormatted},
					map[string]interface{}{"Ref": builder.PrivateSubnetKey + azBFormatted},
				},
			},
		}))
		Expect(vpcTemplate.Outputs.(map[string]interface{})[outputs.ClusterSubnetsPrivate].(map[string]interface{})["Export"]).To(Equal(map[string]interface{}{
			"Name": map[string]interface{}{
				"Fn::Sub": fmt.Sprintf("${AWS::StackName}::%s", outputs.ClusterSubnetsPrivate),
			},
		}))
	})
	When("custom cidr block is provided", func() {
		var (
			cfg *api.ClusterConfig
		)
		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			cfg.KubernetesNetworkConfig.IPFamily = api.IPV6Family
			cfg.AvailabilityZones = []string{azA, azB}
			cfg.VPC.Network.CIDR = ipnet.MustParseCIDR("10.1.0.0/20")
		})
		It("calculates the correct cidr blocks for the subnets", func() {
			vpcRs := builder.NewIPv6VPCResourceSet(builder.NewRS(), cfg, nil)
			_, subnetDetails, err := vpcRs.CreateTemplate(context.Background())
			Expect(err).NotTo(HaveOccurred())

			By("returning the references of public subnets")
			pubRefs := subnetDetails.PublicSubnetRefs()
			Expect(pubRefs).To(HaveLen(2))
			Expect(pubRefs).To(ContainElement(makePrimitive(builder.PublicSubnetKey + azAFormatted)))
			Expect(pubRefs).To(ContainElement(makePrimitive(builder.PublicSubnetKey + azBFormatted)))

			By("returning the references of private subnets")
			privRef := subnetDetails.PrivateSubnetRefs()
			Expect(privRef).To(HaveLen(2))
			Expect(privRef).To(ContainElement(makePrimitive(builder.PrivateSubnetKey + azBFormatted)))
			Expect(privRef).To(ContainElement(makePrimitive(builder.PrivateSubnetKey + azBFormatted)))

			vpcTemplate, err := renderTemplate(vpcRs)
			Expect(err).NotTo(HaveOccurred())

			assertSubnetSet := func(az, subnetKey, kubernetesTag string, cidrBlockIndex float64, mapPublicIpOnLaunch bool) {
				Expect(vpcTemplate.Resources).To(HaveKey(subnetKey))
				Expect(vpcTemplate.Resources[subnetKey].Type).To(Equal("AWS::EC2::Subnet"))
				Expect(vpcTemplate.Resources[subnetKey].DependsOn).To(ConsistOf(builder.IPv6CIDRBlockKey))
				Expect(vpcTemplate.Resources[subnetKey].Properties.AvailabilityZone).To(Equal(az))
				Expect(vpcTemplate.Resources[subnetKey].Properties.MapPublicIPOnLaunch).To(Equal(mapPublicIpOnLaunch))

				Expect(vpcTemplate.Resources[subnetKey].Properties.VpcID).To(Equal(map[string]interface{}{"Ref": "VPC"}))
				Expect(vpcTemplate.Resources[subnetKey].Properties.Tags).To(ConsistOf(
					fakes.Tag{
						Key:   kubernetesTag,
						Value: "1",
					},
					fakes.Tag{
						Key:   "Name",
						Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", subnetKey)},
					},
				))

				// Note, this is the important difference.
				expectedFnIPv4CIDR := `{ "Fn::Cidr": [{ "Fn::GetAtt": ["VPC", "CidrBlock"]}, 6, 9 ]}`
				assertCidrBlockCreatedWithSelect(vpcTemplate.Resources[subnetKey].Properties.CidrBlock, expectedFnIPv4CIDR, cidrBlockIndex)

				expectedFnIPv6CIDR := `{ "Fn::Cidr": [{ "Fn::Select": [ 0, { "Fn::GetAtt": ["VPC", "Ipv6CidrBlocks"] }]}, 6, 64 ]}`
				assertCidrBlockCreatedWithSelect(vpcTemplate.Resources[subnetKey].Properties.Ipv6CidrBlock, expectedFnIPv6CIDR, cidrBlockIndex)
			}
			assertSubnetSet(azA, builder.PublicSubnetKey+azAFormatted, "kubernetes.io/role/elb", float64(0), true)
			Expect(vpcTemplate.Resources[builder.PublicSubnetKey+azAFormatted].Properties.AssignIpv6AddressOnCreation).To(BeNil())
			assertSubnetSet(azB, builder.PublicSubnetKey+azBFormatted, "kubernetes.io/role/elb", float64(1), true)
			Expect(vpcTemplate.Resources[builder.PublicSubnetKey+azBFormatted].Properties.AssignIpv6AddressOnCreation).To(BeNil())

			assertSubnetSet(azA, builder.PrivateSubnetKey+azAFormatted, "kubernetes.io/role/internal-elb", float64(2), false)
			Expect(*vpcTemplate.Resources[builder.PrivateSubnetKey+azAFormatted].Properties.AssignIpv6AddressOnCreation).To(Equal(true))
			assertSubnetSet(azB, builder.PrivateSubnetKey+azBFormatted, "kubernetes.io/role/internal-elb", float64(3), false)
			Expect(*vpcTemplate.Resources[builder.PrivateSubnetKey+azAFormatted].Properties.AssignIpv6AddressOnCreation).To(Equal(true))
		})
	})

	When("private cluster is enabled", func() {
		It("creates only private IPv6 resources", func() {
			cfg := cfg.DeepCopy()
			cfg.PrivateCluster = &api.PrivateCluster{
				Enabled: true,
			}
			cfg.AvailabilityZones = []string{azA, azB}
			vpcRs := builder.NewIPv6VPCResourceSet(builder.NewRS(), cfg, nil)

			_, subnetDetails, err := vpcRs.CreateTemplate(context.Background())
			Expect(err).NotTo(HaveOccurred())

			By("returning the references of public subnets")
			pubRefs := subnetDetails.PublicSubnetRefs()
			Expect(pubRefs).To(BeEmpty())

			By("returning the references of private subnets")
			privRef := subnetDetails.PrivateSubnetRefs()
			Expect(privRef).To(HaveLen(2))
			Expect(privRef).To(ContainElement(makePrimitive(builder.PrivateSubnetKey + azBFormatted)))
			Expect(privRef).To(ContainElement(makePrimitive(builder.PrivateSubnetKey + azBFormatted)))

			vpcTemplate := &fakes.FakeTemplate{}
			templateBody, err := vpcRs.RenderJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(json.Unmarshal(templateBody, vpcTemplate)).To(Succeed())

			By("creating the VPC resource")
			Expect(vpcTemplate.Resources).To(HaveKey(builder.VPCResourceKey))
			Expect(vpcTemplate.Resources[builder.VPCResourceKey].Type).To(Equal("AWS::EC2::VPC"))
			defaultCidr := api.DefaultCIDR()
			cidr := &defaultCidr
			Expect(vpcTemplate.Resources[builder.VPCResourceKey].Properties).To(Equal(fakes.Properties{
				CidrBlock:          cidr.String(),
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				Tags: []fakes.Tag{
					{
						Key:   "Name",
						Value: map[string]interface{}{"Fn::Sub": "${AWS::StackName}/VPC"},
					},
				},
			}))

			By("creating the IPv6 CIDR")
			Expect(vpcTemplate.Resources).To(HaveKey(builder.IPv6CIDRBlockKey))
			Expect(vpcTemplate.Resources[builder.IPv6CIDRBlockKey].Type).To(Equal("AWS::EC2::VPCCidrBlock"))
			Expect(vpcTemplate.Resources[builder.IPv6CIDRBlockKey].Properties).To(Equal(fakes.Properties{
				AmazonProvidedIpv6CidrBlock: true,
				VpcID:                       map[string]interface{}{"Ref": "VPC"},
			}))

			By("creating the internet gateway")
			Expect(vpcTemplate.Resources).ToNot(HaveKey(builder.IGWKey))

			By("creating a VPC gateway attachment to associate the IGW with the VPC")
			Expect(vpcTemplate.Resources).ToNot(HaveKey(builder.GAKey))

			By("creating a VPC gateway attachment to associate the IGW with the VPC")
			Expect(vpcTemplate.Resources).ToNot(HaveKey(builder.EgressOnlyInternetGatewayKey))

			By("creating the NAT gateway")
			Expect(vpcTemplate.Resources).ToNot(HaveKey(builder.NATGatewayKey))

			By("creating an Elastic IP for the Nat Gateway")
			Expect(vpcTemplate.Resources).ToNot(HaveKey(builder.ElasticIPKey))

			By("creating a public Route Table")
			Expect(vpcTemplate.Resources).ToNot(HaveKey(builder.PubRouteTableKey))

			By("creating public subnet route for IPv4 traffic to IPv4 CIDR")
			Expect(vpcTemplate.Resources).ToNot(HaveKey(builder.PubSubRouteKey))

			By("creating public subnet route for IPv6 traffic to IPv6 CIDR")
			Expect(vpcTemplate.Resources).ToNot(HaveKey(builder.PubSubIPv6RouteKey))

			By("creating a private route table for each AZ")
			privateRouteTableA := builder.PrivateRouteTableKey + azAFormatted
			Expect(vpcTemplate.Resources).To(HaveKey(privateRouteTableA))
			Expect(vpcTemplate.Resources[privateRouteTableA].Type).To(Equal("AWS::EC2::RouteTable"))
			Expect(vpcTemplate.Resources[privateRouteTableA].Properties).To(Equal(fakes.Properties{
				VpcID: map[string]interface{}{"Ref": "VPC"},
				Tags: []fakes.Tag{
					{
						Key:   "Name",
						Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", privateRouteTableA)},
					},
				},
			}))
			privateRouteTableB := builder.PrivateRouteTableKey + azBFormatted
			Expect(vpcTemplate.Resources).To(HaveKey(privateRouteTableB))
			Expect(vpcTemplate.Resources[privateRouteTableB].Type).To(Equal("AWS::EC2::RouteTable"))
			Expect(vpcTemplate.Resources[privateRouteTableB].Properties).To(Equal(fakes.Properties{
				VpcID: map[string]interface{}{"Ref": "VPC"},
				Tags: []fakes.Tag{
					{
						Key:   "Name",
						Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", privateRouteTableB)},
					},
				},
			}))

			By("creating a route to the NAT gateway for each private subnet in the AZs")
			privateRouteA := builder.PrivateSubnetRouteKey + azAFormatted
			Expect(vpcTemplate.Resources).NotTo(HaveKey(privateRouteA))

			privateRouteB := builder.PrivateSubnetRouteKey + azBFormatted
			Expect(vpcTemplate.Resources).NotTo(HaveKey(privateRouteB))

			By("creating a ipv6 route to the ingress only internet gateway for each private subnet in the AZs")
			privateRouteA = builder.PrivateSubnetIpv6RouteKey + azAFormatted
			Expect(vpcTemplate.Resources).NotTo(HaveKey(privateRouteA))
			Expect(vpcTemplate.Resources).NotTo(HaveKey(privateRouteB))

			By("creating a private subnet for each AZ")
			assertSubnetSet := func(az, subnetKey, kubernetesTag string, cidrBlockIndex float64, mapPublicIpOnLaunch bool) {
				Expect(vpcTemplate.Resources).To(HaveKey(subnetKey))
				Expect(vpcTemplate.Resources[subnetKey].Type).To(Equal("AWS::EC2::Subnet"))
				Expect(vpcTemplate.Resources[subnetKey].DependsOn).To(ConsistOf(builder.IPv6CIDRBlockKey))
				Expect(vpcTemplate.Resources[subnetKey].Properties.AvailabilityZone).To(Equal(az))
				Expect(vpcTemplate.Resources[subnetKey].Properties.MapPublicIPOnLaunch).To(Equal(mapPublicIpOnLaunch))

				Expect(vpcTemplate.Resources[subnetKey].Properties.VpcID).To(Equal(map[string]interface{}{"Ref": "VPC"}))
				Expect(vpcTemplate.Resources[subnetKey].Properties.Tags).To(ConsistOf(
					fakes.Tag{
						Key:   kubernetesTag,
						Value: "1",
					},
					fakes.Tag{
						Key:   "Name",
						Value: map[string]interface{}{"Fn::Sub": fmt.Sprintf("${AWS::StackName}/%s", subnetKey)},
					},
				))

				expectedFnIPv4CIDR := `{ "Fn::Cidr": [{ "Fn::GetAtt": ["VPC", "CidrBlock"]}, 6, 13 ]}`
				Expect(vpcTemplate.Resources[subnetKey].Properties.CidrBlock.(map[string]interface{})["Fn::Select"]).To(HaveLen(2))
				Expect(vpcTemplate.Resources[subnetKey].Properties.CidrBlock.(map[string]interface{})["Fn::Select"].([]interface{})[0].(float64)).To(Equal(cidrBlockIndex))
				actualFnCIDR, err := json.Marshal(vpcTemplate.Resources[subnetKey].Properties.CidrBlock.(map[string]interface{})["Fn::Select"].([]interface{})[1])
				Expect(err).NotTo(HaveOccurred())
				Expect(actualFnCIDR).To(MatchJSON([]byte(expectedFnIPv4CIDR)))

				expectedFnIPv6CIDR := `{ "Fn::Cidr": [{ "Fn::Select": [ 0, { "Fn::GetAtt": ["VPC", "Ipv6CidrBlocks"] }]}, 6, 64 ]}`
				assertCidrBlockCreatedWithSelect(vpcTemplate.Resources[subnetKey].Properties.Ipv6CidrBlock, expectedFnIPv6CIDR, cidrBlockIndex)
			}
			assertSubnetSet(azA, builder.PrivateSubnetKey+azAFormatted, "kubernetes.io/role/internal-elb", float64(2), false)
			Expect(*vpcTemplate.Resources[builder.PrivateSubnetKey+azAFormatted].Properties.AssignIpv6AddressOnCreation).To(Equal(true))
			assertSubnetSet(azB, builder.PrivateSubnetKey+azBFormatted, "kubernetes.io/role/internal-elb", float64(3), false)
			Expect(*vpcTemplate.Resources[builder.PrivateSubnetKey+azAFormatted].Properties.AssignIpv6AddressOnCreation).To(Equal(true))

			By("creating route table associations", func() {
				assertSubnetRouteTableAssociation := func(routeTableAssociationKey, subnetKey, routeTableKey string) {
					Expect(vpcTemplate.Resources).To(HaveKey(routeTableAssociationKey))
					Expect(vpcTemplate.Resources[routeTableAssociationKey].Type).To(Equal("AWS::EC2::SubnetRouteTableAssociation"))
					Expect(vpcTemplate.Resources[routeTableAssociationKey].Properties).To(Equal(fakes.Properties{
						RouteTableID: map[string]interface{}{"Ref": routeTableKey},
						SubnetID:     map[string]interface{}{"Ref": subnetKey},
					}))
				}

				By("associating each private subnet with its private route table", func() {
					assertSubnetRouteTableAssociation(builder.PrivateRouteTableAssociation+azAFormatted, builder.PrivateSubnetKey+azAFormatted, builder.PrivateRouteTableKey+azAFormatted)
					assertSubnetRouteTableAssociation(builder.PrivateRouteTableAssociation+azBFormatted, builder.PrivateSubnetKey+azBFormatted, builder.PrivateRouteTableKey+azBFormatted)
				})
			})

			By("outputting the VPC on the stack")
			Expect(vpcTemplate.Outputs).To(HaveKey(builder.VPCResourceKey))
			Expect(vpcTemplate.Outputs.(map[string]interface{})[builder.VPCResourceKey].(map[string]interface{})["Value"]).To(Equal(map[string]interface{}{"Ref": builder.VPCResourceKey}))
			Expect(vpcTemplate.Outputs.(map[string]interface{})[builder.VPCResourceKey].(map[string]interface{})["Export"]).To(Equal(map[string]interface{}{
				"Name": map[string]interface{}{
					"Fn::Sub": fmt.Sprintf("${AWS::StackName}::%s", builder.VPCResourceKey),
				},
			}))

			By("outputting the public subnets on the stack")
			Expect(vpcTemplate.Outputs).ToNot(HaveKey(outputs.ClusterSubnetsPublic))

			By("outputting the private subnets on the stack")
			Expect(vpcTemplate.Outputs).To(HaveKey(outputs.ClusterSubnetsPrivate))
			Expect(vpcTemplate.Outputs.(map[string]interface{})[outputs.ClusterSubnetsPrivate].(map[string]interface{})["Value"]).To(Equal(map[string]interface{}{
				"Fn::Join": []interface{}{
					",",
					[]interface{}{
						map[string]interface{}{"Ref": builder.PrivateSubnetKey + azAFormatted},
						map[string]interface{}{"Ref": builder.PrivateSubnetKey + azBFormatted},
					},
				},
			}))
			Expect(vpcTemplate.Outputs.(map[string]interface{})[outputs.ClusterSubnetsPrivate].(map[string]interface{})["Export"]).To(Equal(map[string]interface{}{
				"Name": map[string]interface{}{
					"Fn::Sub": fmt.Sprintf("${AWS::StackName}::%s", outputs.ClusterSubnetsPrivate),
				},
			}))
		})
	})

	Context("when there are 3 AZs", func() {
		BeforeEach(func() {
			cfg.AvailabilityZones = []string{azA, azB, azC}
		})

		It("scales the CIDR blocks accordingly", func() {
			vpcRs := builder.NewIPv6VPCResourceSet(builder.NewRS(), cfg, nil)
			vpcTemplate, err := createAndRenderTemplate(vpcRs)
			Expect(err).NotTo(HaveOccurred())

			assertSubnetSet := func(az, subnetKey string, cidrBlockIndex float64) {
				Expect(vpcTemplate.Resources).To(HaveKey(subnetKey))
				expectedFnIPv4CIDR := `{ "Fn::Cidr": [{ "Fn::GetAtt": ["VPC", "CidrBlock"]}, 8, 13 ]}`
				assertCidrBlockCreatedWithSelect(vpcTemplate.Resources[subnetKey].Properties.CidrBlock, expectedFnIPv4CIDR, cidrBlockIndex)

				expectedFnIPv6CIDR := `{ "Fn::Cidr": [{ "Fn::Select": [ 0, { "Fn::GetAtt": ["VPC", "Ipv6CidrBlocks"] }]}, 8, 64 ]}`
				assertCidrBlockCreatedWithSelect(vpcTemplate.Resources[subnetKey].Properties.Ipv6CidrBlock, expectedFnIPv6CIDR, cidrBlockIndex)
			}
			assertSubnetSet(azA, builder.PublicSubnetKey+azAFormatted, float64(0))
			assertSubnetSet(azB, builder.PublicSubnetKey+azBFormatted, float64(1))
			assertSubnetSet(azC, builder.PublicSubnetKey+azCFormatted, float64(2))

			assertSubnetSet(azA, builder.PrivateSubnetKey+azAFormatted, float64(3))
			assertSubnetSet(azB, builder.PrivateSubnetKey+azBFormatted, float64(4))
			assertSubnetSet(azC, builder.PrivateSubnetKey+azCFormatted, float64(5))

		})
	})

	When("a user provides a custom ipv6 block", func() {
		BeforeEach(func() {
			cfg.VPC.IPv6Cidr = "my-cidr"
			cfg.VPC.IPv6Pool = "my-cidr-pool"
		})

		It("creates the IPv6CidrBlock resource with the users ipv6 pool", func() {
			vpcRs := builder.NewIPv6VPCResourceSet(builder.NewRS(), cfg, nil)
			vpcTemplate, err := createAndRenderTemplate(vpcRs)
			Expect(err).NotTo(HaveOccurred())

			By("creating the IPv6 CIDR")
			Expect(vpcTemplate.Resources).To(HaveKey(builder.IPv6CIDRBlockKey))
			Expect(vpcTemplate.Resources[builder.IPv6CIDRBlockKey].Type).To(Equal("AWS::EC2::VPCCidrBlock"))
			Expect(vpcTemplate.Resources[builder.IPv6CIDRBlockKey].Properties).To(Equal(fakes.Properties{
				Ipv6CidrBlock: "my-cidr",
				Ipv6Pool:      "my-cidr-pool",
				VpcID:         map[string]interface{}{"Ref": "VPC"},
			}))
		})
	})

	When("a user provides a custom ipv4 cidr", func() {
		var customCidr = &ipnet.IPNet{
			IPNet: net.IPNet{
				IP:   []byte{192, 168, 1, 1},
				Mask: []byte{255, 255, 0, 0},
			},
		}

		BeforeEach(func() {
			cfg.VPC.CIDR = customCidr
		})

		It("creates the VPC resource with the users provided ipv4 cidr", func() {
			vpcRs := builder.NewIPv6VPCResourceSet(builder.NewRS(), cfg, nil)
			vpcTemplate, err := createAndRenderTemplate(vpcRs)
			Expect(err).NotTo(HaveOccurred())

			Expect(vpcTemplate.Resources).To(HaveKey(builder.VPCResourceKey))
			Expect(vpcTemplate.Resources[builder.VPCResourceKey].Type).To(Equal("AWS::EC2::VPC"))
			Expect(vpcTemplate.Resources[builder.VPCResourceKey].Properties).To(Equal(fakes.Properties{
				CidrBlock:          customCidr.String(),
				EnableDNSHostnames: true,
				EnableDNSSupport:   true,
				Tags: []fakes.Tag{
					{
						Key:   "Name",
						Value: map[string]interface{}{"Fn::Sub": "${AWS::StackName}/VPC"},
					},
				},
			}))
		})
	})
})

func createAndRenderTemplate(vpcRs *builder.IPv6VPCResourceSet) (*fakes.FakeTemplate, error) {
	_, _, err := vpcRs.CreateTemplate(context.Background())
	if err != nil {
		return nil, err
	}
	return renderTemplate(vpcRs)
}

func renderTemplate(vpcRs *builder.IPv6VPCResourceSet) (*fakes.FakeTemplate, error) {
	vpcTemplate := &fakes.FakeTemplate{}
	templateBody, err := vpcRs.RenderJSON()
	if err != nil {
		return nil, err
	}
	ExpectWithOffset(1, json.Unmarshal(templateBody, vpcTemplate)).To(Succeed())
	return vpcTemplate, nil
}

func assertCidrBlockCreatedWithSelect(cidrBlock interface{}, expectedFnCIDR string, cidrBlockIndex float64) {
	ExpectWithOffset(1, cidrBlock.(map[string]interface{})).To(HaveKey("Fn::Select"))
	fnSelectValue := cidrBlock.(map[string]interface{})["Fn::Select"].([]interface{})
	ExpectWithOffset(1, fnSelectValue).To(HaveLen(2))
	ExpectWithOffset(1, fnSelectValue[0].(float64)).To(Equal(cidrBlockIndex))
	actualFnCIDR, err := json.Marshal(fnSelectValue[1])
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ExpectWithOffset(1, actualFnCIDR).To(MatchJSON([]byte(expectedFnCIDR)))
}
