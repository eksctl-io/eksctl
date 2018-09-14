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

func (c *clusterResourceSet) addResourcesForVPC(globalCIDR *net.IPNet, subnets map[string]*net.IPNet) {
	refVPC := c.newResource("VPC", &gfn.AWSEC2VPC{
		CidrBlock:          globalCIDR.String(),
		EnableDnsSupport:   true,
		EnableDnsHostnames: true,
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
		DestinationCidrBlock: "0.0.0.0/0",
		GatewayId:            refIG,
	})

	for az, subnet := range subnets {
		alias := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))
		refSubnet := c.newResource("Subnet"+alias, &gfn.AWSEC2Subnet{
			AvailabilityZone: az,
			CidrBlock:        subnet.String(),
			VpcId:            refVPC,
		})
		c.newResource("RouteTableAssociation"+alias, &gfn.AWSEC2SubnetRouteTableAssociation{
			SubnetId:     refSubnet,
			RouteTableId: refRT,
		})
		c.subnets = append(c.subnets, refSubnet)
	}

	refSG := c.newResource("ControlPlaneSecurityGroup", &gfn.AWSEC2SecurityGroup{
		GroupDescription: "Communication between the control plane and worker node groups",
		VpcId:            refVPC,
	})
	c.securityGroups = []string{refSG}

	c.rs.newOutput(cfnOutputClusterVPC, refVPC, true)
	c.rs.newJoinedOutput(cfnOutputClusterSecurityGroup, c.securityGroups, true)
	c.rs.newJoinedOutput(cfnOutputClusterSubnets, c.subnets, true)
}

func (n *nodeGroupResourceSet) addResourcesForSecurityGroups() {
	desc := "worker nodes in group " + nodeGroupNameFmt

	tcp := "tcp"
	anywhereIPv4 := "0.0.0.0/0"
	anywhereIPv6 := "::/0"
	const (
		apiPort = 443
		sshPort = 22

		nodeMinPort = 1025
		nodeMaxPort = 65535
	)

	refCP := makeImportValue(ParamClusterStackName, cfnOutputClusterSecurityGroup)
	refSG := n.newResource("SG", &gfn.AWSEC2SecurityGroup{
		VpcId:            makeImportValue(ParamClusterStackName, cfnOutputClusterVPC),
		GroupDescription: gfn.Sub("Communication between the control plane and " + desc),
		Tags:             []gfn.Tag{clusterOwnedTag},
	})
	n.securityGroups = []string{refSG}

	n.newResource("IngressInterSG", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refSG,
		SourceSecurityGroupId: refSG,
		Description:           gfn.Sub("Allow " + desc + " to communicate with each other (all ports)"),
		IpProtocol:            "-1",
		FromPort:              0,
		ToPort:                nodeMaxPort,
	})
	n.newResource("IngressInterCluster", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refSG,
		SourceSecurityGroupId: refCP,
		Description:           gfn.Sub("Allow " + desc + " to communicate with control plane (kubelet and workload TCP ports)"),
		IpProtocol:            tcp,
		FromPort:              nodeMinPort,
		ToPort:                nodeMaxPort,
	})
	n.newResource("EgressInterCluster", &gfn.AWSEC2SecurityGroupEgress{
		GroupId:                    refCP,
		DestinationSecurityGroupId: refSG,
		Description:                gfn.Sub("Allow " + desc + " to communicate with control plane (kubelet and workload TCP ports)"),
		IpProtocol:                 tcp,
		FromPort:                   nodeMinPort,
		ToPort:                     nodeMaxPort,
	})
	n.newResource("IngressInterClusterCP", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:               refCP,
		SourceSecurityGroupId: refSG,
		Description:           gfn.Sub("Allow control plane to recieve API requests from " + desc),
		IpProtocol:            tcp,
		FromPort:              apiPort,
		ToPort:                apiPort,
	})
	if n.spec.NodeSSH {
		n.newResource("SSHIPv4", &gfn.AWSEC2SecurityGroupIngress{
			GroupId:     refSG,
			CidrIp:      anywhereIPv4,
			Description: gfn.Sub("Allow SSH access to " + desc),
			IpProtocol:  tcp,
			FromPort:    sshPort,
			ToPort:      sshPort,
		})
		n.newResource("SSHIPv6", &gfn.AWSEC2SecurityGroupIngress{
			GroupId:     refSG,
			CidrIpv6:    anywhereIPv6,
			Description: gfn.Sub("Allow SSH access to " + desc),
			IpProtocol:  tcp,
			FromPort:    sshPort,
			ToPort:      sshPort,
		})
	}
}
