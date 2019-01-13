package builder

import (
	"fmt"
	"net"
	"strings"

	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
	"k8s.io/kops/pkg/util/subnet"
)

func (c *ClusterResourceSet) addSubnet(CIDR *net.IPNet, az string, topology api.SubnetTopology) {
	alias := string(topology) + strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))
	subnet := &gfn.AWSEC2Subnet{
		AvailabilityZone: gfn.NewString(az),
		CidrBlock:        gfn.NewString(CIDR.String()),
		VpcId:            c.vpc,
	}
	if topology == api.SubnetTopologyPrivate {
		subnet.Tags = []gfn.Tag{{
			Key:   gfn.NewString("kubernetes.io/role/internal-elb"),
			Value: gfn.NewString("1"),
		}}
	}
	refSubnet := c.newResource("Subnet"+alias, subnet)
	c.subnets[topology] = append(c.subnets[topology], refSubnet)
}

//nolint:interfacer
func (c *ClusterResourceSet) addResourcesForVPC() {
	c.vpc = c.newResource("VPC", &gfn.AWSEC2VPC{
		CidrBlock:          gfn.NewString(c.spec.VPC.CIDR.String()),
		EnableDnsSupport:   gfn.True(),
		EnableDnsHostnames: gfn.True(),
	})

}

//nolint:interfacer
func (c *ClusterResourceSet) addResourcesForIGW() {
	c.igw = c.newResource("InternetGateway", &gfn.AWSEC2InternetGateway{})
	c.newResource("VPCGatewayAttachment", &gfn.AWSEC2VPCGatewayAttachment{
		InternetGatewayId: c.igw,
		VpcId:             c.vpc,
	})
}

//nolint:interfacer
func (c *ClusterResourceSet) addResourcesForRouting() {
	internetCIDR := gfn.NewString("0.0.0.0/0")
	routeTables := make(map[api.SubnetTopology]*gfn.Value)
	routeTables[api.SubnetTopologyPublic] = c.newResource("PublicRouteTable", &gfn.AWSEC2RouteTable{
		VpcId: c.vpc,
	})

	c.newResource("PublicSubnetRoute", &gfn.AWSEC2Route{
		RouteTableId:         routeTables[api.SubnetTopologyPublic],
		DestinationCidrBlock: internetCIDR,
		GatewayId:            c.igw,
	})
	c.newResource("NATIP", &gfn.AWSEC2EIP{
		Domain: gfn.NewString("vpc"),
	})
	refNG := c.newResource("NATGateway", &gfn.AWSEC2NatGateway{
		AllocationId: gfn.MakeFnGetAttString("NATIP.AllocationId"),
		// A multi-AZ NAT Gateway is possible, but it's not very
		// clear from the docs how to achieve it
		SubnetId: c.subnets[api.SubnetTopologyPublic][0],
	})

	routeTables[api.SubnetTopologyPrivate] = c.newResource("PrivateRouteTable", &gfn.AWSEC2RouteTable{
		VpcId: c.vpc,
	})

	c.newResource("PrivateSubnetRoute", &gfn.AWSEC2Route{
		RouteTableId:         routeTables[api.SubnetTopologyPrivate],
		DestinationCidrBlock: internetCIDR,
		NatGatewayId:         refNG,
	})
	for topology, subnets := range c.subnets {
		for i, subnet := range subnets {
			c.newResource(fmt.Sprintf("RouteTableAssociation%s%v", string(topology), i), &gfn.AWSEC2SubnetRouteTableAssociation{
				SubnetId:     subnet,
				RouteTableId: routeTables[topology],
			})
		}
	}
}

//nolint:interfacer
func (c *ClusterResourceSet) addResourcesForSubnets() error {
	var err error

	c.subnets = make(map[api.SubnetTopology][]*gfn.Value)
	prefix, _ := c.spec.VPC.CIDR.Mask.Size()
	if (prefix < 16) || (prefix > 24) {
		return fmt.Errorf("VPC CIDR prefix must be betwee /16 and /24")
	}
	zoneCIDRs, err := subnet.SplitInto8(&c.spec.VPC.CIDR.IPNet)
	if err != nil {
		return err
	}

	logger.Debug("VPC CIDR (%s) was divided into 8 subnets %v", c.spec.VPC.CIDR.String(), zoneCIDRs)

	zonesTotal := len(c.spec.AvailabilityZones)
	if 2*zonesTotal > len(zoneCIDRs) {
		return fmt.Errorf("insufficient number of subnets (have %d, but need %d) for %d availability zones", len(zoneCIDRs), 2*zonesTotal, zonesTotal)
	}

	for i, zone := range c.spec.AvailabilityZones {
		public := zoneCIDRs[i]
		private := zoneCIDRs[i+zonesTotal]
		c.addSubnet(public, zone, api.SubnetTopologyPublic)
		c.addSubnet(private, zone, api.SubnetTopologyPrivate)
		logger.Info("subnets for %s - public:%s private:%s", zone, public.String(), private.String())
	}

	return nil

}

func (c *ClusterResourceSet) importResourcesForVPC() {
	c.vpc = gfn.NewString(c.spec.VPC.ID)
}
func (c *ClusterResourceSet) importResourcesForIGW() {
	c.igw = gfn.NewString(c.spec.VPC.IGW.ID)
}
func (c *ClusterResourceSet) importResourcesForSubnets() {
	c.subnets = make(map[api.SubnetTopology][]*gfn.Value)
	for topology := range c.spec.VPC.Subnets {
		for _, subnet := range c.spec.SubnetIDs(topology) {
			c.subnets[topology] = append(c.subnets[topology], gfn.NewString(subnet))
		}
	}
}

func (c *ClusterResourceSet) addOutputsForVPC() {
	c.rs.newOutput(cfnOutputClusterVPC, c.vpc, true)
	for topology := range c.subnets {
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
	udp := gfn.NewString("udp")
	allInternalIPv4 := gfn.NewString(n.clusterSpec.VPC.CIDR.String())
	anywhereIPv4 := gfn.NewString("0.0.0.0/0")
	anywhereIPv6 := gfn.NewString("::/0")
	var (
		apiPort = gfn.NewInteger(443)
		dnsPort = gfn.NewInteger(53)
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
			Key:   gfn.NewString("kubernetes.io/cluster/" + n.clusterSpec.Metadata.Name),
			Value: gfn.NewString("owned"),
		}},
	})
	n.securityGroups = []*gfn.Value{refSG}

	if len(n.spec.SecurityGroups) > 0 {
		for _, arn := range n.spec.SecurityGroups {
			n.securityGroups = append(n.securityGroups, gfn.NewString(arn))
		}
	}

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
				CidrIp:      allInternalIPv4,
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
	n.newResource("DNSUDPIPv4", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:     refSG,
		CidrIp:      allInternalIPv4,
		Description: gfn.NewString("Allow DNS access to " + desc + " inside VPC"),
		IpProtocol:  udp,
		FromPort:    dnsPort,
		ToPort:      dnsPort,
	})
	n.newResource("DNSTCPIPv4", &gfn.AWSEC2SecurityGroupIngress{
		GroupId:     refSG,
		CidrIp:      allInternalIPv4,
		Description: gfn.NewString("Allow DNS access to " + desc + " inside VPC"),
		IpProtocol:  tcp,
		FromPort:    dnsPort,
		ToPort:      dnsPort,
	})
}
