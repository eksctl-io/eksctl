package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/weaveworks/eksctl/pkg/awsapi"

	gfncfn "github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

const (
	cfnControlPlaneSGResource         = "ControlPlaneSecurityGroup"
	cfnSharedNodeSGResource           = "ClusterSharedNodeSecurityGroup"
	cfnIngressClusterToNodeSGResource = "IngressDefaultClusterToNodeSG"
)

// A IPv4VPCResourceSet builds the resources required for the specified VPC
type IPv4VPCResourceSet struct {
	rs            *resourceSet
	clusterConfig *api.ClusterConfig
	ec2API        awsapi.EC2
	vpcID         *gfnt.Value
	subnetDetails *SubnetDetails
}

type SubnetResource struct {
	Subnet           *gfnt.Value
	RouteTable       *gfnt.Value
	AvailabilityZone string
}

type SubnetDetails struct {
	Private          []SubnetResource
	Public           []SubnetResource
	PrivateLocalZone []SubnetResource
	PublicLocalZone  []SubnetResource
}

// NewIPv4VPCResourceSet creates and returns a new VPCResourceSet
func NewIPv4VPCResourceSet(rs *resourceSet, clusterConfig *api.ClusterConfig, ec2API awsapi.EC2) *IPv4VPCResourceSet {
	return &IPv4VPCResourceSet{
		rs:            rs,
		clusterConfig: clusterConfig,
		ec2API:        ec2API,
		subnetDetails: &SubnetDetails{},
	}
}

func (v *IPv4VPCResourceSet) CreateTemplate(ctx context.Context) (*gfnt.Value, *SubnetDetails, error) {
	if err := v.addResources(); err != nil {
		return nil, nil, err
	}
	v.addOutputs(ctx)
	return v.vpcID, v.subnetDetails, nil
}

// AddResources adds all required resources
func (v *IPv4VPCResourceSet) addResources() error {
	vpc := v.clusterConfig.VPC

	v.vpcID = v.rs.newResource("VPC", &gfnec2.VPC{
		CidrBlock:          gfnt.NewString(vpc.CIDR.String()),
		EnableDnsSupport:   gfnt.True(),
		EnableDnsHostnames: gfnt.True(),
	})

	if v.isFullyPrivate() {
		v.noNAT()
		v.subnetDetails.Private = v.addSubnets(nil, api.SubnetTopologyPrivate, vpc.Subnets.Private)
		return nil
	}

	refIG := v.rs.newResource("InternetGateway", &gfnec2.InternetGateway{})
	vpcGA := "VPCGatewayAttachment"
	refPublicRT := v.rs.newResource("PublicRouteTable", &gfnec2.RouteTable{
		VpcId: v.vpcID,
	})

	if api.IsEnabled(vpc.AutoAllocateIPv6) {
		v.rs.newResource("AutoAllocatedCIDRv6", &gfnec2.VPCCidrBlock{
			VpcId:                       v.vpcID,
			AmazonProvidedIpv6CidrBlock: gfnt.True(),
		})

		v.rs.newResource(PubSubIPv6RouteKey, &gfnec2.Route{
			RouteTableId:               refPublicRT,
			DestinationIpv6CidrBlock:   gfnt.NewString(InternetIPv6CIDR),
			GatewayId:                  refIG,
			AWSCloudFormationDependsOn: []string{vpcGA},
		})
	}

	v.rs.newResource(vpcGA, &gfnec2.VPCGatewayAttachment{
		InternetGatewayId: refIG,
		VpcId:             v.vpcID,
	})

	v.rs.newResource("PublicSubnetRoute", &gfnec2.Route{
		RouteTableId:               refPublicRT,
		DestinationCidrBlock:       gfnt.NewString(InternetCIDR),
		GatewayId:                  refIG,
		AWSCloudFormationDependsOn: []string{vpcGA},
	})

	v.subnetDetails.Public = v.addSubnets(refPublicRT, api.SubnetTopologyPublic, vpc.Subnets.Public)

	if err := v.addNATGateways(); err != nil {
		return err
	}

	v.subnetDetails.Private = v.addSubnets(nil, api.SubnetTopologyPrivate, vpc.Subnets.Private)
	if vpc.LocalZoneSubnets != nil {
		if len(vpc.LocalZoneSubnets.Public) > 0 {
			v.subnetDetails.PublicLocalZone = v.addSubnets(refPublicRT, api.SubnetTopologyPublic, vpc.LocalZoneSubnets.Public)
		}
		if len(vpc.LocalZoneSubnets.Public) > 0 {
			v.subnetDetails.PrivateLocalZone = v.addSubnets(nil, api.SubnetTopologyPrivate, vpc.LocalZoneSubnets.Private)
		}
	}

	return nil
}

