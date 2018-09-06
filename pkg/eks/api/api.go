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
}

type AddonIAM struct {
	PolicyAmazonEC2ContainerRegistryPowerUser bool
}
