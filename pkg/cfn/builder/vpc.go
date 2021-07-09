package builder

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
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

// A VPCResourceSet builds the resources required for the specified VPC
type VPCResourceSet struct {
	rs            *resourceSet
	clusterConfig *api.ClusterConfig
	ec2API        ec2iface.EC2API

	vpcResource *VPCResource
}

// VPCResource represents a VPC resource
type VPCResource struct {
	VPC           *gfnt.Value
	SubnetDetails *subnetDetails
}

type SubnetResource struct {
	Subnet           *gfnt.Value
	RouteTable       *gfnt.Value
	AvailabilityZone string
}

type subnetDetails struct {
	Private []SubnetResource
	Public  []SubnetResource
}

// NewVPCResourceSet creates and returns a new VPCResourceSet
func NewVPCResourceSet(rs *resourceSet, clusterConfig *api.ClusterConfig, ec2API ec2iface.EC2API) *VPCResourceSet {
	var vpcRef *gfnt.Value
	if clusterConfig.VPC.ID == "" {
		vpcRef = rs.newResource("VPC", &gfnec2.VPC{
			CidrBlock:          gfnt.NewString(clusterConfig.VPC.CIDR.String()),
			EnableDnsSupport:   gfnt.True(),
			EnableDnsHostnames: gfnt.True(),
		})
	} else {
		vpcRef = gfnt.NewString(clusterConfig.VPC.ID)
	}

	return &VPCResourceSet{
		rs:            rs,
		clusterConfig: clusterConfig,
		ec2API:        ec2API,

		vpcResource: &VPCResource{
			VPC:           vpcRef,
			SubnetDetails: &subnetDetails{},
		},
	}
}

// AddResources adds all required resources
func (v *VPCResourceSet) AddResources() (*VPCResource, error) {
	vpc := v.clusterConfig.VPC
	if vpc.ID != "" { // custom VPC has been set
		if err := v.importResources(); err != nil {
			return nil, errors.Wrap(err, "error importing VPC resources")
		}
		return v.vpcResource, nil
	}

	if api.IsEnabled(vpc.AutoAllocateIPv6) {
		v.rs.newResource("AutoAllocatedCIDRv6", &gfnec2.VPCCidrBlock{
			VpcId:                       v.vpcResource.VPC,
			AmazonProvidedIpv6CidrBlock: gfnt.True(),
		})
	}

	if v.isFullyPrivate() {
		v.noNAT()
		v.vpcResource.SubnetDetails.Private = v.addSubnets(nil, api.SubnetTopologyPrivate, vpc.Subnets.Private)
		return v.vpcResource, nil
	}

	refIG := v.rs.newResource("InternetGateway", &gfnec2.InternetGateway{})
	vpcGA := "VPCGatewayAttachment"
	v.rs.newResource(vpcGA, &gfnec2.VPCGatewayAttachment{
		InternetGatewayId: refIG,
		VpcId:             v.vpcResource.VPC,
	})

	refPublicRT := v.rs.newResource("PublicRouteTable", &gfnec2.RouteTable{
		VpcId: v.vpcResource.VPC,
	})

	v.rs.newResource("PublicSubnetRoute", &gfnec2.Route{
		RouteTableId:               refPublicRT,
		DestinationCidrBlock:       internetCIDR,
		GatewayId:                  refIG,
		AWSCloudFormationDependsOn: []string{vpcGA},
	})

	v.vpcResource.SubnetDetails.Public = v.addSubnets(refPublicRT, api.SubnetTopologyPublic, vpc.Subnets.Public)

	if err := v.addNATGateways(); err != nil {
		return nil, err
	}

	v.vpcResource.SubnetDetails.Private = v.addSubnets(nil, api.SubnetTopologyPrivate, vpc.Subnets.Private)
	return v.vpcResource, nil
}

func (s *subnetDetails) PublicSubnetRefs() []*gfnt.Value {
	var subnetRefs []*gfnt.Value
	for _, subnetAZ := range s.Public {
		subnetRefs = append(subnetRefs, subnetAZ.Subnet)
	}
	return subnetRefs
}

func (s *subnetDetails) PrivateSubnetRefs() []*gfnt.Value {
	var subnetRefs []*gfnt.Value
	for _, subnetAZ := range s.Private {
		subnetRefs = append(subnetRefs, subnetAZ.Subnet)
	}
	return subnetRefs
}