func (s *SubnetDetails) PublicSubnetRefs() []*gfnt.Value {
	return collectSubnetRefs(s.Public)
}

func (s *SubnetDetails) PrivateSubnetRefs() []*gfnt.Value {
	return collectSubnetRefs(s.Private)
}

func (s *SubnetDetails) PublicLocalZoneSubnetRefs() []*gfnt.Value {
	return collectSubnetRefs(s.PublicLocalZone)
}

func (s *SubnetDetails) PrivateLocalZoneSubnetRefs() []*gfnt.Value {
	return collectSubnetRefs(s.PrivateLocalZone)
}

func collectSubnetRefs(subnetResources []SubnetResource) []*gfnt.Value {
	var subnetRefs []*gfnt.Value
	for _, subnetAZ := range subnetResources {
		subnetRefs = append(subnetRefs, subnetAZ.Subnet)
	}
	return subnetRefs
}

// addOutputs adds VPC resource outputs
func (v *IPv4VPCResourceSet) addOutputs(ctx context.Context) {
	v.rs.defineOutput(outputs.ClusterVPC, v.vpcID, true, func(val string) error {
		v.clusterConfig.VPC.ID = val
		return nil
	})
	if v.clusterConfig.VPC.NAT != nil {
		v.rs.defineOutputWithoutCollector(outputs.ClusterFeatureNATMode, v.clusterConfig.VPC.NAT.Gateway, false)
	}

	addSubnetOutput := func(subnetRefs []*gfnt.Value, subnetMapping api.AZSubnetMapping, outputName string) {
		v.rs.defineJoinedOutput(outputName, subnetRefs, true, func(value string) error {
			return vpc.ImportSubnetsFromIDList(ctx, v.ec2API, v.clusterConfig, subnetMapping, strings.Split(value, ","))
		})
	}

	clusterVPC := v.clusterConfig.VPC
	if subnetAZs := v.subnetDetails.PrivateSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, clusterVPC.Subnets.Private, outputs.ClusterSubnetsPrivate)
	}

	if subnetAZs := v.subnetDetails.PublicSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, clusterVPC.Subnets.Public, outputs.ClusterSubnetsPublic)
	}

	if subnetAZs := v.subnetDetails.PrivateLocalZoneSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, clusterVPC.LocalZoneSubnets.Private, outputs.ClusterSubnetsPrivateLocal)
	}

	if subnetAZs := v.subnetDetails.PublicLocalZoneSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, clusterVPC.LocalZoneSubnets.Public, outputs.ClusterSubnetsPublicLocal)
	}

	if v.isFullyPrivate() {
		v.rs.defineOutputWithoutCollector(outputs.ClusterFullyPrivate, true, true)
	}
}

func (v *IPv4VPCResourceSet) isFullyPrivate() bool {
	return v.clusterConfig.PrivateCluster.Enabled
}

// RenderJSON returns the rendered JSON
func (v *IPv4VPCResourceSet) RenderJSON() ([]byte, error) {
	return v.rs.renderJSON()
}

func (v *IPv4VPCResourceSet) addSubnets(refRT *gfnt.Value, topology api.SubnetTopology, subnets map[string]api.AZSubnetSpec) []SubnetResource {
	autoAllocateIPV6 := api.IsEnabled(v.clusterConfig.VPC.AutoAllocateIPv6)

	var subnetResources []SubnetResource

	for name, s := range subnets {
		az := s.AZ
		nameAlias := strings.ToUpper(strings.Join(strings.Split(name, "-"), ""))
		subnet := &gfnec2.Subnet{
			AvailabilityZone: gfnt.NewString(az),
			CidrBlock:        gfnt.NewString(s.CIDR.String()),
			VpcId:            v.vpcID,
		}

		switch topology {
		case api.SubnetTopologyPrivate:
			// Choose the appropriate route table for private subnets.
			refRT = gfnt.MakeRef("PrivateRouteTable" + nameAlias)
			subnet.Tags = []gfncfn.Tag{{
				Key:   gfnt.NewString("kubernetes.io/role/internal-elb"),
				Value: gfnt.NewString("1"),
			}}
		case api.SubnetTopologyPublic:
			subnet.Tags = []gfncfn.Tag{{
				Key:   gfnt.NewString("kubernetes.io/role/elb"),
				Value: gfnt.NewString("1"),
			}}
			subnet.MapPublicIpOnLaunch = gfnt.True()
		}

		subnetAlias := string(topology) + nameAlias
		refSubnet := v.rs.newResource("Subnet"+subnetAlias, subnet)
		v.rs.newResource("RouteTableAssociation"+subnetAlias, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     refSubnet,
			RouteTableId: refRT,
		})

		if autoAllocateIPV6 {
			refSubnetSlices := getSubnetIPv6CIDRBlock((len(v.clusterConfig.AvailabilityZones) * 2) + 2)
			v.rs.newResource(subnetAlias+"CIDRv6", &gfnec2.SubnetCidrBlock{
				SubnetId:      refSubnet,
				Ipv6CidrBlock: gfnt.MakeFnSelect(gfnt.NewInteger(s.CIDRIndex), refSubnetSlices),
			})
		}

		subnetResources = append(subnetResources, SubnetResource{
			AvailabilityZone: az,
			RouteTable:       refRT,
			Subnet:           refSubnet,
		})
	}
	return subnetResources
}

