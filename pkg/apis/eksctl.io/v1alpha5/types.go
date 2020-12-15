package v1alpha5

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Values for `KubernetesVersion`
// All valid values should go in this block
const (
	Version1_14 = "1.14"

	Version1_15 = "1.15"

	Version1_16 = "1.16"

	Version1_17 = "1.17"

	Version1_18 = "1.18"

	// DefaultVersion (default)
	DefaultVersion = Version1_18

	LatestVersion = Version1_18
)

// No longer supported versions
const (
	// Version1_10 represents Kubernetes version 1.10.x
	Version1_10 = "1.10"

	// Version1_11 represents Kubernetes version 1.11.x
	Version1_11 = "1.11"

	// Version1_12 represents Kubernetes version 1.12.x
	Version1_12 = "1.12"

	// Version1_13 represents Kubernetes version 1.13.x
	Version1_13 = "1.13"
)

// Not yet supported versions
const (
	// Version1_19 represents Kubernetes version 1.19.x
	Version1_19 = "1.19"
)

const (
	// AWSDebugLevel defines the LogLevel for AWS produced logs
	AWSDebugLevel = 5
)

// Regions
const (
	// RegionUSWest1 represents the US West Region North California
	RegionUSWest1 = "us-west-1"

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

	// RegionEUSouth1 represents te Eu South Region Milan
	RegionEUSouth1 = "eu-south-1"

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

	// RegionAPEast1 represents the Asia Pacific Region Hong Kong
	RegionAPEast1 = "ap-east-1"

	// RegionMESouth1 represents the Middle East Region Bahrain
	RegionMESouth1 = "me-south-1"

	// RegionSAEast1 represents the South America Region Sao Paulo
	RegionSAEast1 = "sa-east-1"

	// RegionAFSouth1 represents the Africa Region Cape Town
	RegionAFSouth1 = "af-south-1"

	// RegionCNNorthwest1 represents the China region Ningxia
	RegionCNNorthwest1 = "cn-northwest-1"

	// RegionCNNorth1 represents the China region Beijing
	RegionCNNorth1 = "cn-north-1"

	// RegionUSGovWest1 represents the region GovCloud (US-West)
	RegionUSGovWest1 = "us-gov-west-1"

	// RegionUSGovEast1 represents the region GovCloud (US-East)
	RegionUSGovEast1 = "us-gov-east-1"

	// DefaultRegion defines the default region, where to deploy the EKS cluster
	DefaultRegion = RegionUSWest2
)

// Partitions
const (
	PartitionAWS   = "aws"
	PartitionChina = "aws-cn"
	PartitionUSGov = "aws-us-gov"
)

// Values for `NodeAMIFamily`
// All valid values should go in this block
const (
	// DefaultNodeImageFamily (default)
	DefaultNodeImageFamily      = NodeImageFamilyAmazonLinux2
	NodeImageFamilyAmazonLinux2 = "AmazonLinux2"
	NodeImageFamilyUbuntu2004   = "Ubuntu2004"
	NodeImageFamilyUbuntu1804   = "Ubuntu1804"
	NodeImageFamilyBottlerocket = "Bottlerocket"

	NodeImageFamilyWindowsServer2019CoreContainer = "WindowsServer2019CoreContainer"
	NodeImageFamilyWindowsServer2019FullContainer = "WindowsServer2019FullContainer"
	NodeImageFamilyWindowsServer1909CoreContainer = "WindowsServer1909CoreContainer"
	NodeImageFamilyWindowsServer2004CoreContainer = "WindowsServer2004CoreContainer"
)

