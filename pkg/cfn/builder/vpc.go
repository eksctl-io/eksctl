package builder

import (
	"net"
	"strings"

	gfn "github.com/awslabs/goformation/cloudformation"
)

const (
	cfnOutputClusterVPC           = "VPC"
	cfnOutputClusterSubnets       = "Subnets"
	cfnOutputClusterSecurityGroup = "SecurityGroup"
)

func (c *ClusterResourceSet) addResourcesForVPC(globalCIDR *net.IPNet, subnets map[string]*net.IPNet) {
	refVPC := c.newResource("VPC", &gfn.AWSEC2VPC{
		CidrBlock:          gfn.NewString(globalCIDR.String()),
		EnableDnsSupport:   gfn.True(),
		EnableDnsHostnames: gfn.True(),
	})

	refIG := c.newResource("InternetGateway", &gfn.AWSEC2InternetGateway{})
	c.newResource("VPCGatewayAttachment", &gfn.AWSEC2VPCGatewayAttachment{
		InternetGatewayId: refIG,
		VpcId:             refVPC,
	})

	refRT := c.newResource("RouteTable", &gfn.AWSEC2RouteTable{
		VpcId: refVPC,
	})

	c.newResource("PublicSubnetRoute", &gfn.AWSEC2Route{
		RouteTableId:         refRT,
		DestinationCidrBlock: gfn.NewString("0.0.0.0/0"),
		GatewayId:            refIG,
	})

	for az, subnet := range subnets {
		alias := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))
		refSubnet := c.newResource("Subnet"+alias, &gfn.AWSEC2Subnet{
			AvailabilityZone: gfn.NewString(az),
			CidrBlock:        gfn.NewString(subnet.String()),
			VpcId:            refVPC,
		})
		c.newResource("RouteTableAssociation"+alias, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     refSubnet,
			RouteTableId: refRT,
		})
		c.subnets = append(c.subnets, refSubnet)
	}

	refSG := c.newResource("ControlPlaneSecurityGroup", &gfn.AWSEC2SecurityGroup{
		GroupDescription: gfn.NewString("Communication between the control plane and worker node groups"),
		VpcId:            refVPC,
	})
	c.securityGroups = []*gfn.Value{refSG}

	c.rs.newOutput(cfnOutputClusterVPC, refVPC, true)
	c.rs.newJoinedOutput(cfnOutputClusterSecurityGroup, c.securityGroups, true)
	c.rs.newJoinedOutput(cfnOutputClusterSubnets, c.subnets, true)
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
			Key:   gfn.NewString("kubernetes.io/cluster/" + n.spec.ClusterName),
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
	if n.spec.NodeSSH {
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