func (v *IPv4VPCResourceSet) addNATGateways() error {
	switch *v.clusterConfig.VPC.NAT.Gateway {
	case api.ClusterHighlyAvailableNAT:
		v.haNAT()
	case api.ClusterSingleNAT:
		v.singleNAT()
	case api.ClusterDisableNAT:
		v.noNAT()
	default:
		// TODO validate this before starting to add resources
		return fmt.Errorf("%s is not a valid NAT gateway mode", *v.clusterConfig.VPC.NAT.Gateway)
	}
	return nil
}

var (
	sgProtoTCP           = gfnt.NewString("tcp")
	sgSourceAnywhereIPv4 = gfnt.NewString("0.0.0.0/0")
	sgSourceAnywhereIPv6 = gfnt.NewString("::/0")

	sgPortZero    = gfnt.NewInteger(0)
	sgMinNodePort = gfnt.NewInteger(1025)
	sgMaxNodePort = gfnt.NewInteger(65535)

	sgPortHTTPS = gfnt.NewInteger(443)
	sgPortSSH   = gfnt.NewInteger(22)
)

type clusterSecurityGroup struct {
	ControlPlane      *gfnt.Value
	ClusterSharedNode *gfnt.Value
}

func (v *IPv4VPCResourceSet) haNAT() {
	for _, az := range v.clusterConfig.AvailabilityZones {
		alphanumericUpperAZ := formatAZ(az)

		// Allocate an EIP
		v.rs.newResource("NATIP"+alphanumericUpperAZ, &gfnec2.EIP{
			Domain: gfnt.NewString("vpc"),
		})
		// Allocate a NAT gateway in the public subnet
		refNG := v.rs.newResource("NATGateway"+alphanumericUpperAZ, &gfnec2.NatGateway{
			AllocationId: gfnt.MakeFnGetAttString("NATIP"+alphanumericUpperAZ, "AllocationId"),
			SubnetId:     gfnt.MakeRef("SubnetPublic" + alphanumericUpperAZ),
		})

		// Allocate a routing table for the private subnet
		refRT := v.rs.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfnec2.RouteTable{
			VpcId: v.vpcID,
		})
		// Create a route that sends Internet traffic through the NAT gateway
		v.rs.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfnec2.Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: gfnt.NewString(InternetCIDR),
			NatGatewayId:         refNG,
		})
		// Associate the routing table with the subnet
		v.rs.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func (v *IPv4VPCResourceSet) singleNAT() {
	sortedAZs := v.clusterConfig.AvailabilityZones
	firstUpperAZ := strings.ToUpper(strings.Join(strings.Split(sortedAZs[0], "-"), ""))

	v.rs.newResource("NATIP", &gfnec2.EIP{
		Domain: gfnt.NewString("vpc"),
	})
	refNG := v.rs.newResource("NATGateway", &gfnec2.NatGateway{
		AllocationId: gfnt.MakeFnGetAttString("NATIP", "AllocationId"),
		SubnetId:     gfnt.MakeRef("SubnetPublic" + firstUpperAZ),
	})

	for _, zone := range makeAllZones(v.clusterConfig) {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(zone, "-"), ""))

		refRT := v.rs.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfnec2.RouteTable{
			VpcId: v.vpcID,
		})

		v.rs.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfnec2.Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: gfnt.NewString(InternetCIDR),
			NatGatewayId:         refNG,
		})
		v.rs.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func (v *IPv4VPCResourceSet) noNAT() {
	for _, az := range makeAllZones(v.clusterConfig) {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := v.rs.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfnec2.RouteTable{
			VpcId: v.vpcID,
		})
		v.rs.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func makeAllZones(clusterConfig *api.ClusterConfig) []string {
	zones := make([]string, len(clusterConfig.AvailabilityZones)+len(clusterConfig.LocalZones))
	copy(zones, clusterConfig.AvailabilityZones)
	copy(zones[len(clusterConfig.AvailabilityZones):], clusterConfig.LocalZones)
	return zones
}
