package api

import (
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

	// EKSRegionUSEast2 represents the US East Region Ohio
	EKSRegionUSEast2 = "us-east-2"

	// EKSRegionEUWest1 represents the EU West Region Ireland
	EKSRegionEUWest1 = "eu-west-1"

	// DefaultEKSRegion defines the default region, where to deploy the EKS cluster
	DefaultEKSRegion = EKSRegionUSWest2
)

// SupportedRegions are the regions where EKS is available
var SupportedRegions = []string{
	EKSRegionUSWest2,
	EKSRegionUSEast1,
	EKSRegionUSEast2,
	EKSRegionEUWest1,
}

// IsSupportedRegion returns true if EKS is available in the region
func IsSupportedRegion(region string) bool {
	for _, r := range SupportedRegions {
		if r == region {
			return true
		}
	}
	return false
}

// DefaultWaitTimeout defines the default wait timeout
var DefaultWaitTimeout = 20 * time.Minute

// DefaultNodeCount defines the default number of nodes to be created
const DefaultNodeCount = 2

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

	VPC *ClusterVPC

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
	cfg := &ClusterConfig{
		VPC: &ClusterVPC{},
	}

	cidr := DefaultCIDR()
	cfg.VPC.CIDR = &cidr

	return cfg
}

// NewNodeGroup crears new nodegroup inside cluster config,
// it returns pointer to the nodegroup for convenience
func (c *ClusterConfig) NewNodeGroup() *NodeGroup {
	ng := &NodeGroup{
		ID:                len(c.NodeGroups),
		PrivateNetworking: false,
	}

	c.NodeGroups = append(c.NodeGroups, ng)

	return ng
}

// NodeGroup holds all configuration attributes that are
// specific to a nodegroup
type NodeGroup struct {
	ID int

	AMI               string
	AMIFamily         string
	InstanceType      string
	AvailabilityZones []string
	Tags              map[string]string
	PrivateNetworking bool

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

// SubnetTopology check which topology is used for the subnet of
// the given nodegroup
func (n *NodeGroup) SubnetTopology() SubnetTopology {
	if n.PrivateNetworking {
		return SubnetTopologyPrivate
	}
	return SubnetTopologyPublic
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
