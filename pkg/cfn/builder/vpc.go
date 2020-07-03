package builder

import (
	"fmt"
	"strings"

	gfncfn "github.com/weaveworks/goformation/v4/cloudformation/cloudformation"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

var internetCIDR = gfnt.NewString("0.0.0.0/0")

const (
	cfnControlPlaneSGResource         = "ControlPlaneSecurityGroup"
	cfnSharedNodeSGResource           = "ClusterSharedNodeSecurityGroup"
	cfnIngressClusterToNodeSGResource = "IngressDefaultClusterToNodeSG"
)

func (c *ClusterResourceSet) addSubnets(refRT *gfnt.Value, topology api.SubnetTopology, subnets map[string]api.Network) {
	var subnetIndexForIPv6 int
	if api.IsEnabled(c.spec.VPC.AutoAllocateIPv6) {
		// this is same kind of indexing we have in vpc.SetSubnets
		switch topology {
		case api.SubnetTopologyPrivate:
			subnetIndexForIPv6 = len(c.spec.AvailabilityZones)
		case api.SubnetTopologyPublic:
			subnetIndexForIPv6 = 0
		}
	}

	for az, subnet := range subnets {
		alias := string(topology) + strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))
		subnet := &gfnec2.Subnet{
			AvailabilityZone: gfnt.NewString(az),
			CidrBlock:        gfnt.NewString(subnet.CIDR.String()),
			VpcId:            c.vpc,
		}

		switch topology {
		case api.SubnetTopologyPrivate:
			// Choose the appropriate route table for private subnets
			refRT = gfnt.MakeRef("PrivateRouteTable" + strings.ToUpper(strings.Join(strings.Split(az, "-"), "")))
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
		refSubnet := c.newResource("Subnet"+alias, subnet)
		c.newResource("RouteTableAssociation"+alias, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     refSubnet,
			RouteTableId: refRT,
		})

		if api.IsEnabled(c.spec.VPC.AutoAllocateIPv6) {
			// get 8 of /64 subnets from the auto-allocated IPv6 block,
			// and pick one block based on subnetIndexForIPv6 counter;
			// NOTE: this is done inside of CloudFormation using Fn::Cidr,
			// we don't slice it here, just construct the JSON expression
			// that does slicing at runtime.
			refAutoAllocateCIDRv6 := gfnt.MakeFnSelect(
				gfnt.NewInteger(0), gfnt.MakeFnGetAttString("VPC", "Ipv6CidrBlocks"),
			)
			refSubnetSlices := gfnt.MakeFnCIDR(
				refAutoAllocateCIDRv6, gfnt.NewInteger(8), gfnt.NewInteger(64),
			)
			c.newResource(alias+"CIDRv6", &gfnec2.SubnetCidrBlock{
				SubnetId:      refSubnet,
				Ipv6CidrBlock: gfnt.MakeFnSelect(gfnt.NewInteger(subnetIndexForIPv6), refSubnetSlices),
			})
			subnetIndexForIPv6++
		}

		c.subnets[topology] = append(c.subnets[topology], refSubnet)
	}
}

//nolint:interfacer
func (c *ClusterResourceSet) addResourcesForVPC() error {

	c.vpc = c.newResource("VPC", &gfnec2.VPC{
		CidrBlock:          gfnt.NewString(c.spec.VPC.CIDR.String()),
		EnableDnsSupport:   gfnt.True(),
		EnableDnsHostnames: gfnt.True(),
	})

	if api.IsEnabled(c.spec.VPC.AutoAllocateIPv6) {
		c.newResource("AutoAllocatedCIDRv6", &gfnec2.VPCCidrBlock{
			VpcId:                       c.vpc,
			AmazonProvidedIpv6CidrBlock: gfnt.True(),
		})
	}

	c.subnets = make(map[api.SubnetTopology][]*gfnt.Value)

	refIG := c.newResource("InternetGateway", &gfnec2.InternetGateway{})
	vpcGA := "VPCGatewayAttachment"
	c.newResource(vpcGA, &gfnec2.VPCGatewayAttachment{
		InternetGatewayId: refIG,
		VpcId:             c.vpc,
	})

	refPublicRT := c.newResource("PublicRouteTable", &gfnec2.RouteTable{
		VpcId: c.vpc,
	})

	c.newResource("PublicSubnetRoute", &gfnec2.Route{
		RouteTableId:               refPublicRT,
		DestinationCidrBlock:       internetCIDR,
		GatewayId:                  refIG,
		AWSCloudFormationDependsOn: []string{vpcGA},
	})

	c.addSubnets(refPublicRT, api.SubnetTopologyPublic, c.spec.VPC.Subnets.Public)

	if err := c.addNATGateways(); err != nil {
		return err
	}

	c.addSubnets(nil, api.SubnetTopologyPrivate, c.spec.VPC.Subnets.Private)
	return nil
}