// AddOutputs adds VPC resource outputs
func (v *VPCResourceSet) AddOutputs() {
	v.rs.defineOutput(outputs.ClusterVPC, v.vpcResource.VPC, true, func(val string) error {
		v.clusterConfig.VPC.ID = val
		return nil
	})
	if v.clusterConfig.VPC.NAT != nil {
		v.rs.defineOutputWithoutCollector(outputs.ClusterFeatureNATMode, v.clusterConfig.VPC.NAT.Gateway, false)
	}

	addSubnetOutput := func(subnetRefs []*gfnt.Value, topology api.SubnetTopology, outputName string) {
		v.rs.defineJoinedOutput(outputName, subnetRefs, true, func(value string) error {
			return vpc.ImportSubnetsFromIDList(v.ec2API, v.clusterConfig, topology, strings.Split(value, ","))
		})
	}

	if subnetAZs := v.vpcResource.SubnetDetails.PrivateSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, api.SubnetTopologyPrivate, outputs.ClusterSubnetsPrivate)
	}

	if subnetAZs := v.vpcResource.SubnetDetails.PublicSubnetRefs(); len(subnetAZs) > 0 {
		addSubnetOutput(subnetAZs, api.SubnetTopologyPublic, outputs.ClusterSubnetsPublic)
	}

	if v.isFullyPrivate() {
		v.rs.defineOutputWithoutCollector(outputs.ClusterFullyPrivate, true, true)
	}
}

// RenderJSON returns the rendered JSON
func (v *VPCResourceSet) RenderJSON() ([]byte, error) {
	return v.rs.renderJSON()
}

