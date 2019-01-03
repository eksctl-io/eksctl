package v1alpha3

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// AWSDebugLevel defines the LogLevel for AWS produced logs
	AWSDebugLevel = 5

	// RegionUSWest2 represents the US West Region Oregon
	RegionUSWest2 = "us-west-2"

	// RegionUSEast1 represents the US East Region North Virgina
	RegionUSEast1 = "us-east-1"

	// RegionUSEast2 represents the US East Region Ohio
	RegionUSEast2 = "us-east-2"

	// RegionEUWest1 represents the EU West Region Ireland
	RegionEUWest1 = "eu-west-1"

	// RegionEUNorth1 represents the EU North Region Stockholm
	RegionEUNorth1 = "eu-north-1"

	// RegionEUCentral1 represents the EU Central Region Frankfurt
	RegionEUCentral1 = "eu-central-1"

	// RegionAPNorthEast1 represents the Asia-Pacific North East Region Tokyo
	RegionAPNorthEast1 = "ap-northeast-1"

	// RegionAPSouthEast1 represents the Asia-Pacific South East Region Singapore
	RegionAPSouthEast1 = "ap-southeast-1"

	// RegionAPSouthEast2 represents the Asia-Pacific South East Region Sydney
	RegionAPSouthEast2 = "ap-southeast-2"

	// DefaultRegion defines the default region, where to deploy the EKS cluster
	DefaultRegion = RegionUSWest2

	// Version1_10 represents Kubernetes version 1.10.x
	Version1_10 = "1.10"

	// Version1_11 represents Kubernetes version 1.11.x
	Version1_11 = "1.11"

	// LatestVersion represents latest Kubernetes version supported by EKS
	LatestVersion = Version1_11

	// DefaultNodeType is the default instance type to use for nodes
	DefaultNodeType = "m5.large"

	// ClusterNameTag defines the tag of the clsuter name
	ClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"

	// NodeGroupNameTag defines the tag of the node group name
	NodeGroupNameTag = "eksctl.io/v1alpha2/nodegroup-name"
	// OldNodeGroupIDTag defines the old version of tag of the node group name
	OldNodeGroupIDTag = "eksctl.cluster.k8s.io/v1alpha1/nodegroup-id"

	// ClusterNameLabel defines the tag of the clsuter name
	ClusterNameLabel = "alpha.eksctl.io/cluster-name"

	// NodeGroupNameLabel defines the label of the node group name
	NodeGroupNameLabel = "alpha.eksctl.io/nodegroup-name"
)

// SupportedRegions are the regions where EKS is available
func SupportedRegions() []string {
	return []string{
		RegionUSWest2,
		RegionUSEast1,
		RegionUSEast2,
		RegionEUWest1,
		RegionEUNorth1,
		RegionEUCentral1,
		RegionAPNorthEast1,
		RegionAPSouthEast1,
		RegionAPSouthEast2,
	}
}

// SupportedVersions are the versions of Kubernetes that EKS supports
func SupportedVersions() []string {
	return []string{
		Version1_10,
		Version1_11,
	}
}

// DefaultWaitTimeout defines the default wait timeout
var DefaultWaitTimeout = 20 * time.Minute

// DefaultNodeCount defines the default number of nodes to be created
const DefaultNodeCount = 2

