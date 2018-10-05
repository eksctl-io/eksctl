package api

import (
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

const (
	AWSDebugLevel = 5

	EKS_REGION_US_WEST_2 = "us-west-2"
	EKS_REGION_US_EAST_1 = "us-east-1"
	EKS_REGION_EU_WEST_1 = "eu-west-1"
	DEFAULT_EKS_REGION   = EKS_REGION_US_WEST_2
)

var DefaultWaitTimeout = 20 * time.Minute

type ClusterProvider interface {
	CloudFormation() cloudformationiface.CloudFormationAPI
	EKS() eksiface.EKSAPI
	EC2() ec2iface.EC2API
	STS() stsiface.STSAPI
}

// simple config, to be replaced with Cluster API
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

type ClusterAddons struct {
	WithIAM AddonIAM
	Storage bool
}

type AddonIAM struct {
	PolicyAmazonEC2ContainerRegistryPowerUser bool
}
