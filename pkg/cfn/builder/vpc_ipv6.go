package builder

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
	"github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

// A IPv6VPCResourceSet builds the resources required for the specified VPC
type IPv6VPCResourceSet struct {
	rs            *resourceSet
	clusterConfig *api.ClusterConfig
	ec2API        ec2iface.EC2API

	// vpcResource *IPv6VPCResource
}

// // IPv6VPCResource reresents a VPC resource
// type IPv6VPCResource struct {
// 	VPC           *gfnt.Value
// 	SubnetDetails *subnetDetails
// }

// NewIPv6VPCResourceSet creates and returns a new VPCResourceSet
func NewIPv6VPCResourceSet(rs *resourceSet, clusterConfig *api.ClusterConfig, ec2API ec2iface.EC2API) *IPv6VPCResourceSet {
	return &IPv6VPCResourceSet{
		rs:            rs,
		clusterConfig: clusterConfig,
		ec2API:        ec2API,
		// vpcResource: &IPv6VPCResource{
		// 	VPC:           vpcRef,
		// 	SubnetDetails: &subnetDetails{},
		// },
	}
}

func (v *IPv6VPCResourceSet) CreateTemplate() (*gfnt.Value, *SubnetDetails, error) {
	var publicSubnetResourceRefs, privateSubnetResourceRefs []*gfnt.Value
	vpcResourceRef := v.rs.newResource(VPCResourceKey, &gfnec2.VPC{
		CidrBlock:          gfnt.NewString(v.clusterConfig.VPC.CIDR.String()),
		EnableDnsSupport:   gfnt.True(),
		EnableDnsHostnames: gfnt.True(),
	})

	v.rs.newResource(IPv6CIDRBlockKey, &gfnec2.VPCCidrBlock{
		AmazonProvidedIpv6CidrBlock: gfnt.True(),
		VpcId:                       gfnt.MakeRef(VPCResourceKey),
	})

	refIGW := v.rs.newResource(IGWKey, &gfnec2.InternetGateway{})

	v.rs.newResource(GAKey, &gfnec2.VPCGatewayAttachment{
		InternetGatewayId: gfnt.MakeRef(IGWKey),
		VpcId:             gfnt.MakeRef(VPCResourceKey),
	})

	v.rs.newResource(EgressOnlyInternetGatewayKey, &gfnec2.EgressOnlyInternetGateway{
		VpcId: gfnt.MakeRef(VPCResourceKey),
	})

	firstPublicSubnet := PublicSubnetKey + formatAZ(v.clusterConfig.AvailabilityZones[0])
	v.rs.newResource(NATGatewayKey, &gfnec2.NatGateway{
		AWSCloudFormationDependsOn: []string{
			ElasticIPKey,
			firstPublicSubnet,
			GAKey,
		},
		AllocationId: gfnt.MakeFnGetAtt(ElasticIPKey, gfnt.NewString("AllocationId")),
		SubnetId:     gfnt.MakeRef(firstPublicSubnet),
	})

	v.rs.newResource(ElasticIPKey, &gfnec2.EIP{
		Domain:                     gfnt.NewString("vpc"),
		AWSCloudFormationDependsOn: []string{GAKey},
	})

	v.rs.newResource(PubRouteTableKey, &gfnec2.RouteTable{
		VpcId: gfnt.MakeRef(VPCResourceKey),
	})

	v.rs.newResource(PubSubRouteKey, &gfnec2.Route{
		AWSCloudFormationDependsOn: []string{GAKey},
		DestinationCidrBlock:       gfnt.NewString(InternetCIDR),
		GatewayId:                  refIGW,
		RouteTableId:               gfnt.MakeRef(PubRouteTableKey),
	})

	v.rs.newResource(PubSubIPv6RouteKey, &gfnec2.Route{
		AWSCloudFormationDependsOn: []string{GAKey},
		DestinationIpv6CidrBlock:   gfnt.NewString(InternetIPv6CIDR),
		GatewayId:                  refIGW,
		RouteTableId:               gfnt.MakeRef(PubRouteTableKey),
	})

	cidrPartitions := (len(v.clusterConfig.AvailabilityZones) * 2) + 2
	for i, az := range v.clusterConfig.AvailabilityZones {
		azFormatted := formatAZ(az)
		v.rs.newResource(PrivateRouteTableKey+azFormatted, &gfnec2.RouteTable{
			VpcId: gfnt.MakeRef(VPCResourceKey),
		})

		v.rs.newResource(PrivateSubnetIpv6RouteKey+azFormatted, &gfnec2.Route{
			DestinationIpv6CidrBlock:    gfnt.NewString(InternetIPv6CIDR),
			EgressOnlyInternetGatewayId: gfnt.MakeRef(EgressOnlyInternetGatewayKey),
			RouteTableId:                gfnt.MakeRef(PrivateRouteTableKey + azFormatted),
		})

		v.rs.newResource(PrivateSubnetRouteKey+azFormatted, &gfnec2.Route{
			AWSCloudFormationDependsOn: []string{NATGatewayKey, GAKey},
			DestinationCidrBlock:       gfnt.NewString(InternetCIDR),
			NatGatewayId:               gfnt.MakeRef(NATGatewayKey),
			RouteTableId:               gfnt.MakeRef(PrivateRouteTableKey + azFormatted),
		})

		publicSubnetResourceRefs = append(publicSubnetResourceRefs, v.createSubnet(az, azFormatted, i, cidrPartitions, false))
		privateSubnetResourceRefs = append(privateSubnetResourceRefs, v.createSubnet(az, azFormatted, i+len(v.clusterConfig.AvailabilityZones), cidrPartitions, true))

		v.rs.newResource(PubRouteTableAssociation+azFormatted, &gfnec2.SubnetRouteTableAssociation{
			RouteTableId: gfnt.MakeRef(PubRouteTableKey),
			SubnetId:     gfnt.MakeRef(PublicSubnetKey + azFormatted),
		})

		v.rs.newResource(PrivateRouteTableAssociation+azFormatted, &gfnec2.SubnetRouteTableAssociation{
			RouteTableId: gfnt.MakeRef(PrivateRouteTableKey + azFormatted),
			SubnetId:     gfnt.MakeRef(PrivateSubnetKey + azFormatted),
		})
	}

	v.rs.defineOutput(outputs.ClusterVPC, vpcResourceRef, true, func(val string) error {
		v.clusterConfig.VPC.ID = val
		return nil
	})

	addSubnetOutput := func(subnetRefs []*gfnt.Value, topology api.SubnetTopology, outputName string) {
		v.rs.defineJoinedOutput(outputName, subnetRefs, true, func(value string) error {
			return vpc.ImportSubnetsFromIDList(v.ec2API, v.clusterConfig, topology, strings.Split(value, ","))
		})
	}

	addSubnetOutput(publicSubnetResourceRefs, api.SubnetTopologyPublic, outputs.ClusterSubnetsPublic)
	addSubnetOutput(privateSubnetResourceRefs, api.SubnetTopologyPrivate, outputs.ClusterSubnetsPrivate)

	var publicSubnets, privateSubnets []SubnetResource
	for _, s := range publicSubnetResourceRefs {
		publicSubnets = append(publicSubnets, SubnetResource{Subnet: s})
	}
	for _, s := range privateSubnetResourceRefs {
		privateSubnets = append(privateSubnets, SubnetResource{Subnet: s})
	}
	return vpcResourceRef, &SubnetDetails{
		Private: privateSubnets,
		Public:  publicSubnets,
	}, nil
}