const (
	// DefaultNodeType is the default instance type to use for nodes
	DefaultNodeType = "m5.large"

	// DefaultNodeCount defines the default number of nodes to be created
	DefaultNodeCount = 2

	// NodeImageResolverAuto represents auto AMI resolver (see ami package)
	NodeImageResolverAuto = "auto"
	// NodeImageResolverAutoSSM is used to indicate that the latest EKS AMIs should be used for the nodes. The AMI is selected
	// using an SSM GetParameter query
	NodeImageResolverAutoSSM = "auto-ssm"

	// EksctlVersionTag defines the version of eksctl which is used to provision or update EKS cluster
	EksctlVersionTag = "alpha.eksctl.io/eksctl-version"

	// ClusterNameTag defines the tag of the cluster name
	ClusterNameTag = "alpha.eksctl.io/cluster-name"

	// OldClusterNameTag defines the tag of the cluster name
	OldClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"

	// NodeGroupNameTag defines the tag of the nodegroup name
	NodeGroupNameTag = "alpha.eksctl.io/nodegroup-name"

	// NodeGroupTypeTag defines the nodegroup type as managed or unmanaged
	NodeGroupTypeTag = "alpha.eksctl.io/nodegroup-type"

	// OldNodeGroupNameTag defines the tag of the nodegroup name
	OldNodeGroupNameTag = "eksctl.io/v1alpha2/nodegroup-name"

	// OldNodeGroupIDTag defines the old version of tag of the nodegroup name
	OldNodeGroupIDTag = "eksctl.cluster.k8s.io/v1alpha1/nodegroup-id"

	// IAMServiceAccountNameTag defines the tag of the IAM service account name
	IAMServiceAccountNameTag = "alpha.eksctl.io/iamserviceaccount-name"

	// AddonNameTag defines the tag of the IAM service account name
	AddonNameTag = "alpha.eksctl.io/addon-name"

	// ClusterNameLabel defines the tag of the cluster name
	ClusterNameLabel = "alpha.eksctl.io/cluster-name"

	// NodeGroupNameLabel defines the label of the nodegroup name
	NodeGroupNameLabel = "alpha.eksctl.io/nodegroup-name"

	// SpotAllocationStrategyLowestPrice defines the ASG spot allocation strategy of lowest-price
	SpotAllocationStrategyLowestPrice = "lowest-price"

	// SpotAllocationStrategyCapacityOptimized defines the ASG spot allocation strategy of capacity-optimized
	SpotAllocationStrategyCapacityOptimized = "capacity-optimized"

	// eksResourceAccountStandard defines the AWS EKS account ID that provides node resources in default regions
	// for standard AWS partition
	eksResourceAccountStandard = "602401143452"

	// eksResourceAccountAPEast1 defines the AWS EKS account ID that provides node resources in ap-east-1 region
	eksResourceAccountAPEast1 = "800184023465"

	// eksResourceAccountMESouth1 defines the AWS EKS account ID that provides node resources in me-south-1 region
	eksResourceAccountMESouth1 = "558608220178"

	// eksResourceAccountCNNorthWest1 defines the AWS EKS account ID that provides node resources in cn-northwest-1 region
	eksResourceAccountCNNorthWest1 = "961992271922"

	// eksResourceAccountCNNorth1 defines the AWS EKS account ID that provides node resources in cn-north-1
	eksResourceAccountCNNorth1 = "918309763551"

	// eksResourceAccountAFSouth1 defines the AWS EKS account ID that provides node resources in af-south-1
	eksResourceAccountAFSouth1 = "877085696533"

	// eksResourceAccountEUSouth1 defines the AWS EKS account ID that provides node resources in eu-south-1
	eksResourceAccountEUSouth1 = "590381155156"

	// eksResourceAccountUSGovWest1 defines the AWS EKS account ID that provides node resources in us-gov-west-1
	eksResourceAccountUSGovWest1 = "013241004608"

	// eksResourceAccountUSGovEast1 defines the AWS EKS account ID that provides node resources in us-gov-east-1
	eksResourceAccountUSGovEast1 = "151742754352"
)

// Values for `VolumeType`
const (
	// NodeVolumeTypeGP2 is General Purpose SSD (default)
	NodeVolumeTypeGP2 = "gp2"
	// NodeVolumeTypeIO1 is Provisioned IOPS SSD
	NodeVolumeTypeIO1 = "io1"
	// NodeVolumeTypeSC1 is Cold HDD
	NodeVolumeTypeSC1 = "sc1"
	// NodeVolumeTypeST1 is Throughput Optimized HDD
	NodeVolumeTypeST1 = "st1"
)

// NodeGroupType defines the nodegroup type
type NodeGroupType string

const (
	// NodeGroupTypeManaged defines a managed nodegroup
	NodeGroupTypeManaged NodeGroupType = "managed"
	// NodeGroupTypeUnmanaged defines an unmanaged nodegroup
	NodeGroupTypeUnmanaged NodeGroupType = "unmanaged"
)

var (
	// DefaultWaitTimeout defines the default wait timeout
	DefaultWaitTimeout = 25 * time.Minute

	// DefaultNodeSSHPublicKeyPath is the default path to SSH public key
	DefaultNodeSSHPublicKeyPath = "~/.ssh/id_rsa.pub"

	// DefaultNodeVolumeType defines the default root volume type to use
	DefaultNodeVolumeType = NodeVolumeTypeGP2

	// DefaultNodeVolumeSize defines the default root volume size
	DefaultNodeVolumeSize = 80
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
		RegionUSWest1,
		RegionUSWest2,
		RegionUSEast1,
		RegionUSEast2,
		RegionCACentral1,
		RegionEUWest1,
		RegionEUWest2,
		RegionEUWest3,
		RegionEUNorth1,
		RegionEUCentral1,
		RegionEUSouth1,
		RegionAPNorthEast1,
		RegionAPNorthEast2,
		RegionAPSouthEast1,
		RegionAPSouthEast2,
		RegionAPSouth1,
		RegionAPEast1,
		RegionMESouth1,
		RegionSAEast1,
		RegionAFSouth1,
		RegionCNNorthwest1,
		RegionCNNorth1,
		RegionUSGovWest1,
		RegionUSGovEast1,
	}
}

// Partition gives the partition a region belongs to
func Partition(region string) string {
	switch region {
	case RegionUSGovWest1, RegionUSGovEast1:
		return PartitionUSGov
	case RegionCNNorth1, RegionCNNorthwest1:
		return PartitionChina
	default:
		return PartitionAWS
	}
}

