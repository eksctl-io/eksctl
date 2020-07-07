package builder

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/pkg/errors"

	gfn "github.com/weaveworks/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

var internetCIDR = gfn.NewString("0.0.0.0/0")

const (
	cfnControlPlaneSGResource         = "ControlPlaneSecurityGroup"
	cfnSharedNodeSGResource           = "ClusterSharedNodeSecurityGroup"
	cfnIngressClusterToNodeSGResource = "IngressDefaultClusterToNodeSG"
)

// A VPCResourceSet builds the resources required for the specified VPC
type VPCResourceSet struct {
	*resourceSet
	clusterConfig *api.ClusterConfig
	provider      api.ClusterProvider

	vpcResource *VPCResource
}

// VPCResource represents a VPC resource
type VPCResource struct {
	VPC           *gfn.Value
	SubnetDetails *subnetDetails
}

type subnetResource struct {
	Subnet           *gfn.Value
	RouteTable       *gfn.Value
	AvailabilityZone string
}

type subnetDetails struct {
	Private []subnetResource
	Public  []subnetResource
}

func (s *subnetDetails) PublicSubnetRefs() []*gfn.Value {
	var subnetRefs []*gfn.Value
	for _, subnetAZ := range s.Public {
		subnetRefs = append(subnetRefs, subnetAZ.Subnet)
	}
	return subnetRefs
}

func (s *subnetDetails) PrivateSubnetRefs() []*gfn.Value {
	var subnetRefs []*gfn.Value
	for _, subnetAZ := range s.Private {
		subnetRefs = append(subnetRefs, subnetAZ.Subnet)
	}
	return subnetRefs
}

// NewVPCResourceSet creates and returns a new VPCResourceSet
func NewVPCResourceSet(rs *resourceSet, clusterConfig *api.ClusterConfig, provider api.ClusterProvider) *VPCResourceSet {
	var vpcRef *gfn.Value
	if clusterConfig.VPC.ID == "" {
		vpcRef = rs.newResource("VPC", &gfn.AWSEC2VPC{
			CidrBlock:          gfn.NewString(clusterConfig.VPC.CIDR.String()),
			EnableDnsSupport:   gfn.True(),
			EnableDnsHostnames: gfn.True(),
		})
	} else {
		vpcRef = gfn.NewString(clusterConfig.VPC.ID)
	}

	return &VPCResourceSet{
		resourceSet:   rs,
		clusterConfig: clusterConfig,
		provider:      provider,

		vpcResource: &VPCResource{
			VPC:           vpcRef,
			SubnetDetails: &subnetDetails{},
		},
	}
}

// AddResources adds all required resources
func (v *VPCResourceSet) AddResources() (*VPCResource, error) {
	vpc := v.clusterConfig.VPC
	if customVPC := vpc.ID != ""; customVPC {
		if err := v.importResources(); err != nil {
			return nil, errors.Wrap(err, "error importing VPC resources")
		}
		return v.vpcResource, nil
	}

	if api.IsEnabled(vpc.AutoAllocateIPv6) {
		v.newResource("AutoAllocatedCIDRv6", &gfn.AWSEC2VPCCidrBlock{
			VpcId:                       v.vpcResource.VPC,
			AmazonProvidedIpv6CidrBlock: gfn.True(),
		})
	}

	if v.isFullyPrivate() {
		v.noNAT()
		v.vpcResource.SubnetDetails.Private = v.addSubnets(nil, api.SubnetTopologyPrivate, vpc.Subnets.Private)
		return v.vpcResource, nil
	}

	refIG := v.newResource("InternetGateway", &gfn.AWSEC2InternetGateway{})
	vpcGA := "VPCGatewayAttachment"
	v.newResource(vpcGA, &gfn.AWSEC2VPCGatewayAttachment{
		InternetGatewayId: refIG,
		VpcId:             v.vpcResource.VPC,
	})

	refPublicRT := v.newResource("PublicRouteTable", &gfn.AWSEC2RouteTable{
		VpcId: v.vpcResource.VPC,
	})

	v.newResource("PublicSubnetRoute", &route{
		AWSEC2Route: gfn.AWSEC2Route{
			RouteTableId:         refPublicRT,
			DestinationCidrBlock: internetCIDR,
			GatewayId:            refIG,
		},
		DependsOn: []string{vpcGA},
	})

	v.vpcResource.SubnetDetails.Public = v.addSubnets(refPublicRT, api.SubnetTopologyPublic, vpc.Subnets.Public)

	if err := v.addNATGateways(); err != nil {
		return nil, err
	}

	v.vpcResource.SubnetDetails.Private = v.addSubnets(nil, api.SubnetTopologyPrivate, vpc.Subnets.Private)
	return v.vpcResource, nil
}