func (v *IPv6VPCResourceSet) RenderJSON() ([]byte, error) {
	return v.rs.renderJSON()
}

func (v *IPv6VPCResourceSet) createSubnet(az, azFormatted string, i, cidrPartitions int, private bool) *gfnt.Value {
	var assignIpv6AddressOnCreation *gfnt.Value
	subnetKey := PublicSubnetKey + azFormatted
	mapPublicIPOnLaunch := gfnt.True()
	elbTagKey := "kubernetes.io/role/elb"

	if private {
		subnetKey = PrivateSubnetKey + azFormatted
		mapPublicIPOnLaunch = nil
		assignIpv6AddressOnCreation = gfnt.True()
		elbTagKey = "kubernetes.io/role/internal-elb"
	}

	return v.rs.newResource(subnetKey, &gfnec2.Subnet{
		AWSCloudFormationDependsOn:  []string{IPv6CIDRBlockKey},
		AvailabilityZone:            gfnt.NewString(az),
		CidrBlock:                   gfnt.MakeFnSelect(gfnt.NewInteger(i), getSubnetIPv4CIDRBlock(cidrPartitions)),
		Ipv6CidrBlock:               gfnt.MakeFnSelect(gfnt.NewInteger(i), getSubnetIPv6CIDRBlock(cidrPartitions)),
		MapPublicIpOnLaunch:         mapPublicIPOnLaunch,
		AssignIpv6AddressOnCreation: assignIpv6AddressOnCreation,
		VpcId:                       gfnt.MakeRef(VPCResourceKey),
		Tags: []cloudformation.Tag{{
			Key:   gfnt.NewString(elbTagKey),
			Value: gfnt.NewString("1"),
		}},
	})
}
