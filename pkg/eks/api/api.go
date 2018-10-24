package api

import (
	"net"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

const (
	// AWSDebugLevel defines the LogLevel for AWS produced logs
	AWSDebugLevel = 5

	// EKSRegionUSWest2 represents the US West Region Oregon
	EKSRegionUSWest2 = "us-west-2"

	// EKSRegionUSEast1 represents the US East Region North Virgina
	EKSRegionUSEast1 = "us-east-1"

	// EKSRegionEUWest1 represents the EU West Region Ireland
	EKSRegionEUWest1 = "eu-west-1"

	// DefaultEKSRegion defines the default region, where to deploy the EKS cluster
	DefaultEKSRegion = EKSRegionUSWest2
)

// SupportedRegions are the regions where EKS is available
var SupportedRegions = []string{
	EKSRegionUSWest2,
	EKSRegionUSEast1,
	EKSRegionEUWest1,
}

// DefaultWaitTimeout defines the default wait timeout
var DefaultWaitTimeout = 20 * time.Minute

// DefaultNodeCount defines the default number of nodes to be created
var DefaultNodeCount = 2

// ClusterProvider provides an interface with the needed AWS APIs
type ClusterProvider interface {
	CloudFormation() cloudformationiface.CloudFormationAPI
	EKS() eksiface.EKSAPI
	EC2() ec2iface.EC2API
	STS() stsiface.STSAPI
}

// ClusterConfig is a simple config, to be replaced with Cluster API
type ClusterConfig struct {
	Region      string
	Profile     string
	Tags        map[string]string
	ClusterName string

	WaitTimeout time.Duration

	VPC ClusterVPC

	NodeGroups []*NodeGroup

	Endpoint                 string
	CertificateAuthorityData []byte
	ARN                      string

	ClusterStackName string

	AvailabilityZones []string

	Addons ClusterAddons
}

// NewClusterConfig create new config for a cluster;
// it doesn't include initial nodegroup, so user must
// call NewNodeGroup to create one
func NewClusterConfig() *ClusterConfig {
	return &ClusterConfig{}
}

// SetSubnets defines CIDRs for each of the subnets,
// it must be called after SetAvailabilityZones
func (c *ClusterConfig) SetSubnets() {
	_, c.VPC.CIDR, _ = net.ParseCIDR("192.168.0.0/16")

	c.VPC.Subnets = map[SubnetTopology]map[string]Network{
		SubnetTopologyPublic: map[string]Network{},
	}

	zoneCIDRs := []string{
		"192.168.64.0/18",
		"192.168.128.0/18",
		"192.168.192.0/18",
	}
	for i, zone := range c.AvailabilityZones {
		_, zoneCIDR, _ := net.ParseCIDR(zoneCIDRs[i])
		c.VPC.Subnets[SubnetTopologyPublic][zone] = Network{
			CIDR: zoneCIDR,
		}
	}
}

// NewNodeGroup crears new nodegroup inside cluster config,
// it returns pointer to the nodegroup for convenience
func (c *ClusterConfig) NewNodeGroup() *NodeGroup {
	ng := &NodeGroup{
		ID:             len(c.NodeGroups),
		SubnetTopology: SubnetTopologyPublic,
	}

	c.NodeGroups = append(c.NodeGroups, ng)

	return ng
}

// NodeGroup holds all configuration attributes that are
// specific to a nodegroup
type NodeGroup struct {
	ID int

	AMI               string
	InstanceType      string
	AvailabilityZones []string
	Tags              map[string]string
	SubnetTopology    SubnetTopology

	DesiredCapacity int
	MinSize         int
	MaxSize         int

	VolumeSize int

	MaxPodsPerNode int

	PolicyARNs      []string
	InstanceRoleARN string

	AllowSSH         bool
	SSHPublicKeyPath string
	SSHPublicKey     []byte
	SSHPublicKeyName string
}

type (
	// ClusterVPC holds global subnet and all child public/private subnet
	ClusterVPC struct {
		Network              // global CIRD and VPC ID
		SecurityGroup string // cluster SG
		// subnets are either public or private for use with separate nodegroups
		// these are keyed by AZ for conveninece
		Subnets map[SubnetTopology]map[string]Network
		// for additional CIRD associations, e.g. to use with separate CIDR for
		// private subnets or any ad-hoc subnets
		ExtraCIDRs []*net.IPNet
	}
	// SubnetTopology can be SubnetTopologyPrivate or SubnetTopologyPublic
	SubnetTopology string
	// Network holds ID and CIDR
	Network struct {
		ID   string
		CIDR *net.IPNet
	}
)

const (
	// SubnetTopologyPrivate repesents privately-routed subnets
	SubnetTopologyPrivate SubnetTopology = "Private"
	// SubnetTopologyPublic repesents publicly-routed subnets
	SubnetTopologyPublic SubnetTopology = "Public"
)

// SubnetIDs returns list of subnets
func (c ClusterVPC) SubnetIDs(topology SubnetTopology) []string {
	subnets := []string{}
	for _, s := range c.Subnets[topology] {
		subnets = append(subnets, s.ID)
	}
	return subnets
}

// SubnetTopologies returns list of topologies supported
// by a given cluster config
func (c ClusterVPC) SubnetTopologies() []SubnetTopology {
	topologies := []SubnetTopology{}
	for topology := range c.Subnets {
		topologies = append(topologies, topology)
	}
	return topologies
}

// ImportSubnet loads a given subnet into cluster config
func (c ClusterVPC) ImportSubnet(topology SubnetTopology, az, subnetID string) {
	if _, ok := c.Subnets[topology]; !ok {
		c.Subnets[topology] = map[string]Network{}
	}
	if network, ok := c.Subnets[topology][az]; !ok {
		c.Subnets[topology][az] = Network{ID: subnetID}
	} else {
		network.ID = subnetID
		c.Subnets[topology][az] = network
	}
}

// HasSufficientPublicSubnets validates if there is a suffiecent
// number of subnets available to create a cluster
func (c ClusterVPC) HasSufficientPublicSubnets() bool {
	return len(c.SubnetIDs(SubnetTopologyPublic)) >= 3
}

type (
	// ClusterAddons provides addons for the created EKS cluster
	ClusterAddons struct {
		WithIAM AddonIAM
		Storage bool
	}
	// AddonIAM provides an addon for the AWS IAM integration
	AddonIAM struct {
		PolicyAmazonEC2ContainerRegistryPowerUser bool
		PolicyAutoScaling                         bool
	}
)