func (v *VPCResourceSet) addSubnets(refRT *gfn.Value, topology api.SubnetTopology, subnets map[string]api.Network) []subnetResource {
	var subnetIndexForIPv6 int
	if api.IsEnabled(v.clusterConfig.VPC.AutoAllocateIPv6) {
		// this is same kind of indexing we have in vpc.SetSubnets
		switch topology {
		case api.SubnetTopologyPrivate:
			subnetIndexForIPv6 = len(v.clusterConfig.AvailabilityZones)
		case api.SubnetTopologyPublic:
			subnetIndexForIPv6 = 0
		}
	}

	var subnetResources []subnetResource

	for az, subnet := range subnets {
		alias := string(topology) + strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))
		subnet := &gfn.AWSEC2Subnet{
			AvailabilityZone: gfn.NewString(az),
			CidrBlock:        gfn.NewString(subnet.CIDR.String()),
			VpcId:            v.vpcResource.VPC,
		}

		switch topology {
		case api.SubnetTopologyPrivate:
			// Choose the appropriate route table for private subnets
			refRT = gfn.MakeRef("PrivateRouteTable" + strings.ToUpper(strings.Join(strings.Split(az, "-"), "")))
			subnet.Tags = []gfn.Tag{{
				Key:   gfn.NewString("kubernetes.io/role/internal-elb"),
				Value: gfn.NewString("1"),
			}}
		case api.SubnetTopologyPublic:
			subnet.Tags = []gfn.Tag{{
				Key:   gfn.NewString("kubernetes.io/role/elb"),
				Value: gfn.NewString("1"),
			}}
			subnet.MapPublicIpOnLaunch = gfn.True()
		}
		refSubnet := v.newResource("Subnet"+alias, subnet)
		v.newResource("RouteTableAssociation"+alias, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     refSubnet,
			RouteTableId: refRT,
		})

		if api.IsEnabled(v.clusterConfig.VPC.AutoAllocateIPv6) {
			// get 8 of /64 subnets from the auto-allocated IPv6 block,
			// and pick one block based on subnetIndexForIPv6 counter;
			// NOTE: this is done inside of CloudFormation using Fn::Cidr,
			// we don't slice it here, just construct the JSON expression
			// that does slicing at runtime.
			refAutoAllocateCIDRv6 := gfn.MakeFnSelect(
				0, gfn.MakeFnGetAttString("VPC.Ipv6CidrBlocks"),
			)
			refSubnetSlices := gfn.MakeFnCIDR(
				refAutoAllocateCIDRv6, 8, 64,
			)
			v.newResource(alias+"CIDRv6", &gfn.AWSEC2SubnetCidrBlock{
				SubnetId:      refSubnet,
				Ipv6CidrBlock: gfn.MakeFnSelect(subnetIndexForIPv6, refSubnetSlices),
			})
			subnetIndexForIPv6++
		}

		subnetResources = append(subnetResources, subnetResource{
			AvailabilityZone: az,
			RouteTable:       refRT,
			Subnet:           refSubnet,
		})
	}
	return subnetResources
}

// route adds DependsOn support to the AWSEC2Route struct
type route struct {
	AWSEC2Route gfn.AWSEC2Route
	DependsOn   []string
}

