package builder

import (
	"strings"

	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/weaveworks/eksctl/pkg/eks/api"
)

func (c *ClusterResourceSet) addSubnets(refRT *gfn.Value, topology api.SubnetTopology) {
	for az, subnet := range c.spec.VPC.Subnets[topology] {
		alias := string(topology) + strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))
		subnet := &gfn.AWSEC2Subnet{
			AvailabilityZone: gfn.NewString(az),
			CidrBlock:        gfn.NewString(subnet.CIDR.String()),
			VpcId:            c.vpc,
		}
		if topology == api.SubnetTopologyPrivate {
			subnet.Tags = []gfn.Tag{{
				Key:   gfn.NewString("kubernetes.io/role/internal-elb"),
				Value: gfn.NewString("1"),
			}}
		}
		refSubnet := c.newResource("Subnet"+alias, subnet)
		c.newResource("RouteTableAssociation"+alias, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     refSubnet,
			RouteTableId: refRT,
		})
		c.subnets[topology] = append(c.subnets[topology], refSubnet)
	}
}

//nolint:interfacer
func (c *ClusterResourceSet) addResourcesForVPC() {
	internetCIDR := gfn.NewString("0.0.0.0/0")

	c.vpc = c.newResource("VPC", &gfn.AWSEC2VPC{
		CidrBlock:          gfn.NewString(c.spec.VPC.CIDR.String()),
		EnableDnsSupport:   gfn.True(),
		EnableDnsHostnames: gfn.True(),
	})

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

	c.addSubnets(refPublicRT, api.SubnetTopologyPublic)

	c.newResource("NATIP", &gfn.AWSEC2EIP{
		Domain: gfn.NewString("vpc"),
	})
	refNG := c.newResource("NATGateway", &gfn.AWSEC2NatGateway{
		AllocationId: gfn.MakeFnGetAttString("NATIP.AllocationId"),
		// A multi-AZ NAT Gateway is possible, but it's not very
		// clear from the docs how to achieve it
		SubnetId: c.subnets[api.SubnetTopologyPublic][0],
	})

	refPrivateRT := c.newResource("PrivateRouteTable", &gfn.AWSEC2RouteTable{
		VpcId: c.vpc,
	})

	c.newResource("PrivateSubnetRoute", &gfn.AWSEC2Route{
		RouteTableId:         refPrivateRT,
		DestinationCidrBlock: internetCIDR,
		NatGatewayId:         refNG,
	})

	c.addSubnets(refPrivateRT, api.SubnetTopologyPrivate)
}

func (c *ClusterResourceSet) importResourcesForVPC() {
	c.vpc = gfn.NewString(c.spec.VPC.ID)
	c.subnets = make(map[api.SubnetTopology][]*gfn.Value)
	for topology := range c.spec.VPC.Subnets {
		for _, subnet := range c.spec.SubnetIDs(topology) {
			c.subnets[topology] = append(c.subnets[topology], gfn.NewString(subnet))
		}
	}
}

func (c *ClusterResourceSet) addOutputsForVPC() {
	c.rs.newOutput(cfnOutputClusterVPC, c.vpc, true)
	for topology := range c.spec.VPC.Subnets {
		c.rs.newJoinedOutput(cfnOutputClusterSubnets+string(topology), c.subnets[topology], true)
	}
}

func (c *ClusterResourceSet) addResourcesForSecurityGroups() {
	refSG := c.newResource("ControlPlaneSecurityGroup", &gfn.AWSEC2SecurityGroup{
		GroupDescription: gfn.NewString("Communication between the control plane and worker node groups"),
		VpcId:            c.vpc,
	})
	c.securityGroups = []*gfn.Value{refSG}
	c.rs.newJoinedOutput(cfnOutputClusterSecurityGroup, c.securityGroups, true)
}

