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

	// EKSRegionEUWest1 represents the EU West Region Ireland
	EKSRegionEUWest1 = "eu-west-1"

	// DefaultEKSRegion defines the default region, where to deploy the EKS cluster
	DefaultEKSRegion = EKSRegionUSWest2
)

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

	NodeAMI  string
	NodeType string
	Nodes    int
	MinNodes int
	MaxNodes int

	NodeVolumeSize int

	MaxPodsPerNode int

	NodePolicyARNs []string

	NodeSSH          bool
	SSHPublicKeyPath string
	SSHPublicKey     []byte
	SSHPublicKeyName string

	WaitTimeout time.Duration

	SecurityGroup string
	Subnets       []string
	VPC           string

	Endpoint                 string
	CertificateAuthorityData []byte
	ARN                      string

	ClusterStackName string

	NodeInstanceRoleARN string

	AvailabilityZones []string

	Addons ClusterAddons
}

// ClusterAddons provides addons for the created EKS cluster
type ClusterAddons struct {
	WithIAM AddonIAM
	Storage bool
}

// AddonIAM provides an addon for the AWS IAM integration
type AddonIAM struct {
	PolicyAmazonEC2ContainerRegistryPowerUser bool
}