// MarshalJSON is a custom JSON marshalling hook that adds DependsOn to the
// legacy goformation struct AWSEC2Route
func (r *route) MarshalJSON() ([]byte, error) {
	type Properties gfn.AWSEC2Route
	return json.Marshal(&struct {
		Type       string
		Properties Properties
		DependsOn  []string
	}{
		Type:       r.AWSEC2Route.AWSCloudFormationType(),
		Properties: (Properties)(r.AWSEC2Route),
		DependsOn:  r.DependsOn,
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that adds DependsOn to the
// legacy goformation struct AWSEC2Route
func (r *route) UnmarshalJSON(b []byte) error {
	type Properties gfn.AWSEC2Route
	res := &struct {
		Type       string
		Properties *Properties
		DependsOn  *[]string
	}{}
	if err := json.Unmarshal(b, &res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		r.AWSEC2Route = gfn.AWSEC2Route(*res.Properties)
	}
	if res.DependsOn != nil {
		r.DependsOn = *res.DependsOn
	}

	return nil
}

func (v *VPCResourceSet) addNATGateways() error {
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

func (v *VPCResourceSet) importResources() error {
	makeSubnetResources := func(subnets map[string]api.Network, subnetRoutes map[string]string) ([]subnetResource, error) {
		subnetResources := make([]subnetResource, len(subnets))
		i := 0
		for az, network := range subnets {
			sr := subnetResource{
				AvailabilityZone: az,
				Subnet:           gfn.NewString(network.ID),
			}

			if subnetRoutes != nil {
				rt, ok := subnetRoutes[network.ID]
				if !ok {
					return nil, errors.Errorf("failed to find an explicit route table associated with subnet: %q;"+
						"eksctl does not modify the main route table if a subnet is not associated with an explicit route table", network.ID)
				}
				sr.RouteTable = gfn.NewString(rt)
			}
			subnetResources[i] = sr
			i++
		}
		return subnetResources, nil
	}

	if subnets := v.clusterConfig.VPC.Subnets.Private; subnets != nil {
		var (
			subnetRoutes map[string]string
			err          error
		)
		if v.isFullyPrivate() {
			subnetRoutes, err = importRouteTables(v.provider.EC2(), v.clusterConfig.VPC.Subnets.Private)
			if err != nil {
				return err
			}
		}

		subnetResources, err := makeSubnetResources(subnets, subnetRoutes)
		if err != nil {
			return err
		}
		v.vpcResource.SubnetDetails.Private = subnetResources
	}

	if subnets := v.clusterConfig.VPC.Subnets.Public; subnets != nil {
		subnetResources, err := makeSubnetResources(subnets, nil)
		if err != nil {
			return err
		}
		v.vpcResource.SubnetDetails.Public = subnetResources
	}

	return nil
}

func importRouteTables(ec2API ec2iface.EC2API, subnets map[string]api.Network) (map[string]string, error) {
	var subnetIDs []string
	for id := range subnets {
		subnetIDs = append(subnetIDs, id)
	}

	var routeTables []*ec2.RouteTable
	var nextToken *string

	for {
		output, err := ec2API.DescribeRouteTables(&ec2.DescribeRouteTablesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("association.subnet-id"),
					Values: aws.StringSlice(subnetIDs),
				},
			},
			NextToken: nextToken,
		})

		if err != nil {
			return nil, errors.Wrap(err, "error describing route tables")
		}

		routeTables = append(routeTables, output.RouteTables...)

		if nextToken = output.NextToken; nextToken == nil {
			break
		}
	}

	subnetRoutes := make(map[string]string)
	for _, rt := range routeTables {
		for _, rta := range rt.Associations {
			subnetRoutes[*rta.SubnetId] = *rt.RouteTableId
		}
	}
	return subnetRoutes, nil
}

// AddOutputs adds VPC resource outputs
func (v *VPCResourceSet) AddOutputs() {
	v.defineOutput(outputs.ClusterVPC, v.vpcResource.VPC, true, func(val string) error {
		v.clusterConfig.VPC.ID = val
		return nil
	})
	if v.clusterConfig.VPC.NAT != nil {
		v.defineOutputWithoutCollector(outputs.ClusterFeatureNATMode, v.clusterConfig.VPC.NAT.Gateway, false)
	}

	addSubnetOutput := func(subnetRefs []*gfn.Value, topology api.SubnetTopology, outputName string) {
		v.defineJoinedOutput(outputName, subnetRefs, true, func(value string) error {
			return vpc.ImportSubnetsFromList(v.provider, v.clusterConfig, topology, strings.Split(value, ","))
		})
	}

	if subnetAZs := v.vpcResource.SubnetDetails.PrivateSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, api.SubnetTopologyPrivate, outputs.ClusterSubnetsPrivate)
	}

	if subnetAZs := v.vpcResource.SubnetDetails.PublicSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, api.SubnetTopologyPublic, outputs.ClusterSubnetsPublic)
	}

	if v.isFullyPrivate() {
		v.defineOutputWithoutCollector(outputs.ClusterFullyPrivate, true, true)
	}
}

