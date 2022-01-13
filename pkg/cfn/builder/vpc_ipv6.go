package builder

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// A IPv6VPCResourceSet builds the resources required for the specified VPC
type IPv6VPCResourceSet struct {
	rs            *resourceSet
	clusterConfig *api.ClusterConfig
	ec2API        ec2iface.EC2API
}

// NewIPv6VPCResourceSet creates and returns a new VPCResourceSet
func NewIPv6VPCResourceSet(rs *resourceSet, clusterConfig *api.ClusterConfig, ec2API ec2iface.EC2API) *IPv6VPCResourceSet {
	return &IPv6VPCResourceSet{
		rs:            rs,
		clusterConfig: clusterConfig,
		ec2API:        ec2API,
	}
}

func (v *IPv6VPCResourceSet) CreateTemplate() (*gfnt.Value, *SubnetDetails, error) {
	var publicSubnetResourceRefs, privateSubnetResourceRefs []*gfnt.Value
	vpcResourceRef := v.rs.newResource(VPCResourceKey, &gfnec2.VPC{
		CidrBlock:          gfnt.NewString(v.clusterConfig.VPC.CIDR.String()),
		EnableDnsSupport:   gfnt.True(),
		EnableDnsHostnames: gfnt.True(),
	})
	v.rs.defineOutput(outputs.ClusterVPC, vpcResourceRef, true, func(val string) error {
		v.clusterConfig.VPC.ID = val
		return nil
	})

	v.addIpv6CidrBlock()

	addSubnetOutput := func(subnetRefs []*gfnt.Value, topology api.SubnetTopology, outputName string) {
		v.rs.defineJoinedOutput(outputName, subnetRefs, true, func(value string) error {
			return vpc.ImportSubnetsFromIDList(v.ec2API, v.clusterConfig, topology, strings.Split(value, ","))
		})
	}

	var privateSubnets []SubnetResource
	cidrPartitions := (len(v.clusterConfig.AvailabilityZones) * 2) + 2
	for i, az := range v.clusterConfig.AvailabilityZones {
		azFormatted := formatAZ(az)
		rtRef := v.rs.newResource(PrivateRouteTableKey+azFormatted, &gfnec2.RouteTable{
			VpcId: gfnt.MakeRef(VPCResourceKey),
		})

		subnet := v.createSubnet(az, azFormatted, i+len(v.clusterConfig.AvailabilityZones), cidrPartitions, true)
		privateSubnetResourceRefs = append(privateSubnetResourceRefs, subnet)
		privateSubnets = append(privateSubnets, SubnetResource{
			Subnet:           subnet,
			AvailabilityZone: az,
			RouteTable:       rtRef,
		})

		v.rs.newResource(PrivateRouteTableAssociation+azFormatted, &gfnec2.SubnetRouteTableAssociation{
			RouteTableId: rtRef,
			SubnetId:     subnet,
		})
	}
	addSubnetOutput(privateSubnetResourceRefs, api.SubnetTopologyPrivate, outputs.ClusterSubnetsPrivate)

	if v.isFullyPrivate() {
		return vpcResourceRef, &SubnetDetails{
			Private: privateSubnets,
		}, nil
	}

	// add the rest of the public resources.
	refIGW := v.rs.newResource(IGWKey, &gfnec2.InternetGateway{})

	v.rs.newResource(GAKey, &gfnec2.VPCGatewayAttachment{
		InternetGatewayId: gfnt.MakeRef(IGWKey),
		VpcId:             gfnt.MakeRef(VPCResourceKey),
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

	v.rs.newResource(EgressOnlyInternetGatewayKey, &gfnec2.EgressOnlyInternetGateway{
		VpcId: gfnt.MakeRef(VPCResourceKey),
	})
	var publicSubnets []SubnetResource
	for i, az := range v.clusterConfig.AvailabilityZones {
		azFormatted := formatAZ(az)

		subnet := v.createSubnet(az, azFormatted, i, cidrPartitions, false)
		publicSubnets = append(publicSubnets, SubnetResource{Subnet: subnet, AvailabilityZone: az})
		publicSubnetResourceRefs = append(publicSubnetResourceRefs, subnet)

		v.rs.newResource(PubRouteTableAssociation+azFormatted, &gfnec2.SubnetRouteTableAssociation{
			RouteTableId: gfnt.MakeRef(PubRouteTableKey),
			SubnetId:     gfnt.MakeRef(PublicSubnetKey + azFormatted),
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
	}
	addSubnetOutput(publicSubnetResourceRefs, api.SubnetTopologyPublic, outputs.ClusterSubnetsPublic)

	return vpcResourceRef, &SubnetDetails{
		Private: privateSubnets,
		Public:  publicSubnets,
	}, nil
}

func (v *IPv6VPCResourceSet) addIpv6CidrBlock() {
	if v.clusterConfig.VPC.IPv6Cidr != "" {
		v.rs.newResource(IPv6CIDRBlockKey, &gfnec2.VPCCidrBlock{
			AmazonProvidedIpv6CidrBlock: gfnt.False(),
			Ipv6CidrBlock:               gfnt.NewString(v.clusterConfig.VPC.IPv6Cidr),
			Ipv6Pool:                    gfnt.NewString(v.clusterConfig.VPC.IPv6Pool),
			VpcId:                       gfnt.MakeRef(VPCResourceKey),
		})
		return
	}

	v.rs.newResource(IPv6CIDRBlockKey, &gfnec2.VPCCidrBlock{
		AmazonProvidedIpv6CidrBlock: gfnt.True(),
		VpcId:                       gfnt.MakeRef(VPCResourceKey),
	})
}

func (v *IPv6VPCResourceSet) RenderJSON() ([]byte, error) {
	return v.rs.renderJSON()
}

func (v *IPv6VPCResourceSet) isFullyPrivate() bool {
	return v.clusterConfig.PrivateCluster.Enabled
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