func (c *ClusterResourceSet) addNATGateways() error {

	switch *c.spec.VPC.NAT.Gateway {

	case api.ClusterHighlyAvailableNAT:
		c.haNAT()
	case api.ClusterSingleNAT:
		c.singleNAT()
	case api.ClusterDisableNAT:
		c.noNAT()
	default:
		// TODO validate this before starting to add resources
		return fmt.Errorf("%s is not a valid NAT gateway mode", *c.spec.VPC.NAT.Gateway)
	}
	return nil
}

func (c *ClusterResourceSet) importResourcesForVPC() {
	c.vpc = gfnt.NewString(c.spec.VPC.ID)
	c.subnets = make(map[api.SubnetTopology][]*gfnt.Value)
	for _, subnet := range c.spec.PrivateSubnetIDs() {
		c.subnets[api.SubnetTopologyPrivate] = append(c.subnets[api.SubnetTopologyPrivate], gfnt.NewString(subnet))
	}
	for _, subnet := range c.spec.PublicSubnetIDs() {
		c.subnets[api.SubnetTopologyPublic] = append(c.subnets[api.SubnetTopologyPublic], gfnt.NewString(subnet))
	}

}

func (c *ClusterResourceSet) addOutputsForVPC() {
	if c.spec.VPC == nil {
		c.spec.VPC = &api.ClusterVPC{}
	}
	c.rs.defineOutput(outputs.ClusterVPC, c.vpc, true, func(v string) error {
		c.spec.VPC.ID = v
		return nil
	})
	if c.spec.VPC.NAT != nil {
		c.rs.defineOutputWithoutCollector(outputs.ClusterFeatureNATMode, c.spec.VPC.NAT.Gateway, false)
	}
	if refs, ok := c.subnets[api.SubnetTopologyPrivate]; ok {
		c.rs.defineJoinedOutput(outputs.ClusterSubnetsPrivate, refs, true, func(v string) error {
			return vpc.ImportSubnetsFromList(c.provider, c.spec, api.SubnetTopologyPrivate, strings.Split(v, ","))
		})
	}
	if refs, ok := c.subnets[api.SubnetTopologyPublic]; ok {
		c.rs.defineJoinedOutput(outputs.ClusterSubnetsPublic, refs, true, func(v string) error {
			return vpc.ImportSubnetsFromList(c.provider, c.spec, api.SubnetTopologyPublic, strings.Split(v, ","))
		})
	}
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

func (c *ClusterResourceSet) addResourcesForSecurityGroups() {
	var refControlPlaneSG, refClusterSharedNodeSG *gfnt.Value

	if c.spec.VPC.SecurityGroup == "" {
		refControlPlaneSG = c.newResource(cfnControlPlaneSGResource, &gfnec2.SecurityGroup{
			GroupDescription: gfnt.NewString("Communication between the control plane and worker nodegroups"),
			VpcId:            c.vpc,
		})
	} else {
		refControlPlaneSG = gfnt.NewString(c.spec.VPC.SecurityGroup)
	}
	c.securityGroups = []*gfnt.Value{refControlPlaneSG} // only this one SG is passed to EKS API, nodes are isolated

	if c.spec.VPC.SharedNodeSecurityGroup == "" {
		refClusterSharedNodeSG = c.newResource(cfnSharedNodeSGResource, &gfnec2.SecurityGroup{
			GroupDescription: gfnt.NewString("Communication between all nodes in the cluster"),
			VpcId:            c.vpc,
		})
		c.newResource("IngressInterNodeGroupSG", &gfnec2.SecurityGroupIngress{
			GroupId:               refClusterSharedNodeSG,
			SourceSecurityGroupId: refClusterSharedNodeSG,
			Description:           gfnt.NewString("Allow nodes to communicate with each other (all ports)"),
			IpProtocol:            gfnt.NewString("-1"),
			FromPort:              sgPortZero,
			ToPort:                sgMaxNodePort,
		})
		if c.supportsManagedNodes {
			// To enable communication between both managed and unmanaged nodegroups, this allows ingress traffic from
			// the default cluster security group ID that EKS creates by default
			// EKS attaches this to Managed Nodegroups by default, but we need to handle this for unmanaged nodegroups
			c.newResource(cfnIngressClusterToNodeSGResource, &gfnec2.SecurityGroupIngress{
				GroupId:               refClusterSharedNodeSG,
				SourceSecurityGroupId: gfnt.MakeFnGetAttString("ControlPlane", outputs.ClusterDefaultSecurityGroup),
				Description:           gfnt.NewString("Allow managed and unmanaged nodes to communicate with each other (all ports)"),
				IpProtocol:            gfnt.NewString("-1"),
				FromPort:              sgPortZero,
				ToPort:                sgMaxNodePort,
			})
			c.newResource("IngressNodeToDefaultClusterSG", &gfnec2.SecurityGroupIngress{
				GroupId:               gfnt.MakeFnGetAttString("ControlPlane", outputs.ClusterDefaultSecurityGroup),
				SourceSecurityGroupId: refClusterSharedNodeSG,
				Description:           gfnt.NewString("Allow unmanaged nodes to communicate with control plane (all ports)"),
				IpProtocol:            gfnt.NewString("-1"),
				FromPort:              sgPortZero,
				ToPort:                sgMaxNodePort,
			})
		}
	} else {
		refClusterSharedNodeSG = gfnt.NewString(c.spec.VPC.SharedNodeSecurityGroup)
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
}

func (n *NodeGroupResourceSet) addResourcesForSecurityGroups() {
	for _, id := range n.spec.SecurityGroups.AttachIDs {
		n.securityGroups = append(n.securityGroups, gfnt.NewString(id))
	}

	if api.IsEnabled(n.spec.SecurityGroups.WithShared) {
		refClusterSharedNodeSG := makeImportValue(n.clusterStackName, outputs.ClusterSharedNodeSecurityGroup)
		n.securityGroups = append(n.securityGroups, refClusterSharedNodeSG)
	}

	if api.IsDisabled(n.spec.SecurityGroups.WithLocal) {
		return
	}

	desc := "worker nodes in group " + n.nodeGroupName

	allInternalIPv4 := gfnt.NewString(n.clusterSpec.VPC.CIDR.String())

	refControlPlaneSG := makeImportValue(n.clusterStackName, outputs.ClusterSecurityGroup)

	refNodeGroupLocalSG := n.newResource("SG", &gfnec2.SecurityGroup{
		VpcId:            makeImportValue(n.clusterStackName, outputs.ClusterVPC),
		GroupDescription: gfnt.NewString("Communication between the control plane and " + desc),
		Tags: []gfncfn.Tag{{
			Key:   gfnt.NewString("kubernetes.io/cluster/" + n.clusterSpec.Metadata.Name),
			Value: gfnt.NewString("owned"),
		}},
	})

	n.securityGroups = append(n.securityGroups, refNodeGroupLocalSG)

	n.newResource("IngressInterCluster", &gfnec2.SecurityGroupIngress{
		GroupId:               refNodeGroupLocalSG,
		SourceSecurityGroupId: refControlPlaneSG,
		Description:           gfnt.NewString("Allow " + desc + " to communicate with control plane (kubelet and workload TCP ports)"),
		IpProtocol:            sgProtoTCP,
		FromPort:              sgMinNodePort,
		ToPort:                sgMaxNodePort,
	})
	n.newResource("EgressInterCluster", &gfnec2.SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                gfnt.NewString("Allow control plane to communicate with " + desc + " (kubelet and workload TCP ports)"),
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgMinNodePort,
		ToPort:                     sgMaxNodePort,
	})
	n.newResource("IngressInterClusterAPI", &gfnec2.SecurityGroupIngress{
		GroupId:               refNodeGroupLocalSG,
		SourceSecurityGroupId: refControlPlaneSG,
		Description:           gfnt.NewString("Allow " + desc + " to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)"),
		IpProtocol:            sgProtoTCP,
		FromPort:              sgPortHTTPS,
		ToPort:                sgPortHTTPS,
	})
	n.newResource("EgressInterClusterAPI", &gfnec2.SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                gfnt.NewString("Allow control plane to communicate with " + desc + " (workloads using HTTPS port, commonly used with extension API servers)"),
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgPortHTTPS,
		ToPort:                     sgPortHTTPS,
	})
	n.newResource("IngressInterClusterCP", &gfnec2.SecurityGroupIngress{
		GroupId:               refControlPlaneSG,
		SourceSecurityGroupId: refNodeGroupLocalSG,
		Description:           gfnt.NewString("Allow control plane to receive API requests from " + desc),
		IpProtocol:            sgProtoTCP,
		FromPort:              sgPortHTTPS,
		ToPort:                sgPortHTTPS,
	})
	if *n.spec.SSH.Allow {
		if n.spec.PrivateNetworking {
			n.newResource("SSHIPv4", &gfnec2.SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIp:      allInternalIPv4,
				Description: gfnt.NewString("Allow SSH access to " + desc + " (private, only inside VPC)"),
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
		} else {
			n.newResource("SSHIPv4", &gfnec2.SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIp:      sgSourceAnywhereIPv4,
				Description: gfnt.NewString("Allow SSH access to " + desc),
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
			n.newResource("SSHIPv6", &gfnec2.SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIpv6:    sgSourceAnywhereIPv6,
				Description: gfnt.NewString("Allow SSH access to " + desc),
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
		}
	}
}

func (c *ClusterResourceSet) haNAT() {

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		// Allocate an EIP
		c.newResource("NATIP"+alphanumericUpperAZ, &gfnec2.EIP{
			Domain: gfnt.NewString("vpc"),
		})
		// Allocate a NAT gateway in the public subnet
		refNG := c.newResource("NATGateway"+alphanumericUpperAZ, &gfnec2.NatGateway{
			AllocationId: gfnt.MakeFnGetAttString("NATIP"+alphanumericUpperAZ, "AllocationId"),
			SubnetId:     gfnt.MakeRef("SubnetPublic" + alphanumericUpperAZ),
		})

		// Allocate a routing table for the private subnet
		refRT := c.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfnec2.RouteTable{
			VpcId: c.vpc,
		})
		// Create a route that sends Internet traffic through the NAT gateway
		c.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfnec2.Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		// Associate the routing table with the subnet
		c.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}

}

func (c *ClusterResourceSet) singleNAT() {

	sortedAZs := c.spec.AvailabilityZones
	firstUpperAZ := strings.ToUpper(strings.Join(strings.Split(sortedAZs[0], "-"), ""))

	c.newResource("NATIP", &gfnec2.EIP{
		Domain: gfnt.NewString("vpc"),
	})
	refNG := c.newResource("NATGateway", &gfnec2.NatGateway{
		AllocationId: gfnt.MakeFnGetAttString("NATIP", "AllocationId"),
		SubnetId:     gfnt.MakeRef("SubnetPublic" + firstUpperAZ),
	})

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := c.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfnec2.RouteTable{
			VpcId: c.vpc,
		})

		c.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfnec2.Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		c.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func (c *ClusterResourceSet) noNAT() {

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := c.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfnec2.RouteTable{
			VpcId: c.vpc,
		})
		c.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}