func (v *VPCResourceSet) isFullyPrivate() bool {
	return v.clusterConfig.PrivateCluster.Enabled
}

var (
	sgProtoTCP           = gfn.NewString("tcp")
	sgSourceAnywhereIPv4 = gfn.NewString("0.0.0.0/0")
	sgSourceAnywhereIPv6 = gfn.NewString("::/0")

	sgPortZero    = gfn.NewInteger(0)
	sgMinNodePort = gfn.NewInteger(1025)
	sgMaxNodePort = gfn.NewInteger(65535)

	sgPortHTTPS = gfn.NewInteger(443)
	sgPortSSH   = gfn.NewInteger(22)
)

type clusterSecurityGroup struct {
	ControlPlane      *gfn.Value
	ClusterSharedNode *gfn.Value
}

func (c *ClusterResourceSet) addResourcesForSecurityGroups(vpcResource *VPCResource) *clusterSecurityGroup {
	var refControlPlaneSG, refClusterSharedNodeSG *gfn.Value

	if c.spec.VPC.SecurityGroup == "" {
		refControlPlaneSG = c.newResource(cfnControlPlaneSGResource, &gfn.AWSEC2SecurityGroup{
			GroupDescription: gfn.NewString("Communication between the control plane and worker nodegroups"),
			VpcId:            vpcResource.VPC,
		})
	} else {
		refControlPlaneSG = gfn.NewString(c.spec.VPC.SecurityGroup)
	}
	c.securityGroups = []*gfn.Value{refControlPlaneSG} // only this one SG is passed to EKS API, nodes are isolated

	if c.spec.VPC.SharedNodeSecurityGroup == "" {
		refClusterSharedNodeSG = c.newResource(cfnSharedNodeSGResource, &gfn.AWSEC2SecurityGroup{
			GroupDescription: gfn.NewString("Communication between all nodes in the cluster"),
			VpcId:            vpcResource.VPC,
		})
		c.newResource("IngressInterNodeGroupSG", &gfn.AWSEC2SecurityGroupIngress{
			GroupId:               refClusterSharedNodeSG,
			SourceSecurityGroupId: refClusterSharedNodeSG,
			Description:           gfn.NewString("Allow nodes to communicate with each other (all ports)"),
			IpProtocol:            gfn.NewString("-1"),
			FromPort:              sgPortZero,
			ToPort:                sgMaxNodePort,
		})
		if c.supportsManagedNodes {
			// To enable communication between both managed and unmanaged nodegroups, this allows ingress traffic from
			// the default cluster security group ID that EKS creates by default
			// EKS attaches this to Managed Nodegroups by default, but we need to handle this for unmanaged nodegroups
			c.newResource(cfnIngressClusterToNodeSGResource, &gfn.AWSEC2SecurityGroupIngress{
				GroupId:               refClusterSharedNodeSG,
				SourceSecurityGroupId: gfn.MakeFnGetAttString(makeAttrAccessor("ControlPlane", outputs.ClusterDefaultSecurityGroup)),
				Description:           gfn.NewString("Allow managed and unmanaged nodes to communicate with each other (all ports)"),
				IpProtocol:            gfn.NewString("-1"),
				FromPort:              sgPortZero,
				ToPort:                sgMaxNodePort,
			})
			c.newResource("IngressNodeToDefaultClusterSG", &gfn.AWSEC2SecurityGroupIngress{
				GroupId:               gfn.MakeFnGetAttString(makeAttrAccessor("ControlPlane", outputs.ClusterDefaultSecurityGroup)),
				SourceSecurityGroupId: refClusterSharedNodeSG,
				Description:           gfn.NewString("Allow unmanaged nodes to communicate with control plane (all ports)"),
				IpProtocol:            gfn.NewString("-1"),
				FromPort:              sgPortZero,
				ToPort:                sgMaxNodePort,
			})
		}
	} else {
		refClusterSharedNodeSG = gfn.NewString(c.spec.VPC.SharedNodeSecurityGroup)
	}

	if c.spec.VPC == nil {
		c.spec.VPC = &api.ClusterVPC{}
	}
	c.rs.defineOutput(outputs.ClusterSecurityGroup, refControlPlaneSG, true, func(v string) error {
		c.spec.VPC.SecurityGroup = v
		return nil
	})
	c.rs.defineOutput(outputs.ClusterSharedNodeSecurityGroup, refClusterSharedNodeSG, true, func(v string) error {
		c.spec.VPC.SharedNodeSecurityGroup = v
		return nil
	})

	return &clusterSecurityGroup{
		ControlPlane:      refControlPlaneSG,
		ClusterSharedNode: refClusterSharedNodeSG,
	}
}