func (n *NodeGroupResourceSet) addResourcesForSecurityGroups() {
	desc := "worker nodes in group " + n.nodeGroupName

	tcp := gfn.NewString("tcp")
	anywhereIPv4 := gfn.NewString("0.0.0.0/0")
	anywhereIPv6 := gfn.NewString("::/0")
	var (
		apiPort = gfn.NewInteger(443)
		sshPort = gfn.NewInteger(22)

		portZero    = gfn.NewInteger(0)
		nodeMinPort = gfn.NewInteger(1025)
		nodeMaxPort = gfn.NewInteger(65535)
	)

	refCP := makeImportValue(n.clusterStackName, cfnOutputClusterSecurityGroup)
	refSG := n.newResource("SG", &gfn.AWSEC2SecurityGroup{
		VpcId:            makeImportValue(n.clusterStackName, cfnOutputClusterVPC),
		GroupDescription: gfn.NewString("Communication between the control plane and " + desc),
		Tags: []gfn.Tag{{
			Key:   gfn.NewString("kubernetes.io/cluster/" + n.clusterSpec.ClusterName),
			Value: gfn.NewString("owned"),
		}},
	})
	n.securityGroups = []*gfn.Value{refSG}

	n.newResource("IngressInterSG", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refSG,
		SourceSecurityGroupId: refSG,
		Description:           gfn.NewString("Allow " + desc + " to communicate with each other (all ports)"),
		IpProtocol:            gfn.NewString("-1"),
		FromPort:              portZero,
		ToPort:                nodeMaxPort,
	})
	n.newResource("IngressInterCluster", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refSG,
		SourceSecurityGroupId: refCP,
		Description:           gfn.NewString("Allow " + desc + " to communicate with control plane (kubelet and workload TCP ports)"),
		IpProtocol:            tcp,
		FromPort:              nodeMinPort,
		ToPort:                nodeMaxPort,
	})
	n.newResource("EgressInterCluster", &gfn.AWSEC2SecurityGroupEgress{
		GroupId:                    refCP,
		DestinationSecurityGroupId: refSG,
		Description:                gfn.NewString("Allow control plane to communicate with " + desc + " (kubelet and workload TCP ports)"),
		IpProtocol:                 tcp,
		FromPort:                   nodeMinPort,
		ToPort:                     nodeMaxPort,
	})
	n.newResource("IngressInterClusterAPI", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refSG,
		SourceSecurityGroupId: refCP,
		Description:           gfn.NewString("Allow " + desc + " to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)"),
		IpProtocol:            tcp,
		FromPort:              apiPort,
		ToPort:                apiPort,
	})
	n.newResource("EgressInterClusterAPI", &gfn.AWSEC2SecurityGroupEgress{
		GroupId:                    refCP,
		DestinationSecurityGroupId: refSG,
		Description:                gfn.NewString("Allow control plane to communicate with " + desc + " (workloads using HTTPS port, commonly used with extension API servers)"),
		IpProtocol:                 tcp,
		FromPort:                   apiPort,
		ToPort:                     apiPort,
	})
	n.newResource("IngressInterClusterCP", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refCP,
		SourceSecurityGroupId: refSG,
		Description:           gfn.NewString("Allow control plane to receive API requests from " + desc),
		IpProtocol:            tcp,
		FromPort:              apiPort,
		ToPort:                apiPort,
	})
	if n.spec.AllowSSH {
		if n.spec.PrivateNetworking {
			n.newResource("SSHIPv4", &gfn.AWSEC2SecurityGroupIngress{
				GroupId:     refSG,
				CidrIp:      gfn.NewString(n.clusterSpec.VPC.CIDR.String()),
				Description: gfn.NewString("Allow SSH access to " + desc + " (private, only inside VPC)"),
				IpProtocol:  tcp,
				FromPort:    sshPort,
				ToPort:      sshPort,
			})
		} else {
			n.newResource("SSHIPv4", &gfn.AWSEC2SecurityGroupIngress{
				GroupId:     refSG,
				CidrIp:      anywhereIPv4,
				Description: gfn.NewString("Allow SSH access to " + desc),
				IpProtocol:  tcp,
				FromPort:    sshPort,
				ToPort:      sshPort,
			})
			n.newResource("SSHIPv6", &gfn.AWSEC2SecurityGroupIngress{
				GroupId:     refSG,
				CidrIpv6:    anywhereIPv6,
				Description: gfn.NewString("Allow SSH access to " + desc),
				IpProtocol:  tcp,
				FromPort:    sshPort,
				ToPort:      sshPort,
			})
		}
	}
}
