package builder

import (
	"strings"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

const (
	VPCResourceKey, IGWKey, GAKey                          = "VPC", "InternetGateway", "VPCGatewayAttachment"
	IPv6CIDRBlockKey                                       = "IPv6CidrBlock"
	EgressOnlyInternetGatewayKey                           = "EgressOnlyInternetGateway"
	ElasticIPKey                                           = "EIP"
	InternetCIDR, InternetIPv6CIDR                         = "0.0.0.0/0", "::/0"
	PubRouteTableKey, PrivateRouteTableKey                 = "PublicRouteTable", "PrivateRouteTable"
	PubRouteTableAssociation, PrivateRouteTableAssociation = "RouteTableAssociationPublic", "RouteTableAssociationPrivate"
	PubSubRouteKey, PubSubIPv6RouteKey                     = "PublicSubnetDefaultRoute", "PublicSubnetIPv6DefaultRoute"
	PrivateSubnetRouteKey, PrivateSubnetIpv6RouteKey       = "PrivateSubnetDefaultRoute", "PrivateSubnetDefaultIpv6Route"
	PublicSubnetKey, PrivateSubnetKey                      = "PublicSubnet", "PrivateSubnet"
	NATGatewayKey                                          = "NATGateway"
	PublicSubnetsOutputKey, PrivateSubnetsOutputKey        = "SubnetsPublic", "SubnetsPrivate"
)

type VPCResourceSet interface {
	CreateTemplate() (*gfnt.Value, *SubnetDetails, error)
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