// ClusterMeta is what identifies a cluster
type ClusterMeta struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	// +optional
	Version string `json:"version,omitempty"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// String returns canonical representation of ClusterMeta
func (c *ClusterMeta) String() string {
	return fmt.Sprintf("%s.%s.eksctl.io", c.Name, c.Region)
}

// LogString returns representation of ClusterMeta for logs
func (c *ClusterMeta) LogString() string {
	return fmt.Sprintf("EKS cluster %q in %q region", c.Name, c.Region)
}

// ClusterProvider is the interface to AWS APIs
type ClusterProvider interface {
	CloudFormation() cloudformationiface.CloudFormationAPI
	CloudFormationRoleARN() string
	EKS() eksiface.EKSAPI
	EC2() ec2iface.EC2API
	STS() stsiface.STSAPI
	Region() string
	Profile() string
	WaitTimeout() time.Duration
}

// ProviderConfig holds global parameters for all interactions with AWS APIs
type ProviderConfig struct {
	CloudFormationRoleARN string

	Region      string
	Profile     string
	WaitTimeout time.Duration
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfig is a simple config, to be replaced with Cluster API
type ClusterConfig struct {
	metav1.TypeMeta `json:",inline"`

	Metadata *ClusterMeta `json:"metadata"`

	// +optional
	VPC *ClusterVPC `json:"vpc,omitempty"`

	// +optional
	NodeGroups []*NodeGroup `json:"nodeGroups,omitempty"`

	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// TODO: refactor and move IAM addons to nodegroup
	Addons ClusterAddons

	// TODO: move under status
	Endpoint                 string
	CertificateAuthorityData []byte
	ARN                      string
	ClusterStackName         string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfigList is a list of ClusterConfigs
type ClusterConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ClusterConfig `json:"items"`
}

// ClusterConfigTypeMeta constructs TypeMeta for ClusterConfig
func ClusterConfigTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ClusterConfigKind,
		APIVersion: SchemeGroupVersion.String(),
	}
}

// NewClusterConfig creates new config for a cluster;
// it doesn't include initial nodegroup, so user must
// call NewNodeGroup to create one
func NewClusterConfig() *ClusterConfig {
	cfg := &ClusterConfig{
		TypeMeta: ClusterConfigTypeMeta(),
		Metadata: &ClusterMeta{
			Version: LatestVersion,
		},
		VPC: NewClusterVPC(),
	}

	return cfg
}

// NewClusterVPC creates new VPC config for a cluster
func NewClusterVPC() *ClusterVPC {
	cidr := DefaultCIDR()
	return &ClusterVPC{
		Network: Network{
			CIDR: &cidr,
		},
	}
}

// AppendAvailabilityZone appends a new AZ to the set
func (c *ClusterConfig) AppendAvailabilityZone(newAZ string) {
	for _, az := range c.AvailabilityZones {
		if az == newAZ {
			return
		}
	}
	c.AvailabilityZones = append(c.AvailabilityZones, newAZ)
}

// NewNodeGroup creates new nodegroup inside cluster config,
// it returns pointer to the nodegroup for convenience
func (c *ClusterConfig) NewNodeGroup() *NodeGroup {
	ng := &NodeGroup{
		PrivateNetworking: false,
		DesiredCapacity:   DefaultNodeCount,
		InstanceType:      DefaultNodeType,
	}

	c.NodeGroups = append(c.NodeGroups, ng)

	return ng
}

// NodeGroup holds all configuration attributes that are
// specific to a nodegroup
type NodeGroup struct {
	Name string `json:"name"`
	// +optional
	AMI string `json:"ami,omitempty"`
	// +optional
	AMIFamily string `json:"amiFamily,omitempty"`
	// +optional
	InstanceType string `json:"instanceType,omitempty"`
	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// +optional
	PrivateNetworking bool `json:"privateNetworking"`

	// +optional
	DesiredCapacity int `json:"desiredCapacity"`
	// +optional
	MinSize int `json:"minSize,omitempty"`
	// +optional
	MaxSize int `json:"maxSize,omitempty"`

	// +optional
	VolumeSize int `json:"volumeSize"`
	// +optional
	MaxPodsPerNode int `json:"maxPodsPerNode,omitempty"`

	// +optional
	Labels NodeLabels `json:"labels,omitempty"`

	// +optional
	PolicyARNs []string `json:"policyARNs,omitempty"`
	// +optional
	InstanceRoleARN string `json:"instanceRoleARN,omitempty"`

	// TODO move to separate struct
	// +optional
	AllowSSH bool `json:"allowSSH"`
	// +optional
	SSHPublicKeyPath string `json:"sshPublicKeyPath,omitempty"`
	// +optional
	SSHPublicKey []byte `json:"SSHPublicKey,omitempty"` // TODO: right now it's kind of read-only, but one may wish to use key body in a config file so we will need recognise that
	// +optional
	SSHPublicKeyName string `json:"sshPublicKeyName,omitempty"`
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
		PolicyExternalDNS                         bool
	}
)
