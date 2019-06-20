package v1alpha5

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// AWSDebugLevel defines the LogLevel for AWS produced logs
	AWSDebugLevel = 5

	// RegionUSWest2 represents the US West Region Oregon
	RegionUSWest2 = "us-west-2"

	// RegionUSEast1 represents the US East Region North Virginia
	RegionUSEast1 = "us-east-1"

	// RegionUSEast2 represents the US East Region Ohio
	RegionUSEast2 = "us-east-2"

	// RegionCACentral1 represents the Canada Central Region
	RegionCACentral1 = "ca-central-1"

	// RegionEUWest1 represents the EU West Region Ireland
	RegionEUWest1 = "eu-west-1"

	// RegionEUWest2 represents the EU West Region London
	RegionEUWest2 = "eu-west-2"

	// RegionEUWest3 represents the EU West Region Paris
	RegionEUWest3 = "eu-west-3"

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

	// RegionAPSouth1 represents the Asia-Pacific South Region Mumbai
	RegionAPSouth1 = "ap-south-1"

	// DefaultRegion defines the default region, where to deploy the EKS cluster
	DefaultRegion = RegionUSWest2

	// Version1_10 represents Kubernetes version 1.10.x
	Version1_10 = "1.10"

	// Version1_11 represents Kubernetes version 1.11.x
	Version1_11 = "1.11"

	// Version1_12 represents Kubernetes version 1.12.x
	Version1_12 = "1.12"

	// Version1_13 represents Kubernetes version 1.13.x
	Version1_13 = "1.13"

	// DefaultVersion represents default Kubernetes version supported by EKS
	DefaultVersion = Version1_12

	// LatestVersion represents latest Kubernetes version supported by EKS
	LatestVersion = Version1_13

	// DefaultNodeType is the default instance type to use for nodes
	DefaultNodeType = "m5.large"

	// DefaultNodeCount defines the default number of nodes to be created
	DefaultNodeCount = 2

	// NodeVolumeTypeGP2 is General Purpose SSD
	NodeVolumeTypeGP2 = "gp2"
	// NodeVolumeTypeIO1 is Provisioned IOPS SSD
	NodeVolumeTypeIO1 = "io1"
	// NodeVolumeTypeSC1 is Throughput Optimized HDD
	NodeVolumeTypeSC1 = "sc1"
	// NodeVolumeTypeST1 is Cold HDD
	NodeVolumeTypeST1 = "st1"

	// DefaultNodeImageFamily defines the default image family for the worker nodes
	DefaultNodeImageFamily = NodeImageFamilyAmazonLinux2
	// NodeImageFamilyAmazonLinux2 represents Amazon Linux 2 family
	NodeImageFamilyAmazonLinux2 = "AmazonLinux2"
	// NodeImageFamilyUbuntu1804 represents Ubuntu 18.04 family
	NodeImageFamilyUbuntu1804 = "Ubuntu1804"
	// NodeImageResolverStatic represents static AMI resolver (see ami package)
	NodeImageResolverStatic = "static"
	// NodeImageResolverAuto represents auto AMI resolver (see ami package)
	NodeImageResolverAuto = "auto"

	// ClusterNameTag defines the tag of the cluster name
	ClusterNameTag = "alpha.eksctl.io/cluster-name"

	// OldClusterNameTag defines the tag of the cluster name
	OldClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"

	// NodeGroupNameTag defines the tag of the node group name
	NodeGroupNameTag = "alpha.eksctl.io/nodegroup-name"

	// OldNodeGroupNameTag defines the tag of the node group name
	OldNodeGroupNameTag = "eksctl.io/v1alpha2/nodegroup-name"

	// OldNodeGroupIDTag defines the old version of tag of the node group name
	OldNodeGroupIDTag = "eksctl.cluster.k8s.io/v1alpha1/nodegroup-id"

	// ClusterNameLabel defines the tag of the cluster name
	ClusterNameLabel = "alpha.eksctl.io/cluster-name"

	// NodeGroupNameLabel defines the label of the node group name
	NodeGroupNameLabel = "alpha.eksctl.io/nodegroup-name"
)

