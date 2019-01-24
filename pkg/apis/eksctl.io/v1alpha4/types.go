package v1alpha4

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

	// RegionAPNorthEast2 represents the Asia-Pacific North East Region Seoul
	RegionAPNorthEast2 = "ap-northeast-2"

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

	// DefaultNodeCount defines the default number of nodes to be created
	DefaultNodeCount = 2

	// DefaultNodeVolumeType defines the default root volume type to use
	DefaultNodeVolumeType = NodeVolumeTypeGP2
	// NodeVolumeTypeGP2 is General Purpose SSD
	NodeVolumeTypeGP2 = "gp2"
	// NodeVolumeTypeIO1 is Provisioned IOPS SSD
	NodeVolumeTypeIO1 = "io1"
	// NodeVolumeTypeSC1 is Throughput Optimized HDD
	NodeVolumeTypeSC1 = "sc1"
	// NodeVolumeTypeST1 is Cold HDD
	NodeVolumeTypeST1 = "st1"

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

var (
	// DefaultWaitTimeout defines the default wait timeout
	DefaultWaitTimeout = 20 * time.Minute
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
		RegionAPNorthEast2,
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

// SupportedNodeVolumeTypes are the volume types that can be used for a node root volume
func SupportedNodeVolumeTypes() []string {
	return []string{
		NodeVolumeTypeGP2,
		NodeVolumeTypeIO1,
		NodeVolumeTypeSC1,
		NodeVolumeTypeST1,
	}
}

// ClusterMeta is what identifies a cluster
type ClusterMeta struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	// +optional
	Version string `json:"version,omitempty"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// ClusterStatus hold read-only attributes of a cluster
type ClusterStatus struct {
	Endpoint                 string `json:"endpoint,omitempty"`
	CertificateAuthorityData []byte `json:"certificateAuthorityData,omitempty"`
	ARN                      string `json:"arn,omitempty"`
	StackName                string `json:"stackName,omitempty"`
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
	IAM ClusterIAM `json:"iam"`

	// +optional
	VPC *ClusterVPC `json:"vpc,omitempty"`

	// +optional
	NodeGroups []*NodeGroup `json:"nodeGroups,omitempty"`

	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	Status *ClusterStatus `json:"status,omitempty"`
}

// ClusterIAM holds all IAM attributes of a cluster
type ClusterIAM struct {
	// +optional
	ServiceRoleARN string `json:"serviceRoleARN,omitempty"`
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
		SecurityGroups: &NodeGroupSGs{
			WithShared: true,
			WithLocal:  true,
			AttachIDs:  []string{},
		},
		DesiredCapacity: DefaultNodeCount,
		InstanceType:    DefaultNodeType,
		VolumeSize:      0,
		VolumeType:      DefaultNodeVolumeType,
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
	SecurityGroups *NodeGroupSGs `json:"securityGroups,omitempty"`

	// +optional
	DesiredCapacity int `json:"desiredCapacity"`
	// +optional
	MinSize int `json:"minSize,omitempty"`
	// +optional
	MaxSize int `json:"maxSize,omitempty"`

	// +optional
	VolumeSize int `json:"volumeSize"`
	// +optional
	VolumeType string `json:"volumeType"`
	// +optional
	MaxPodsPerNode int `json:"maxPodsPerNode,omitempty"`

	// +optional
	Labels NodeLabels `json:"labels,omitempty"`

	// TODO move to separate struct
	// +optional
	AllowSSH bool `json:"allowSSH"`
	// +optional
	SSHPublicKeyPath string `json:"sshPublicKeyPath,omitempty"`
	// +optional
	SSHPublicKey []byte `json:"SSHPublicKey,omitempty"` // TODO: right now it's kind of read-only, but one may wish to use key body in a config file so we will need recognise that
	// +optional
	SSHPublicKeyName string `json:"sshPublicKeyName,omitempty"`

	// +optional
	IAM NodeGroupIAM `json:"iam"`
}

// SubnetTopology check which topology is used for the subnet of
// the given nodegroup
func (n *NodeGroup) SubnetTopology() SubnetTopology {
	if n.PrivateNetworking {
		return SubnetTopologyPrivate
	}
	return SubnetTopologyPublic
}

// ListOptions returns metav1.ListOptions with label selector for the nodegroup
func (n *NodeGroup) ListOptions() metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", NodeGroupNameLabel, n.Name),
	}
}

type (
	// NodeGroupSGs holds all SG attributes of a NodeGroup
	NodeGroupSGs struct {
		// +optional
		AttachIDs []string
		// +optional
		WithShared bool
		// +optional
		WithLocal bool
	}
	// NodeGroupIAM holds all IAM attributes of a NodeGroup
	NodeGroupIAM struct {
		// +optional
		AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`
		// +optional
		InstanceRoleARN string `json:"instanceRoleARN,omitempty"`
		// +optional
		InstanceRoleName string `json:"instanceRoleName,omitempty"`
		// +optional
		WithAddonPolicies NodeGroupIAMAddonPolicies `json:"withAddonPolicies,omitempty"`
	}
	// NodeGroupIAMAddonPolicies holds all IAM addon policies
	NodeGroupIAMAddonPolicies struct {
		// +optional
		ImageBuilder bool `json:"imageBuilder"`
		// +optional
		AutoScaler bool `json:"autoScaler"`
		// +optional
		ExternalDNS bool `json:"externalDNS"`
	}
)
