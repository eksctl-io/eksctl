package builder

import (
	"fmt"
	"strings"

	gfn "github.com/awslabs/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

func (c *ClusterResourceSet) addSubnets(refRT *gfn.Value, topology api.SubnetTopology, subnets map[string]api.Network, kind string) {
	for az, subnet := range subnets {
		alias := string(kind) + string(topology) + strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))
		subnet := &gfn.AWSEC2Subnet{
			AvailabilityZone: gfn.NewString(az),
			CidrBlock:        gfn.NewString(subnet.CIDR.String()),
			VpcId:            c.vpc,
		}
		switch topology {
		case api.SubnetTopologyPrivate:
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
		if alias == "" {
			c.newResource("RouteTableAssociation"+alias, &gfn.AWSEC2SubnetRouteTableAssociation{
				SubnetId:     refSubnet,
				RouteTableId: refRT,
			})
		} else {
			c.newResource("RouteTableAssociation"+alias, &awsCloudFormationResource{
				Type: "AWS::EC2::SubnetRouteTableAssociation",
				Properties: map[string]interface{}{
					"SubnetId":     refSubnet,
					"RouteTableId": refRT,
				},
				DependsOn: []string{alias},
			})
		}
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

	for i, podCIDR := range c.spec.VPC.PodCIDRs {
		c.newResource(fmt.Sprintf("PodCIDR%d", i), &awsCloudFormationResource{
			Type: "AWS::EC2::VPCCidrBlock",
			Properties: map[string]interface{}{
				"CidrBlock": podCIDR.String(),
				"VpcId":     c.vpc,
			},
			DependsOn: []string{"VPC"},
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

	c.addSubnets(refPublicRT, api.SubnetTopologyPublic, c.spec.VPC.Subnets.Public, "")

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

	c.addSubnets(refPrivateRT, api.SubnetTopologyPrivate, c.spec.VPC.Subnets.Private, "")

	// TODO add specific name for pod subnets
	for i := range c.spec.VPC.PodCIDRs {

		c.addSubnets(refPrivateRT, api.SubnetTopologyPrivate, c.spec.VPC.PodSubnets[fmt.Sprintf("eksctlGroup%d", i)].Subnets.Private, fmt.Sprintf("PodCIDR%d", i))
	}
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
		refControlPlaneSG = c.newResource("ControlPlaneSecurityGroup", &gfn.AWSEC2SecurityGroup{
			GroupDescription: gfn.NewString("Communication between the control plane and worker nodegroups"),
			VpcId:            c.vpc,
		})
	} else {
		refControlPlaneSG = gfn.NewString(c.spec.VPC.SecurityGroup)
	}
	c.securityGroups = []*gfn.Value{refControlPlaneSG} // only this one SG is passed to EKS API, nodes are isolated

	if c.spec.VPC.SharedNodeSecurityGroup == "" {
		refClusterSharedNodeSG = c.newResource("ClusterSharedNodeSecurityGroup", &gfn.AWSEC2SecurityGroup{
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