// DeprecatedVersions are the versions of Kubernetes that EKS used to support
// but no longer does. See also:
// https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html
func DeprecatedVersions() []string {
	return []string{
		Version1_10,
		Version1_11,
		Version1_12,
		Version1_13,
		Version1_14,
	}
}

// IsDeprecatedVersion returns true if the given Kubernetes version has been deprecated in EKS
func IsDeprecatedVersion(version string) bool {
	for _, v := range DeprecatedVersions() {
		if version == v {
			return true
		}
	}
	return false
}

// SupportedVersions are the versions of Kubernetes that EKS supports
func SupportedVersions() []string {
	return []string{
		Version1_15,
		Version1_16,
		Version1_17,
		Version1_18,
	}
}

// IsSupportedVersion returns true if the given version is a Kubernetes supported by eksctl and EKS
func IsSupportedVersion(version string) bool {
	for _, v := range SupportedVersions() {
		if version == v {
			return true
		}
	}
	return false
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

// supportedSpotAllocationStrategies are the spot allocation strategies supported by ASG
func supportedSpotAllocationStrategies() []string {
	return []string{
		SpotAllocationStrategyLowestPrice,
		SpotAllocationStrategyCapacityOptimized,
	}
}

// isSpotAllocationStrategySupported returns true if the spot allocation strategy is supported for ASG
func isSpotAllocationStrategySupported(allocationStrategy string) bool {
	for _, strategy := range supportedSpotAllocationStrategies() {
		if strategy == allocationStrategy {
			return true
		}
	}
	return false
}

// EKSResourceAccountID provides worker node resources(ami/ecr image) in different aws account
// for different aws partitions & opt-in regions.
func EKSResourceAccountID(region string) string {
	switch region {
	case RegionAPEast1:
		return eksResourceAccountAPEast1
	case RegionMESouth1:
		return eksResourceAccountMESouth1
	case RegionCNNorthwest1:
		return eksResourceAccountCNNorthWest1
	case RegionCNNorth1:
		return eksResourceAccountCNNorth1
	case RegionUSGovWest1:
		return eksResourceAccountUSGovWest1
	case RegionUSGovEast1:
		return eksResourceAccountUSGovEast1
	case RegionAFSouth1:
		return eksResourceAccountAFSouth1
	case RegionEUSouth1:
		return eksResourceAccountEUSouth1
	default:
		return eksResourceAccountStandard
	}
}

// ClusterMeta contains general cluster information
type ClusterMeta struct {
	// Name of the cluster
	// +required
	Name string `json:"name"`
	// the AWS region hosting this cluster
	// +required
	Region string `json:"region"`
	// Valid variants are `KubernetesVersion` constants
	// +optional
	Version string `json:"version,omitempty"`
	// Tags are used to tag AWS resources created by eksctl
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// Annotations are arbitrary metadata ignored by `eksctl`.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// KubernetesNetworkConfig contains cluster networking options
type KubernetesNetworkConfig struct {
	// ServiceIPv4CIDR is the CIDR range from where `ClusterIP`s are assigned
	ServiceIPv4CIDR string `json:"serviceIPv4CIDR,omitempty"`
}

type EKSCTLCreated string

// ClusterStatus hold read-only attributes of a cluster
type ClusterStatus struct {
	Endpoint                 string        `json:"endpoint,omitempty"`
	CertificateAuthorityData []byte        `json:"certificateAuthorityData,omitempty"`
	ARN                      string        `json:"arn,omitempty"`
	StackName                string        `json:"stackName,omitempty"`
	EKSCTLCreated            EKSCTLCreated `json:"eksctlCreated,omitempty"`
}

// String returns canonical representation of ClusterMeta
func (c *ClusterMeta) String() string {
	return fmt.Sprintf("%s.%s.eksctl.io", c.Name, c.Region)
}

// LogString returns representation of ClusterMeta for logs
func (c *ClusterMeta) LogString() string {
	return fmt.Sprintf("EKS cluster %q in %q region", c.Name, c.Region)
}

// LogString returns representation of ClusterConfig for logs
func (c ClusterConfig) LogString() string {
	modes := []string{}
	if c.IsFargateEnabled() {
		modes = append(modes, "Fargate profile")
	}
	if len(c.ManagedNodeGroups) > 0 {
		modes = append(modes, "managed nodes")
	}
	if len(c.NodeGroups) > 0 {
		modes = append(modes, "un-managed nodes")
	}
	return fmt.Sprintf("%s with %s", c.Metadata.LogString(), strings.Join(modes, " and "))
}

// IsFargateEnabled returns true if Fargate is enabled in this ClusterConfig,
// or false otherwise.
func (c ClusterConfig) IsFargateEnabled() bool {
	return len(c.FargateProfiles) > 0
}

// ClusterProvider is the interface to AWS APIs
type ClusterProvider interface {
	CloudFormation() cloudformationiface.CloudFormationAPI
	CloudFormationRoleARN() string
	CloudFormationDisableRollback() bool
	EKS() eksiface.EKSAPI
	EC2() ec2iface.EC2API
	ELB() elbiface.ELBAPI
	ELBV2() elbv2iface.ELBV2API
	STS() stsiface.STSAPI
	SSM() ssmiface.SSMAPI
	IAM() iamiface.IAMAPI
	CloudTrail() cloudtrailiface.CloudTrailAPI
	Region() string
	Profile() string
	WaitTimeout() time.Duration
}

// ProviderConfig holds global parameters for all interactions with AWS APIs
type ProviderConfig struct {
	CloudFormationRoleARN         string
	CloudFormationDisableRollback bool

	Region      string
	Profile     string
	WaitTimeout time.Duration
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfig is a simple config, to be replaced with Cluster API
type ClusterConfig struct {
	metav1.TypeMeta

	// +required
	Metadata *ClusterMeta `json:"metadata"`

	// +optional
	KubernetesNetworkConfig *KubernetesNetworkConfig `json:"kubernetesNetworkConfig,omitempty"`

	// +optional
	IAM *ClusterIAM `json:"iam,omitempty"`

	// +optional
	VPC *ClusterVPC `json:"vpc,omitempty"`

	// +optional
	Addons []*Addon `json:"addons,omitempty"`

	// PrivateCluster allows configuring a fully-private cluster
	// in which no node has outbound internet access, and private access
	// to AWS services is enabled via VPC endpoints
	// +optional
	PrivateCluster *PrivateCluster `json:"privateCluster,omitempty"`

	// NodeGroups For information and examples see [nodegroups](/usage/managing-nodegroups)
	// +optional
	NodeGroups []*NodeGroup `json:"nodeGroups,omitempty"`

	// ManagedNodeGroups See [Nodegroups usage](/usage/managing-nodegroups)
	// and [managed nodegroups](/usage/eks-managed-nodes/)
	// +optional
	ManagedNodeGroups []*ManagedNodeGroup `json:"managedNodeGroups,omitempty"`

	// +optional
	FargateProfiles []*FargateProfile `json:"fargateProfiles,omitempty"`

	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// See [CloudWatch support](/usage/cloudwatch-cluster-logging/)
	// +optional
	CloudWatch *ClusterCloudWatch `json:"cloudWatch,omitempty"`

	// +optional
	SecretsEncryption *SecretsEncryption `json:"secretsEncryption,omitempty"`

	Status *ClusterStatus `json:"status,omitempty"`

	// +optional
	Git *Git `json:"git,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfigList is a list of ClusterConfigs
type ClusterConfigList struct {
	metav1.TypeMeta
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
		IAM: NewClusterIAM(),
		VPC: NewClusterVPC(),
		CloudWatch: &ClusterCloudWatch{
			ClusterLogging: &ClusterCloudWatchLogging{},
		},
		PrivateCluster: &PrivateCluster{},
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
		NAT:              DefaultClusterNAT(),
		AutoAllocateIPv6: Disabled(),
		ClusterEndpoints: &ClusterEndpoints{},
	}
}

// NewClusterIAM creates a new ClusterIAM for a cluster
func NewClusterIAM() *ClusterIAM {
	return &ClusterIAM{
		WithOIDC: Disabled(),
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

// NewNodeGroup creates a new NodeGroup, and returns a pointer to it
func NewNodeGroup() *NodeGroup {
	return &NodeGroup{
		NodeGroupBase: &NodeGroupBase{
			PrivateNetworking: false,
			InstanceType:      DefaultNodeType,
			VolumeSize:        &DefaultNodeVolumeSize,
			IAM: &NodeGroupIAM{
				WithAddonPolicies: NodeGroupIAMAddonPolicies{
					ImageBuilder:              Disabled(),
					AutoScaler:                Disabled(),
					ExternalDNS:               Disabled(),
					CertManager:               Disabled(),
					AppMesh:                   Disabled(),
					AppMeshPreview:            Disabled(),
					EBS:                       Disabled(),
					FSX:                       Disabled(),
					EFS:                       Disabled(),
					AWSLoadBalancerController: Disabled(),
					XRay:                      Disabled(),
					CloudWatch:                Disabled(),
				},
			},
			ScalingConfig: &ScalingConfig{},
			SSH: &NodeGroupSSH{
				Allow:         Disabled(),
				PublicKeyPath: &DefaultNodeSSHPublicKeyPath,
			},
			VolumeType: &DefaultNodeVolumeType,
			SecurityGroups: &NodeGroupSGs{
				AttachIDs:  []string{},
				WithLocal:  Enabled(),
				WithShared: Enabled(),
			},
			DisableIMDSv1:  Disabled(),
			DisablePodIMDS: Disabled(),
		},
	}
}

// NewManagedNodeGroup creates a new ManagedNodeGroup
func NewManagedNodeGroup() *ManagedNodeGroup {
	var (
		publicKey  = DefaultNodeSSHPublicKeyPath
		volumeSize = DefaultNodeVolumeSize
		volumeType = DefaultNodeVolumeType
	)
	return &ManagedNodeGroup{
		NodeGroupBase: &NodeGroupBase{
			VolumeSize: &volumeSize,
			VolumeType: &volumeType,
			SSH: &NodeGroupSSH{
				Allow:         Disabled(),
				PublicKeyName: &publicKey,
				EnableSSM:     Disabled(),
			},
			IAM: &NodeGroupIAM{
				WithAddonPolicies: NodeGroupIAMAddonPolicies{
					ImageBuilder:              Disabled(),
					AutoScaler:                Disabled(),
					ExternalDNS:               Disabled(),
					CertManager:               Disabled(),
					AppMesh:                   Disabled(),
					AppMeshPreview:            Disabled(),
					EBS:                       Disabled(),
					FSX:                       Disabled(),
					EFS:                       Disabled(),
					AWSLoadBalancerController: Disabled(),
					XRay:                      Disabled(),
					CloudWatch:                Disabled(),
				},
			},
			ScalingConfig:  &ScalingConfig{},
			SecurityGroups: &NodeGroupSGs{},
		},
	}
}

// NewNodeGroup creates new nodegroup inside cluster config,
// it returns pointer to the nodegroup for convenience
func (c *ClusterConfig) NewNodeGroup() *NodeGroup {
	ng := NewNodeGroup()

	c.NodeGroups = append(c.NodeGroups, ng)

	return ng
}

// NodeGroup holds configuration attributes that are
// specific to a nodegroup
type NodeGroup struct {
	*NodeGroupBase

	//+optional
	InstancesDistribution *NodeGroupInstancesDistribution `json:"instancesDistribution,omitempty"`

	// +optional
	ASGMetricsCollection []MetricsCollection `json:"asgMetricsCollection,omitempty"`

	// CPUCredits configures [T3 Unlimited](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/burstable-performance-instances-unlimited-mode.html), valid only for T-type instances
	// +optional
	CPUCredits *string `json:"cpuCredits,omitempty"`

	// +optional
	Taints map[string]string `json:"taints,omitempty"`

	// Associate load balancers with auto scaling group
	// +optional
	ClassicLoadBalancerNames []string `json:"classicLoadBalancerNames,omitempty"`

	// Associate target group with auto scaling group
	// +optional
	TargetGroupARNs []string `json:"targetGroupARNs,omitempty"`

	// +optional
	Bottlerocket *NodeGroupBottlerocket `json:"bottlerocket,omitempty"`

	// [Custom
	// address](/usage/vpc-networking/#custom-cluster-dns-address) used for DNS
	// lookups
	// +optional
	ClusterDNS string `json:"clusterDNS,omitempty"`

	// [Customize `kubelet` config](/usage/customizing-the-kubelet/)
	// +optional
	KubeletExtraConfig *InlineDocument `json:"kubeletExtraConfig,omitempty"`
}

// BaseNodeGroup implements NodePool
func (n *NodeGroup) BaseNodeGroup() *NodeGroupBase {
	return n.NodeGroupBase
}

// Git groups all configuration options related to enabling GitOps on a
// cluster and linking it to a Git repository.
// [Gitops Guide](/gitops-quickstart/)
type Git struct {

	// [Enable Repo](/usage/gitops/#installing-flux)
	Repo *Repo `json:"repo,omitempty"`

	// [Enable Repo](/usage/gitops/#installing-flux)
	// +optional
	Operator Operator `json:"operator,omitempty"`

	// [Installing a Quickstart profile](/usage/gitops/#installing-a-quickstart-profile-in-your-cluster)
	// +optional
	BootstrapProfile *Profile `json:"bootstrapProfile,omitempty"` // one or many profiles to enable on this cluster once it is created
}

// NewGit returns a new empty Git configuration
func NewGit() *Git {
	return &Git{
		Repo:             &Repo{},
		Operator:         Operator{},
		BootstrapProfile: &Profile{},
	}
}

// Repo groups all configuration options related to a Git repository used for
// GitOps.
type Repo struct {
	// The Git SSH URL to the repository which will contain the cluster configuration
	// For example: `git@github.com:org/repo`
	URL string `json:"url,omitempty"`

	// The git branch under which cluster configuration files will be committed & pushed, e.g. master
	// +optional
	Branch string `json:"branch,omitempty"`

	// Relative paths within the Git repository which the GitOps operator will monitor to find Kubernetes manifests to apply, e.g. ["kube-system", "base"]
	//+optional
	Paths []string `json:"paths,omitempty"`

	// The directory under which Flux configuration files will be written, e.g. flux/
	// +optional
	FluxPath string `json:"fluxPath,omitempty"`

	// Git user which will be used to commit changes
	// +optional
	User string `json:"user,omitempty"`

	// Git email which will be used to commit changes
	Email string `json:"email,omitempty"`

	// Path to the private SSH key to use to authenticate
	// +optional
	PrivateSSHKeyPath string `json:"privateSSHKeyPath,omitempty"`
}

// Operator groups all configuration options related to the operator used to
// keep the cluster and the Git repository in sync.
type Operator struct {

	// Commit and push Flux manifests to the Git Repo on install
	// +optional
	CommitOperatorManifests *bool `json:"commitOperatorManifests,omitempty"`

	// Git label to keep track of Flux's sync progress; this is equivalent to overriding --git-sync-tag and --git-notes-ref in Flux
	// +optional
	Label string `json:"label,omitempty"`

	// Cluster namespace where to install Flux and the Helm Operator e.g. flux
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Install the Helm Operator
	// +optional
	WithHelm *bool `json:"withHelm,omitempty"`

	// Instruct Flux to read-only mode and create the deploy key as read-only
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`

	// Additional command line arguments for the Flux daemon
	// +optional
	AdditionalFluxArgs []string `json:"additionalFluxArgs,omitempty"`

	// Additional command line arguments for the Helm Operator
	// +optional
	AdditionalHelmOperatorArgs []string `json:"additionalHelmOperatorArgs,omitempty"`
}

// Profile groups all details on a quickstart profile to enable on the cluster
// and add to the Git repository.
type Profile struct {

	// Name or URL of the Quick Start profile
	// For example: `app-dev`
	Source string `json:"source,omitempty"`

	// Revision of the Quick Start profile. Can be a branch, tag or commit hash
	// +optional
	Revision string `json:"revision,omitempty"`

	// Output directory for the processed profile templates (generate profile command)
	// Defaults to `./<quickstart-repo-name>`
	// +optional
	OutputPath string `json:"outputPath,omitempty"`
}

// HasBootstrapProfile returns true if there is a profile with a source specified
func (c *ClusterConfig) HasBootstrapProfile() bool {
	return c.Git != nil && c.Git.BootstrapProfile != nil && c.Git.BootstrapProfile.Source != ""
}

// HasGitopsRepoConfigured returns true if there is a profile with a source specified
func (c *ClusterConfig) HasGitopsRepoConfigured() bool {
	return c.Git != nil && c.Git.Repo != nil && c.Git.Repo.URL != ""
}

type (
	// NodeGroupSGs controls security groups for this nodegroup
	NodeGroupSGs struct {
		// AttachIDs attaches additional security groups to the nodegroup
		// +optional
		AttachIDs []string `json:"attachIDs,omitempty"`
		// WithShared attach the security group
		// shared among all nodegroups in the cluster
		// Defaults to `true`
		// +optional
		WithShared *bool `json:"withShared"`
		// WithLocal attach a security group
		// local to this nodegroup
		// Not supported for managed nodegroups
		// Defaults to `true`
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
		InstanceRolePermissionsBoundary string `json:"instanceRolePermissionsBoundary,omitempty"`
		// +optional
		WithAddonPolicies NodeGroupIAMAddonPolicies `json:"withAddonPolicies,omitempty"`
	}
	// NodeGroupIAMAddonPolicies holds all IAM addon policies
	NodeGroupIAMAddonPolicies struct {
		// +optional
		// ImageBuilder allows for full ECR (Elastic Container Registry) access. This is useful for building, for
		// example, a CI server that needs to push images to ECR
		ImageBuilder *bool `json:"imageBuilder"`
		// +optional
		// AutoScaler enables IAM policy for cluster-autoscaler
		AutoScaler *bool `json:"autoScaler"`
		// +optional
		// ExternalDNS adds the external-dns project policies for Amazon Route 53
		ExternalDNS *bool `json:"externalDNS"`
		// +optional
		// CertManager enables the ability to add records to Route 53 in order to solve the DNS01 challenge. More information can be found
		// [here](https://cert-manager.io/docs/configuration/acme/dns01/route53/#set-up-a-iam-role)
		CertManager *bool `json:"certManager"`
		// +optional
		// AppMesh enables full access to AppMesh
		AppMesh *bool `json:"appMesh"`
		// +optional
		// AppMeshPreview enables full access to AppMesh Preview
		AppMeshPreview *bool `json:"appMeshPreview"`
		// +optional
		// EBS enables the new EBS CSI (Elastic Block Store Container Storage Interface) driver
		EBS *bool `json:"ebs"`
		// +optional
		FSX *bool `json:"fsx"`
		// +optional
		EFS *bool `json:"efs"`
		// +optional
		AWSLoadBalancerController *bool `json:"albIngress"`
		// +optional
		XRay *bool `json:"xRay"`
		// +optional
		CloudWatch *bool `json:"cloudWatch"`
	}

	// NodeGroupSSH holds all the ssh access configuration to a NodeGroup
	NodeGroupSSH struct {
		// +optional Enables/Disables the security group configuration. Values provided by SourceSecurityGroupIDs
		// are ignored if set to false
		Allow *bool `json:"allow"`
		// +optional
		PublicKeyPath *string `json:"publicKeyPath,omitempty"`
		// +optional
		PublicKey *string `json:"publicKey,omitempty"`
		// +optional
		PublicKeyName *string `json:"publicKeyName,omitempty"`
		// +optional
		SourceSecurityGroupIDs []string `json:"sourceSecurityGroupIds,omitempty"`
		// Enables the ability to [SSH onto nodes using SSM](/introduction#ssh-access)
		// +optional
		EnableSSM *bool `json:"enableSsm,omitempty"`
	}

	// NodeGroupInstancesDistribution holds the configuration for [spot
	// instances](/usage/spot-instances/)
	NodeGroupInstancesDistribution struct {
		// +required
		InstanceTypes []string `json:"instanceTypes,omitempty"`
		// Defaults to `on demand price`
		// +optional
		MaxPrice *float64 `json:"maxPrice,omitempty"`
		// Defaults to `0`
		// +optional
		OnDemandBaseCapacity *int `json:"onDemandBaseCapacity,omitempty"`
		// Range [0-100]
		// Defaults to `100`
		// +optional
		OnDemandPercentageAboveBaseCapacity *int `json:"onDemandPercentageAboveBaseCapacity,omitempty"`
		// Range [1-20]
		// Defaults to `2`
		// +optional
		SpotInstancePools *int `json:"spotInstancePools,omitempty"`
		// +optional
		SpotAllocationStrategy *string `json:"spotAllocationStrategy,omitempty"`
	}

	// NodeGroupBottlerocket holds the configuration for Bottlerocket based
	// NodeGroups.
	NodeGroupBottlerocket struct {
		// +optional
		EnableAdminContainer *bool `json:"enableAdminContainer,omitempty"`
		// Settings contains any [bottlerocket
		// settings](https://github.com/bottlerocket-os/bottlerocket/#description-of-settings)
		// +optional
		Settings *InlineDocument `json:"settings,omitempty"`
	}
)

// MetricsCollection used by the scaling config,
// see [cloudformation
// docs](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-metricscollection.html)
type MetricsCollection struct {
	// +required
	Granularity string `json:"granularity"`
	// +optional
	Metrics []string `json:"metrics,omitempty"`
}

// ScalingConfig defines the scaling config
type ScalingConfig struct {
	// +optional
	DesiredCapacity *int `json:"desiredCapacity,omitempty"`
	// +optional
	MinSize *int `json:"minSize,omitempty"`
	// +optional
	MaxSize *int `json:"maxSize,omitempty"`
}

// NodePool represents a group of nodes that share the same configuration
// Ideally the NodeGroup type should be renamed to UnmanagedNodeGroup or SelfManagedNodeGroup and this interface
// should be called NodeGroup
type NodePool interface {
	// BaseNodeGroup returns the base nodegroup
	BaseNodeGroup() *NodeGroupBase
}

// NodeGroupBase represents the base nodegroup config for self-managed and managed nodegroups
type NodeGroupBase struct {
	// +required
	Name string `json:"name"`

	// Valid variants are `NodeAMIFamily` constants
	// +optional
	AMIFamily string `json:"amiFamily,omitempty"`
	// +optional
	InstanceType string `json:"instanceType,omitempty"`
	// Limit [nodes to specific
	// AZs](/usage/autoscaling/#zone-aware-auto-scaling)
	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`
	// Limit nodes to specific subnets
	// +optional
	Subnets []string `json:"subnets,omitempty"`

	// +optional
	InstancePrefix string `json:"instancePrefix,omitempty"`
	// +optional
	InstanceName string `json:"instanceName,omitempty"`

	// +optional
	*ScalingConfig

	// +optional
	// VolumeSize gigabytes
	// Defaults to `80`
	VolumeSize *int `json:"volumeSize,omitempty"`
	// +optional
	// SSH configures ssh access for this nodegroup
	SSH *NodeGroupSSH `json:"ssh,omitempty"`
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Enable [private
	// networking](/usage/vpc-networking/#use-private-subnets-for-initial-nodegroup)
	// for nodegroup
	// +optional
	PrivateNetworking bool `json:"privateNetworking"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// +optional
	IAM *NodeGroupIAM `json:"iam,omitempty"`

	// Specify [custom AMIs](/usage/custom-ami-support/), `auto-ssm`, `auto`, or `static`
	// +optional
	AMI string `json:"ami,omitempty"`

	// +optional
	SecurityGroups *NodeGroupSGs `json:"securityGroups,omitempty"`

	// +optional
	MaxPodsPerNode int `json:"maxPodsPerNode,omitempty"`

	// See [relevant AWS
	// docs](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-attribute-updatepolicy.html#cfn-attributes-updatepolicy-rollingupdate-suspendprocesses)
	// +optional
	ASGSuspendProcesses []string `json:"asgSuspendProcesses,omitempty"`

	// EBSOptimized enables [EBS
	// optimization](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-optimized.html)
	// +optional
	EBSOptimized *bool `json:"ebsOptimized,omitempty"`

	// Valid variants are `VolumeType` constants
	// +optional
	VolumeType *string `json:"volumeType,omitempty"`
	// +optional
	VolumeName *string `json:"volumeName,omitempty"`
	// Some AMIs (bottlerocket) have a separate volume for the OS
	OSVolumeName *string
	// +optional
	VolumeEncrypted *bool `json:"volumeEncrypted,omitempty"`
	// +optional
	VolumeKmsKeyID *string `json:"volumeKmsKeyID,omitempty"`
	// +optional
	VolumeIOPS *int `json:"volumeIOPS,omitempty"`

	// PreBootstrapCommands are executed before bootstrapping instances to the
	// cluster
	// +optional
	PreBootstrapCommands []string `json:"preBootstrapCommands,omitempty"`

	// Override `eksctl`'s bootstrapping script
	// +optional
	OverrideBootstrapCommand *string `json:"overrideBootstrapCommand,omitempty"`

	// DisableIMDSv1 requires requests to the metadata service to use IMDSv2 tokens
	// Defaults to `false`
	// +optional
	DisableIMDSv1 *bool `json:"disableIMDSv1,omitempty"`

	// DisablePodIMDS blocks all IMDS requests from non host networking pods
	// Defaults to `false`
	// +optional
	DisablePodIMDS *bool `json:"disablePodIMDS,omitempty"`

	// Placement specifies the placement group in which nodes should
	// be spawned
	// +optional
	Placement *Placement `json:"placement,omitempty"`
}

// Placement specifies placement group information
type Placement struct {
	GroupName string `json:"groupName,omitempty"`
}

// ListOptions returns metav1.ListOptions with label selector for the nodegroup
func (n *NodeGroupBase) ListOptions() metav1.ListOptions {
	return makeListOptions(n.Name)
}

// NameString returns the nodegroup name
func (n *NodeGroupBase) NameString() string {
	return n.Name
}

// Size returns the minimum nodegroup size
func (n *NodeGroupBase) Size() int {
	if n.MinSize == nil {
		return 0
	}
	return *n.MinSize
}

// GetAMIFamily returns the AMI family
func (n *NodeGroupBase) GetAMIFamily() string {
	return n.AMIFamily
}

type LaunchTemplate struct {
	// Launch template ID
	// +required
	ID string `json:"id,omitempty"`
	// Launch template version
	// Defaults to the default launch template version
	// TODO support $Default, $Latest
	Version *string `json:"version,omitempty"`
	// TODO support Name?
}

// ManagedNodeGroup represents an EKS-managed nodegroup
// TODO Validate for unmapped fields and throw an error
type ManagedNodeGroup struct {
	*NodeGroupBase

	// InstanceTypes specifies a list of instance types
	InstanceTypes []string `json:"instanceTypes,omitempty"`

	// Spot creates a spot nodegroup
	Spot bool `json:"spot,omitempty"`

	// LaunchTemplate specifies an existing launch template to use
	// for the nodegroup
	LaunchTemplate *LaunchTemplate `json:"launchTemplate,omitempty"`
}

// BaseNodeGroup implements NodePool
func (m *ManagedNodeGroup) BaseNodeGroup() *NodeGroupBase {
	return m.NodeGroupBase
}

func makeListOptions(nodeGroupName string) metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", NodeGroupNameLabel, nodeGroupName),
	}
}

// InlineDocument holds any arbitrary JSON/YAML documents, such as extra config parameters or IAM policies
type InlineDocument map[string]interface{}

// DeepCopy is needed to generate kubernetes types for InlineDocument
func (in *InlineDocument) DeepCopy() *InlineDocument {
	if in == nil {
		return nil
	}
	out := new(InlineDocument)
	*out = runtime.DeepCopyJSON(*in)
	return out
}

// HasMixedInstances checks if a nodegroup has mixed instances option declared
func HasMixedInstances(ng *NodeGroup) bool {
	return ng.InstancesDistribution != nil && len(ng.InstancesDistribution.InstanceTypes) > 0
}

// IsAMI returns true if the argument is an AMI ID
func IsAMI(amiFlag string) bool {
	return strings.HasPrefix(amiFlag, "ami-")
}

// FargateProfile defines the settings used to schedule workload onto Fargate.
type FargateProfile struct {

	// Name of the Fargate profile.
	// +required
	Name string `json:"name"`

	// PodExecutionRoleARN is the IAM role's ARN to use to run pods onto Fargate.
	PodExecutionRoleARN string `json:"podExecutionRoleARN,omitempty"`

	// Selectors define the rules to select workload to schedule onto Fargate.
	Selectors []FargateProfileSelector `json:"selectors"`

	// Subnets which Fargate should use to do network placement of the selected workload.
	// If none provided, all subnets for the cluster will be used.
	// +optional
	Subnets []string `json:"subnets,omitempty"`

	// Used to tag the AWS resources
	// +optional
	Tags map[string]string `json:"tags,omitempty"`

	// The current status of the Fargate profile.
	Status string `json:"status"`
}

// FargateProfileSelector defines rules to select workload to schedule onto Fargate.
type FargateProfileSelector struct {

	// Namespace is the Kubernetes namespace from which to select workload.
	// +required
	Namespace string `json:"namespace"`

	// Labels are the Kubernetes label selectors to use to select workload.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// SecretsEncryption defines the configuration for KMS encryption provider
type SecretsEncryption struct {
	// +required
	KeyARN *string `json:"keyARN,omitempty"`
}

// PrivateCluster defines the configuration for a fully-private cluster
type PrivateCluster struct {

	// Enabled enables creation of a fully-private cluster
	Enabled bool `json:"enabled"`

	// AdditionalEndpointServices specifies additional endpoint services that
	// must be enabled for private access.
	// Valid entries are `AdditionalEndpointServices` constants
	AdditionalEndpointServices []string `json:"additionalEndpointServices,omitempty"`
}
