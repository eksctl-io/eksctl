package builder

import (
	"context"
	"math"
	"strings"

	gfncfn "goformation/v4/cloudformation/cloudformation"
	gfnec2 "goformation/v4/cloudformation/ec2"
	gfnt "goformation/v4/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
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
	PublicSubnetKey        = "PublicSubnet"
	PrivateSubnetKey       = "PrivateSubnet"
	defaultPrefix          = 19
	defaultSubnetMask      = 32
	defaultDesiredMaskSize = defaultSubnetMask - defaultPrefix
)

// VPCResourceSet interface for creating cloudformation resource sets for generating VPC resources
type VPCResourceSet interface {
	// CreateTemplate generates all of the resources & outputs required for the VPC. Returns the
	CreateTemplate(ctx context.Context) (vpcID *gfnt.Value, subnetDetails *SubnetDetails, err error)
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

func getSubnetIPv4CIDRBlock(cidrPartitions int, cidr *ipnet.IPNet) *gfnt.Value {
	desiredMask := calculateDesiredMask(cidrPartitions, cidr)
	refSubnetSlices := gfnt.MakeFnCIDR(gfnt.MakeFnGetAttString("VPC", "CidrBlock"), gfnt.NewInteger(cidrPartitions), gfnt.NewInteger(desiredMask))
	return refSubnetSlices
}

// To calculate the desiredMask -> ip -> 192.168.0.0/20 cidrPartition -> 6
// 32-20 -> 12 -> 2^12 -> 4096 -> 4096/cidrPartitions -> ~682 -> log2(682) -> ~9.2 -> 9 (we always floor)!
// This should result in subnets with bit size of 23! Because 32 - 9 -> 23! This is, however, calculated by
// the cloudformation CIDR function. We just need to pass in 9.
func calculateDesiredMask(cidrPartitions int, cidr *ipnet.IPNet) int {
	// We only calculate it if a custom cidr range was given
	// otherwise the hardcoded one is fine for now. Don't take my word on that.
	if cidr == nil {
		return defaultDesiredMaskSize
	}
	prefixSize, _ := cidr.Mask.Size()
	remainingCIDRBit := defaultSubnetMask - prefixSize
	remainingIPs := math.Pow(2, float64(remainingCIDRBit))
	numberOfIPsPerSubnet := remainingIPs / float64(cidrPartitions)
	return int(math.Floor(math.Log2(numberOfIPsPerSubnet)))
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
