package builder

import (
	"strings"

	gfncfn "github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

const (
	VPCResourceKey = "VPC"

	// Gateways
	IGWKey                       = "InternetGateway"
	GAKey                        = "VPCGatewayAttachment"
	EgressOnlyInternetGatewayKey = "EgressOnlyInternetGateway"
	NATGatewayKey                = "NATGateway"
	ElasticIPKey                 = "EIP"

	// CIDRs
	IPv6CIDRBlockKey = "IPv6CidrBlock"
	InternetCIDR     = "0.0.0.0/0"
	InternetIPv6CIDR = "::/0"

	// Routing
	PubRouteTableKey             = "PublicRouteTable"
	PrivateRouteTableKey         = "PrivateRouteTable"
	PubRouteTableAssociation     = "RouteTableAssociationPublic"
	PrivateRouteTableAssociation = "RouteTableAssociationPrivate"
	PubSubRouteKey               = "PublicSubnetDefaultRoute"
	PubSubIPv6RouteKey           = "PublicSubnetIPv6DefaultRoute"
	PrivateSubnetRouteKey        = "PrivateSubnetDefaultRoute"
	PrivateSubnetIpv6RouteKey    = "PrivateSubnetDefaultIpv6Route"

	// Subnets
	PublicSubnetKey         = "PublicSubnet"
	PrivateSubnetKey        = "PrivateSubnet"
	PublicSubnetsOutputKey  = "SubnetsPublic"
	PrivateSubnetsOutputKey = "SubnetsPrivate"
)

//VPCResourceSet interface for creating cloudformation resource sets for generating VPC resources
type VPCResourceSet interface {
	//CreateTemplate generates all of the resources & outputs required for the VPC. Returns the
	CreateTemplate() (vpcID *gfnt.Value, subnetDetails *SubnetDetails, err error)
}

func formatAZ(az string) string {
	return strings.ToUpper(strings.ReplaceAll(az, "-", ""))
}

func getSubnetIPv6CIDRBlock(cidrPartitions int) *gfnt.Value {
	// get 8 of /64 subnets from the auto-allocated IPv6 block,
	// and pick one block based on subnetIndexForIPv6 counter;
	// NOTE: this is done inside of CloudFormation using Fn::Cidr,
	// we don't slice it here, just construct the JSON expression
	// that does slicing at runtime.
	refIPv6CIDRv6 := gfnt.MakeFnSelect(
		gfnt.NewInteger(0), gfnt.MakeFnGetAttString("VPC", "Ipv6CidrBlocks"),
	)
	refSubnetSlices := gfnt.MakeFnCIDR(refIPv6CIDRv6, gfnt.NewInteger(cidrPartitions), gfnt.NewInteger(64))
	return refSubnetSlices
}

func getSubnetIPv4CIDRBlock(cidrPartitions int) *gfnt.Value {
	//TODO: should we be doing /19? Should we adjust for the partition size?
	desiredMask := 19
	refSubnetSlices := gfnt.MakeFnCIDR(gfnt.MakeFnGetAttString("VPC", "CidrBlock"), gfnt.NewInteger(cidrPartitions), gfnt.NewInteger(32-desiredMask))
	return refSubnetSlices
}

func (rs *resourceSet) addEFASecurityGroup(vpcID *gfnt.Value, clusterName, desc string) *gfnt.Value {
	efaSG := rs.newResource("EFASG", &gfnec2.SecurityGroup{
		VpcId:            vpcID,
		GroupDescription: gfnt.NewString("EFA-enabled security group"),
		Tags: []gfncfn.Tag{{
			Key:   gfnt.NewString("kubernetes.io/cluster/" + clusterName),
			Value: gfnt.NewString("owned"),
		}},
	})
	rs.newResource("EFAIngressSelf", &gfnec2.SecurityGroupIngress{
		GroupId:               efaSG,
		SourceSecurityGroupId: efaSG,
		Description:           gfnt.NewString("Allow " + desc + " to communicate to itself (EFA-enabled)"),
		IpProtocol:            gfnt.NewString("-1"),
	})
	rs.newResource("EFAEgressSelf", &gfnec2.SecurityGroupEgress{
		GroupId:                    efaSG,
		DestinationSecurityGroupId: efaSG,
		Description:                gfnt.NewString("Allow " + desc + " to communicate to itself (EFA-enabled)"),
		IpProtocol:                 gfnt.NewString("-1"),
	})

	return efaSG
}