func (v *VPCResourceSet) addSubnets(refRT *gfnt.Value, topology api.SubnetTopology, subnets map[string]api.AZSubnetSpec) []SubnetResource {
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

	var subnetResources []SubnetResource

	for name, subnet := range subnets {
		az := subnet.AZ
		nameAlias := strings.ToUpper(strings.Join(strings.Split(name, "-"), ""))
		subnet := &gfnec2.Subnet{
			AvailabilityZone: gfnt.NewString(az),
			CidrBlock:        gfnt.NewString(subnet.CIDR.String()),
			VpcId:            v.vpcResource.VPC,
		}

		switch topology {
		case api.SubnetTopologyPrivate:
			// Choose the appropriate route table for private subnets
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

		if api.IsEnabled(v.clusterConfig.VPC.AutoAllocateIPv6) {
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
			v.rs.newResource(subnetAlias+"CIDRv6", &gfnec2.SubnetCidrBlock{
				SubnetId:      refSubnet,
				Ipv6CidrBlock: gfnt.MakeFnSelect(gfnt.NewInteger(subnetIndexForIPv6), refSubnetSlices),
			})
			subnetIndexForIPv6++
		}

		subnetResources = append(subnetResources, SubnetResource{
			AvailabilityZone: az,
			RouteTable:       refRT,
			Subnet:           refSubnet,
		})
	}
	return subnetResources
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
	if subnets := v.clusterConfig.VPC.Subnets.Private; subnets != nil {
		var (
			subnetRoutes map[string]string
			err          error
		)
		if v.isFullyPrivate() {
			subnetRoutes, err = importRouteTables(v.ec2API, v.clusterConfig.VPC.Subnets.Private)
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

func makeSubnetResources(subnets map[string]api.AZSubnetSpec, subnetRoutes map[string]string) ([]SubnetResource, error) {
	subnetResources := make([]SubnetResource, len(subnets))
	i := 0
	for _, network := range subnets {
		az := network.AZ
		sr := SubnetResource{
			AvailabilityZone: az,
			Subnet:           gfnt.NewString(network.ID),
		}

		if subnetRoutes != nil {
			rt, ok := subnetRoutes[network.ID]
			if !ok {
				return nil, errors.Errorf("failed to find an explicit route table associated with subnet %q; "+
					"eksctl does not modify the main route table if a subnet is not associated with an explicit route table", network.ID)
			}
			sr.RouteTable = gfnt.NewString(rt)
		}
		subnetResources[i] = sr
		i++
	}
	return subnetResources, nil
}

func importRouteTables(ec2API ec2iface.EC2API, subnets map[string]api.AZSubnetSpec) (map[string]string, error) {
	var subnetIDs []string
	for _, subnet := range subnets {
		subnetIDs = append(subnetIDs, subnet.ID)
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
			if rta.Main != nil && *rta.Main {
				return nil, errors.New("subnets must be associated with a non-main route table; eksctl does not modify the main route table")
			}
			subnetRoutes[*rta.SubnetId] = *rt.RouteTableId
		}
	}
	return subnetRoutes, nil
}

func (v *VPCResourceSet) isFullyPrivate() bool {
	return v.clusterConfig.PrivateCluster.Enabled
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

// TODO move this
func (c *ClusterResourceSet) addResourcesForSecurityGroups(vpcResource *VPCResource) *clusterSecurityGroup {
	var refControlPlaneSG, refClusterSharedNodeSG *gfnt.Value

	if c.spec.VPC.SecurityGroup == "" {
		refControlPlaneSG = c.newResource(cfnControlPlaneSGResource, &gfnec2.SecurityGroup{
			GroupDescription: gfnt.NewString("Communication between the control plane and worker nodegroups"),
			VpcId:            vpcResource.VPC,
		})

		if len(c.spec.VPC.ExtraCIDRs) > 0 {
			for i, cidr := range c.spec.VPC.ExtraCIDRs {
				c.newResource(fmt.Sprintf("IngressControlPlaneExtraCIDR%d", i), &gfnec2.SecurityGroupIngress{
					GroupId:     refControlPlaneSG,
					CidrIp:      gfnt.NewString(cidr),
					Description: gfnt.NewString(fmt.Sprintf("Allow Extra CIDR %d (%s) to communicate to controlplane", i, cidr)),
					IpProtocol:  gfnt.NewString("tcp"),
					FromPort:    sgPortHTTPS,
					ToPort:      sgPortHTTPS,
				})
			}
		}
	} else {
		refControlPlaneSG = gfnt.NewString(c.spec.VPC.SecurityGroup)
	}
	c.securityGroups = []*gfnt.Value{refControlPlaneSG} // only this one SG is passed to EKS API, nodes are isolated

	if c.spec.VPC.SharedNodeSecurityGroup == "" {
		refClusterSharedNodeSG = c.newResource(cfnSharedNodeSGResource, &gfnec2.SecurityGroup{
			GroupDescription: gfnt.NewString("Communication between all nodes in the cluster"),
			VpcId:            vpcResource.VPC,
		})
		c.newResource("IngressInterNodeGroupSG", &gfnec2.SecurityGroupIngress{
			GroupId:               refClusterSharedNodeSG,
			SourceSecurityGroupId: refClusterSharedNodeSG,
			Description:           gfnt.NewString("Allow nodes to communicate with each other (all ports)"),
			IpProtocol:            gfnt.NewString("-1"),
			FromPort:              sgPortZero,
			ToPort:                sgMaxNodePort,
		})
	} else {
		refClusterSharedNodeSG = gfnt.NewString(c.spec.VPC.SharedNodeSecurityGroup)
	}

	if c.supportsManagedNodes && api.IsEnabled(c.spec.VPC.ManageSharedNodeSecurityGroupRules) {
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

// TODO move this
func (rs *resourceSet) addEFASecurityGroup(vpcID *gfnt.Value, clusterName, desc string) *gfnt.Value {
	efaSG := rs.newResource("EFASG", &gfnec2.SecurityGroup{
		VpcId:            vpcID,
		GroupDescription: gfnt.NewString("EFA-enabled security group"),
		Tags: []gfncfn.Tag{{
			Key:   gfnt.NewString("kubernetes.io/cluster/" + clusterName),
			Value: gfnt.NewString("owned"),
		}},
	})
	rs.newResource("EFAIngressSelf", &gfnec2.SecurityGroupIngress{
		GroupId:               efaSG,
		SourceSecurityGroupId: efaSG,
		Description:           gfnt.NewString("Allow " + desc + " to communicate to itself (EFA-enabled)"),
		IpProtocol:            gfnt.NewString("-1"),
	})
	rs.newResource("EFAEgressSelf", &gfnec2.SecurityGroupEgress{
		GroupId:                    efaSG,
		DestinationSecurityGroupId: efaSG,
		Description:                gfnt.NewString("Allow " + desc + " to communicate to itself (EFA-enabled)"),
		IpProtocol:                 gfnt.NewString("-1"),
	})

	return efaSG
}

// TODO move this
func (n *NodeGroupResourceSet) addResourcesForSecurityGroups() {
	for _, id := range n.spec.SecurityGroups.AttachIDs {
		n.securityGroups = append(n.securityGroups, gfnt.NewString(id))
	}

	if api.IsEnabled(n.spec.SecurityGroups.WithShared) {
		n.securityGroups = append(n.securityGroups, n.vpcImporter.SharedNodeSecurityGroup())
	}

	if api.IsDisabled(n.spec.SecurityGroups.WithLocal) {
		return
	}

	desc := "worker nodes in group " + n.spec.Name
	vpcID := n.vpcImporter.VPC()
	refControlPlaneSG := n.vpcImporter.ControlPlaneSecurityGroup()

	refNodeGroupLocalSG := n.newResource("SG", &gfnec2.SecurityGroup{
		VpcId:            vpcID,
		GroupDescription: gfnt.NewString("Communication between the control plane and " + desc),
		Tags: []gfncfn.Tag{{
			Key:   gfnt.NewString("kubernetes.io/cluster/" + n.clusterSpec.Metadata.Name),
			Value: gfnt.NewString("owned"),
		}},
		SecurityGroupIngress: makeNodeIngressRules(n.spec.NodeGroupBase, refControlPlaneSG, n.clusterSpec.VPC.CIDR.String(), desc),
	})

	n.securityGroups = append(n.securityGroups, refNodeGroupLocalSG)

	if api.IsEnabled(n.spec.EFAEnabled) {
		efaSG := n.rs.addEFASecurityGroup(vpcID, n.clusterSpec.Metadata.Name, desc)
		n.securityGroups = append(n.securityGroups, efaSG)
	}

	n.newResource("EgressInterCluster", &gfnec2.SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                gfnt.NewString("Allow control plane to communicate with " + desc + " (kubelet and workload TCP ports)"),
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgMinNodePort,
		ToPort:                     sgMaxNodePort,
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
}

func makeNodeIngressRules(ng *api.NodeGroupBase, controlPlaneSG *gfnt.Value, vpcCIDR, description string) []gfnec2.SecurityGroup_Ingress {
	ingressRules := []gfnec2.SecurityGroup_Ingress{
		{
			SourceSecurityGroupId: controlPlaneSG,
			Description:           gfnt.NewString(fmt.Sprintf("[IngressInterCluster] Allow %s to communicate with control plane (kubelet and workload TCP ports)", description)),
			IpProtocol:            sgProtoTCP,
			FromPort:              sgMinNodePort,
			ToPort:                sgMaxNodePort,
		},
		{
			SourceSecurityGroupId: controlPlaneSG,
			Description:           gfnt.NewString(fmt.Sprintf("[IngressInterClusterAPI] Allow %s to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)", description)),
			IpProtocol:            sgProtoTCP,
			FromPort:              sgPortHTTPS,
			ToPort:                sgPortHTTPS,
		},
	}

	return append(ingressRules, makeSSHIngressRules(ng, vpcCIDR, description)...)
}

func (v *VPCResourceSet) haNAT() {
	for _, az := range v.clusterConfig.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

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
			VpcId: v.vpcResource.VPC,
		})
		// Create a route that sends Internet traffic through the NAT gateway
		v.rs.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfnec2.Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		// Associate the routing table with the subnet
		v.rs.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func (v *VPCResourceSet) singleNAT() {
	sortedAZs := v.clusterConfig.AvailabilityZones
	firstUpperAZ := strings.ToUpper(strings.Join(strings.Split(sortedAZs[0], "-"), ""))

	v.rs.newResource("NATIP", &gfnec2.EIP{
		Domain: gfnt.NewString("vpc"),
	})
	refNG := v.rs.newResource("NATGateway", &gfnec2.NatGateway{
		AllocationId: gfnt.MakeFnGetAttString("NATIP", "AllocationId"),
		SubnetId:     gfnt.MakeRef("SubnetPublic" + firstUpperAZ),
	})

	for _, az := range v.clusterConfig.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := v.rs.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfnec2.RouteTable{
			VpcId: v.vpcResource.VPC,
		})

		v.rs.newResource("NATPrivateSubnetRoute"+alphanumericUpperAZ, &gfnec2.Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		v.rs.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func (v *VPCResourceSet) noNAT() {
	for _, az := range v.clusterConfig.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := v.rs.newResource("PrivateRouteTable"+alphanumericUpperAZ, &gfnec2.RouteTable{
			VpcId: v.vpcResource.VPC,
		})
		v.rs.newResource("RouteTableAssociationPrivate"+alphanumericUpperAZ, &gfnec2.SubnetRouteTableAssociation{
			SubnetId:     gfnt.MakeRef("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}
