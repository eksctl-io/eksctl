package builder

import (
	"fmt"
	"strings"

	gfn "github.com/awslabs/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

var internetCIDR = gfn.NewString("0.0.0.0/0")

const (
	cfnControlPlaneSGResource = "ControlPlaneSecurityGroup"
	cfnSharedNodeSGResource   = "ClusterSharedNodeSecurityGroup"
)

func (c *ClusterResourceSet) addSubnets(refRT *gfn.Value, topology api.SubnetTopology, subnets map[string]api.Network) {
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
		subnet := &gfn.AWSEC2Subnet{
			AvailabilityZone: gfn.NewString(az),
			CidrBlock:        gfn.NewString(subnet.CIDR.String()),
			VpcId:            c.vpc,
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
		}
		refSubnet := c.newResource("Subnet"+alias, subnet)
		c.newResource("RouteTableAssociation"+alias, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     refSubnet,
			RouteTableId: refRT,
		})

		if api.IsEnabled(c.spec.VPC.AutoAllocateIPv6) {
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
			c.newResource(alias+"CIDRv6", &gfn.AWSEC2SubnetCidrBlock{
				SubnetId:      refSubnet,
				Ipv6CidrBlock: gfn.MakeFnSelect(subnetIndexForIPv6, refSubnetSlices),
			})
			subnetIndexForIPv6++
		}

		c.subnets[topology] = append(c.subnets[topology], refSubnet)
	}
}

//nolint:interfacer
func (c *ClusterResourceSet) addResourcesForVPC() error {

	c.vpc = c.newResource("VPC", &gfn.AWSEC2VPC{
		CidrBlock:          gfn.NewString(c.spec.VPC.CIDR.String()),
		EnableDnsSupport:   gfn.True(),
		EnableDnsHostnames: gfn.True(),
	})

	if api.IsEnabled(c.spec.VPC.AutoAllocateIPv6) {
		c.newResource("AutoAllocatedCIDRv6", &gfn.AWSEC2VPCCidrBlock{
			VpcId: c.vpc,
			AmazonProvidedIpv6CidrBlock: gfn.True(),
		})
	}

	c.subnets = make(map[api.SubnetTopology][]*gfn.Value)

	refIG := c.newResource("InternetGateway", &gfn.AWSEC2InternetGateway{})
	c.newResource("VPCGatewayAttachment", &gfn.AWSEC2VPCGatewayAttachment{
		InternetGatewayId: refIG,
		VpcId:             c.vpc,
	})

	refPublicRT := c.newResource("PublicRouteTable", &gfn.AWSEC2RouteTable{
		VpcId: c.vpc,
	})

	c.newResource("PublicSubnetRoute", &gfn.AWSEC2Route{
		RouteTableId:         refPublicRT,
		DestinationCidrBlock: internetCIDR,
		GatewayId:            refIG,
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
	c.vpc = gfn.NewString(c.spec.VPC.ID)
	c.subnets = make(map[api.SubnetTopology][]*gfn.Value)
	for _, subnet := range c.spec.PrivateSubnetIDs() {
		c.subnets[api.SubnetTopologyPrivate] = append(c.subnets[api.SubnetTopologyPrivate], gfn.NewString(subnet))
	}
	for _, subnet := range c.spec.PublicSubnetIDs() {
		c.subnets[api.SubnetTopologyPublic] = append(c.subnets[api.SubnetTopologyPublic], gfn.NewString(subnet))
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
	c.rs.defineOutputWithoutCollector(outputs.ClusterFeatureNATMode, c.spec.VPC.NAT.Gateway, false)
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
	sgProtoTCP           = gfn.NewString("tcp")
	sgSourceAnywhereIPv4 = gfn.NewString("0.0.0.0/0")
	sgSourceAnywhereIPv6 = gfn.NewString("::/0")

	sgPortZero    = gfn.NewInteger(0)
	sgMinNodePort = gfn.NewInteger(1025)
	sgMaxNodePort = gfn.NewInteger(65535)

	sgPortHTTPS = gfn.NewInteger(443)
	sgPortSSH   = gfn.NewInteger(22)
)

func (c *ClusterResourceSet) addResourcesForSecurityGroups() {
	var refControlPlaneSG, refClusterSharedNodeSG *gfn.Value

	if c.spec.VPC.SecurityGroup == "" {
		refControlPlaneSG = c.newResource(cfnControlPlaneSGResource, &gfn.AWSEC2SecurityGroup{
			GroupDescription: gfn.NewString("Communication between the control plane and worker nodegroups"),
			VpcId:            c.vpc,
		})
	} else {
		refControlPlaneSG = gfn.NewString(c.spec.VPC.SecurityGroup)
	}
	c.securityGroups = []*gfn.Value{refControlPlaneSG} // only this one SG is passed to EKS API, nodes are isolated

	if c.spec.VPC.SharedNodeSecurityGroup == "" {
		refClusterSharedNodeSG = c.newResource(cfnSharedNodeSGResource, &gfn.AWSEC2SecurityGroup{
			GroupDescription: gfn.NewString("Communication between all nodes in the cluster"),
			VpcId:            c.vpc,
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
			c.newResource("IngressDefaultClusterToNodeSG", &gfn.AWSEC2SecurityGroupIngress{
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

func (c *ClusterResourceSet) haNAT() {

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		// Allocate an EIP
		c.newResource("NATIP"+alphanumericUpperAZ, &gfn.AWSEC2EIP{
			Domain: gfn.NewString("vpc"),
		})
		// Allocate a NAT gateway in the public subnet
		refNG := c.newResource("NATGateway"+alphanumericUpperAZ, &gfn.AWSEC2NatGateway{
			AllocationId: gfn.MakeFnGetAttString("NATIP" + alphanumericUpperAZ + ".AllocationId"),
			SubnetId:     gfn.MakeRef("SubnetPublic" + alphanumericUpperAZ),
		})

		// Allocate a routing table for the private subnet
		refRT := c.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfn.AWSEC2RouteTable{
			VpcId: c.vpc,
		})
		// Create a route that sends Internet traffic through the NAT gateway
		c.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfn.AWSEC2Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		// Associate the routing table with the subnet
		c.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     gfn.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}

}

func (c *ClusterResourceSet) singleNAT() {

	sortedAZs := c.spec.AvailabilityZones
	firstUpperAZ := strings.ToUpper(strings.Join(strings.Split(sortedAZs[0], "-"), ""))

	c.newResource("NATIP", &gfn.AWSEC2EIP{
		Domain: gfn.NewString("vpc"),
	})
	refNG := c.newResource("NATGateway", &gfn.AWSEC2NatGateway{
		AllocationId: gfn.MakeFnGetAttString("NATIP.AllocationId"),
		SubnetId:     gfn.MakeRef("SubnetPublic" + firstUpperAZ),
	})

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := c.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfn.AWSEC2RouteTable{
			VpcId: c.vpc,
		})

		c.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfn.AWSEC2Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		c.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     gfn.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func (c *ClusterResourceSet) noNAT() {

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := c.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfn.AWSEC2RouteTable{
			VpcId: c.vpc,
		})
		c.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     gfn.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}