var (
	// DefaultWaitTimeout defines the default wait timeout
	DefaultWaitTimeout = 25 * time.Minute

	// DefaultNodeSSHPublicKeyPath is the default path to SSH public key
	DefaultNodeSSHPublicKeyPath = "~/.ssh/id_rsa.pub"

	// DefaultNodeVolumeType defines the default root volume type to use
	DefaultNodeVolumeType = NodeVolumeTypeGP2

	// DefaultNodeVolumeSize defines the default root volume size
	DefaultNodeVolumeSize = 0
)

// Enabled return pointer to true value
// for use in defaulters of *bool fields
func Enabled() *bool {
	v := true
	return &v
}

// Disabled return pointer to false value
// for use in defaulters of *bool fields
func Disabled() *bool {
	v := false
	return &v
}

// IsEnabled will only return true if v is not nil and true
func IsEnabled(v *bool) bool { return v != nil && *v }

// IsDisabled will only return true if v is not nil and false
func IsDisabled(v *bool) bool { return v != nil && !*v }

// IsSetAndNonEmptyString will only return true if s is not nil and not empty
func IsSetAndNonEmptyString(s *string) bool { return s != nil && *s != "" }

// SupportedRegions are the regions where EKS is available
func SupportedRegions() []string {
	return []string{
		RegionUSWest2,
		RegionUSEast1,
		RegionUSEast2,
		RegionCACentral1,
		RegionEUWest1,
		RegionEUWest2,
		RegionEUWest3,
		RegionEUNorth1,
		RegionEUCentral1,
		RegionAPNorthEast1,
		RegionAPNorthEast2,
		RegionAPSouthEast1,
		RegionAPSouthEast2,
		RegionAPSouth1,
	}
}