func (n *NodeGroupResourceSet) addResourcesForSecurityGroups() {
	for _, id := range n.spec.SecurityGroups.AttachIDs {
		n.securityGroups = append(n.securityGroups, gfn.NewString(id))
	}

	if api.IsEnabled(n.spec.SecurityGroups.WithShared) {
		refClusterSharedNodeSG := makeImportValue(n.clusterStackName, outputs.ClusterSharedNodeSecurityGroup)
		n.securityGroups = append(n.securityGroups, refClusterSharedNodeSG)
	}

	if api.IsDisabled(n.spec.SecurityGroups.WithLocal) {
		return
	}

	desc := "worker nodes in group " + n.nodeGroupName

	allInternalIPv4 := gfn.NewString(n.clusterSpec.VPC.CIDR.String())

	refControlPlaneSG := makeImportValue(n.clusterStackName, outputs.ClusterSecurityGroup)

	refNodeGroupLocalSG := n.newResource("SG", &gfn.AWSEC2SecurityGroup{
		VpcId:            makeImportValue(n.clusterStackName, outputs.ClusterVPC),
		GroupDescription: gfn.NewString("Communication between the control plane and " + desc),
		Tags: []gfn.Tag{{
			Key:   gfn.NewString("kubernetes.io/cluster/" + n.clusterSpec.Metadata.Name),
			Value: gfn.NewString("owned"),
		}},
	})

	n.securityGroups = append(n.securityGroups, refNodeGroupLocalSG)

	n.newResource("IngressInterCluster", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refNodeGroupLocalSG,
		SourceSecurityGroupId: refControlPlaneSG,
		Description:           gfn.NewString("Allow " + desc + " to communicate with control plane (kubelet and workload TCP ports)"),
		IpProtocol:            sgProtoTCP,
		FromPort:              sgMinNodePort,
		ToPort:                sgMaxNodePort,
	})
	n.newResource("EgressInterCluster", &gfn.AWSEC2SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                gfn.NewString("Allow control plane to communicate with " + desc + " (kubelet and workload TCP ports)"),
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgMinNodePort,
		ToPort:                     sgMaxNodePort,
	})
	n.newResource("IngressInterClusterAPI", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refNodeGroupLocalSG,
		SourceSecurityGroupId: refControlPlaneSG,
		Description:           gfn.NewString("Allow " + desc + " to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)"),
		IpProtocol:            sgProtoTCP,
		FromPort:              sgPortHTTPS,
		ToPort:                sgPortHTTPS,
	})
	n.newResource("EgressInterClusterAPI", &gfn.AWSEC2SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                gfn.NewString("Allow control plane to communicate with " + desc + " (workloads using HTTPS port, commonly used with extension API servers)"),
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgPortHTTPS,
		ToPort:                     sgPortHTTPS,
	})
	n.newResource("IngressInterClusterCP", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refControlPlaneSG,
		SourceSecurityGroupId: refNodeGroupLocalSG,
		Description:           gfn.NewString("Allow control plane to receive API requests from " + desc),
		IpProtocol:            sgProtoTCP,
		FromPort:              sgPortHTTPS,
		ToPort:                sgPortHTTPS,
	})
	if *n.spec.SSH.Allow {
		if n.spec.PrivateNetworking {
			n.newResource("SSHIPv4", &gfn.AWSEC2SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIp:      allInternalIPv4,
				Description: gfn.NewString("Allow SSH access to " + desc + " (private, only inside VPC)"),
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
		} else {
			n.newResource("SSHIPv4", &gfn.AWSEC2SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIp:      sgSourceAnywhereIPv4,
				Description: gfn.NewString("Allow SSH access to " + desc),
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
			n.newResource("SSHIPv6", &gfn.AWSEC2SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIpv6:    sgSourceAnywhereIPv6,
				Description: gfn.NewString("Allow SSH access to " + desc),
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
		}
	}
}

func (v *VPCResourceSet) haNAT() {
	for _, az := range v.clusterConfig.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		// Allocate an EIP
		v.newResource("NATIP"+alphanumericUpperAZ, &gfn.AWSEC2EIP{
			Domain: gfn.NewString("vpc"),
		})
		// Allocate a NAT gateway in the public subnet
		refNG := v.newResource("NATGateway"+alphanumericUpperAZ, &gfn.AWSEC2NatGateway{
			AllocationId: gfn.MakeFnGetAttString("NATIP" + alphanumericUpperAZ + ".AllocationId"),
			SubnetId:     gfn.MakeRef("SubnetPublic" + alphanumericUpperAZ),
		})

		// Allocate a routing table for the private subnet
		refRT := v.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfn.AWSEC2RouteTable{
			VpcId: v.vpcResource.VPC,
		})
		// Create a route that sends Internet traffic through the NAT gateway
		v.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfn.AWSEC2Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		// Associate the routing table with the subnet
		v.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     gfn.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}

}

func (v *VPCResourceSet) singleNAT() {
	sortedAZs := v.clusterConfig.AvailabilityZones
	firstUpperAZ := strings.ToUpper(strings.Join(strings.Split(sortedAZs[0], "-"), ""))

	v.newResource("NATIP", &gfn.AWSEC2EIP{
		Domain: gfn.NewString("vpc"),
	})
	refNG := v.newResource("NATGateway", &gfn.AWSEC2NatGateway{
		AllocationId: gfn.MakeFnGetAttString("NATIP.AllocationId"),
		SubnetId:     gfn.MakeRef("SubnetPublic" + firstUpperAZ),
	})

	for _, az := range v.clusterConfig.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := v.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfn.AWSEC2RouteTable{
			VpcId: v.vpcResource.VPC,
		})

		v.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfn.AWSEC2Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		v.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     gfn.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func (v *VPCResourceSet) noNAT() {
	for _, az := range v.clusterConfig.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := v.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfn.AWSEC2RouteTable{
			VpcId: v.vpcResource.VPC,
		})
		v.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     gfn.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}