// SupportedVersions are the versions of Kubernetes that EKS supports
func SupportedVersions() []string {
	return []string{
		Version1_10,
		Version1_11,
		Version1_12,
		Version1_13,
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
	IAM() iamiface.IAMAPI
	CloudTrail() cloudtrailiface.CloudTrailAPI
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
			Version: DefaultVersion,
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
		AutoAllocateIPv6: Disabled(),
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
			AttachIDs:  []string{},
			WithLocal:  Enabled(),
			WithShared: Enabled(),
		},
		DesiredCapacity: nil,
		InstanceType:    DefaultNodeType,
		VolumeSize:      &DefaultNodeVolumeSize,
		VolumeType:      &DefaultNodeVolumeType,
		IAM: &NodeGroupIAM{
			WithAddonPolicies: NodeGroupIAMAddonPolicies{
				ImageBuilder: Disabled(),
				AutoScaler:   Disabled(),
				ExternalDNS:  Disabled(),
				CertManager:  Disabled(),
				AppMesh:      Disabled(),
				EBS:          Disabled(),
				FSX:          Disabled(),
				EFS:          Disabled(),
				ALBIngress:   Disabled(),
				XRay:         Disabled(),
				CloudWatch:   Disabled(),
			},
		},
		SSH: &NodeGroupSSH{
			Allow:         Disabled(),
			PublicKeyPath: &DefaultNodeSSHPublicKeyPath,
		},
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
	//+optional
	InstancesDistribution *NodeGroupInstancesDistribution `json:"instancesDistribution,omitempty"`
	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// +optional
	PrivateNetworking bool `json:"privateNetworking"`

	// +optional
	SecurityGroups *NodeGroupSGs `json:"securityGroups,omitempty"`

	// +optional
	DesiredCapacity *int `json:"desiredCapacity,omitempty"`
	// +optional
	MinSize *int `json:"minSize,omitempty"`
	// +optional
	MaxSize *int `json:"maxSize,omitempty"`

	// +optional
	VolumeSize *int `json:"volumeSize"`
	// +optional
	VolumeType *string `json:"volumeType"`
	// +optional
	VolumeName *string `json:"volumeName,omitempty"`
	// +optional
	VolumeEncrypted *bool `json:"volumeEncrypted,omitempty"`
	// +optional
	MaxPodsPerNode int `json:"maxPodsPerNode,omitempty"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	Taints map[string]string `json:"taints,omitempty"`

	// +optional
	TargetGroupARNs []string `json:"targetGroupARNs,omitempty"`

	SSH *NodeGroupSSH `json:"ssh"`

	// +optional
	IAM *NodeGroupIAM `json:"iam"`

	// +optional
	PreBootstrapCommands []string `json:"preBootstrapCommands,omitempty"`

	// +optional
	OverrideBootstrapCommand *string `json:"overrideBootstrapCommand,omitempty"`

	// +optional
	ClusterDNS string `json:"clusterDNS,omitempty"`

	// +optional
	KubeletExtraConfig *NodeGroupKubeletConfig `json:"kubeletExtraConfig,omitempty"`
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
		AttachIDs []string `json:"attachIDs,omitempty"`
		// +optional
		WithShared *bool `json:"withShared"`
		// +optional
		WithLocal *bool `json:"withLocal"`
	}
	// NodeGroupIAM holds all IAM attributes of a NodeGroup
	NodeGroupIAM struct {
		// +optional
		AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`
		// +optional
		InstanceProfileARN string `json:"instanceProfileARN,omitempty"`
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
		ImageBuilder *bool `json:"imageBuilder"`
		// +optional
		AutoScaler *bool `json:"autoScaler"`
		// +optional
		ExternalDNS *bool `json:"externalDNS"`
		// +optional
		CertManager *bool `json:"certManager"`
		// +optional
		AppMesh *bool `json:"appMesh"`
		// +optional
		EBS *bool `json:"ebs"`
		// +optional
		FSX *bool `json:"fsx"`
		// +optional
		EFS *bool `json:"efs"`
		// +optional
		ALBIngress *bool `json:"albIngress"`
		// +optional
		XRay *bool `json:"xRay"`
		// +optional
		CloudWatch *bool `json:"cloudWatch"`
	}

	// NodeGroupSSH holds all the ssh access configuration to a NodeGroup
	NodeGroupSSH struct {
		// +optional
		Allow *bool `json:"allow"`
		// +optional
		PublicKeyPath *string `json:"publicKeyPath,omitempty"`
		// +optional
		PublicKey *string `json:"publicKey,omitempty"`
		// +optional
		PublicKeyName *string `json:"publicKeyName,omitempty"`
	}

	// NodeGroupInstancesDistribution holds the configuration for spot instances
	NodeGroupInstancesDistribution struct {
		//+required
		InstanceTypes []string `json:"instanceTypes,omitEmpty"`
		// +optional
		MaxPrice *float64 `json:"maxPrice,omitempty"`
		//+optional
		OnDemandBaseCapacity *int `json:"onDemandBaseCapacity,omitEmpty"`
		//+optional
		OnDemandPercentageAboveBaseCapacity *int `json:"onDemandPercentageAboveBaseCapacity,omitEmpty"`
		//+optional
		SpotInstancePools *int `json:"spotInstancePools,omitEmpty"`
	}
)

// NodeGroupKubeletConfig contains extra config parameters for the kubelet.yaml
type NodeGroupKubeletConfig map[string]interface{}

// HasMixedInstances checks if a nodegroup has mixed instances option declared
func HasMixedInstances(ng *NodeGroup) bool {
	return ng.InstancesDistribution != nil && ng.InstancesDistribution.InstanceTypes != nil && len(ng.InstancesDistribution.InstanceTypes) != 0
}
