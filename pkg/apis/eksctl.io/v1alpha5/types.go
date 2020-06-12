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

const (
	// AWSDebugLevel defines the LogLevel for AWS produced logs
	AWSDebugLevel = 5

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

	// Version1_10 represents Kubernetes version 1.10.x
	Version1_10 = "1.10"

	// Version1_11 represents Kubernetes version 1.11.x
	Version1_11 = "1.11"

	// Version1_12 represents Kubernetes version 1.12.x
	Version1_12 = "1.12"

	// Version1_13 represents Kubernetes version 1.13.x
	Version1_13 = "1.13"

	// Version1_14 represents Kubernetes version 1.14.x
	Version1_14 = "1.14"

	// Version1_15 represents Kubernetes version 1.15.x
	Version1_15 = "1.15"

	// Version1_16 represents Kubernetes version 1.16.x
	Version1_16 = "1.16"

	// DefaultVersion represents default Kubernetes version supported by EKS
	DefaultVersion = Version1_16

	// LatestVersion represents latest Kubernetes version supported by EKS
	LatestVersion = Version1_16

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
	// NodeImageFamilyBottlerocket represents Bottlerocket family
	NodeImageFamilyBottlerocket = "Bottlerocket"

	// NodeImageFamilyWindowsServer2019CoreContainer represents Windows 2019 core container family
	NodeImageFamilyWindowsServer2019CoreContainer = "WindowsServer2019CoreContainer"

	// NodeImageFamilyWindowsServer2019FullContainer represents Windows 2019 full container family
	NodeImageFamilyWindowsServer2019FullContainer = "WindowsServer2019FullContainer"

	// NodeImageResolverStatic represents static AMI resolver (see ami package)
	NodeImageResolverStatic = "static"
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

	// IAMServiceAccountNameTag defines the tag of the iamserviceaccount name
	IAMServiceAccountNameTag = "alpha.eksctl.io/iamserviceaccount-name"

	// ClusterNameLabel defines the tag of the cluster name
	ClusterNameLabel = "alpha.eksctl.io/cluster-name"

	// NodeGroupNameLabel defines the label of the nodegroup name
	NodeGroupNameLabel = "alpha.eksctl.io/nodegroup-name"

	// ClusterHighlyAvailableNAT defines the highly available NAT configuration option
	ClusterHighlyAvailableNAT = "HighlyAvailable"

	// ClusterSingleNAT defines the single NAT configuration option
	ClusterSingleNAT = "Single"

	// ClusterDisableNAT defines the disabled NAT configuration option
	ClusterDisableNAT = "Disable"

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

	// eksResourceAccountUSGovWest1 defines the AWS EKS account ID that provides node resources in us-gov-west-1
	eksResourceAccountUSGovWest1 = "013241004608"

	// eksResourceAccountUSGovEast1 defines the AWS EKS account ID that provides node resources in us-gov-east-1
	eksResourceAccountUSGovEast1 = "151742754352"
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
		RegionAPNorthEast1,
		RegionAPNorthEast2,
		RegionAPSouthEast1,
		RegionAPSouthEast2,
		RegionAPSouth1,
		RegionAPEast1,
		RegionMESouth1,
		RegionSAEast1,
		RegionCNNorthwest1,
		RegionCNNorth1,
		RegionUSGovWest1,
		RegionUSGovEast1,
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
	}
}

// SupportedVersions are the versions of Kubernetes that EKS supports
func SupportedVersions() []string {
	return []string{
		Version1_13,
		Version1_14,
		Version1_15,
		Version1_16,
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
	default:
		return eksResourceAccountStandard
	}
}

// ClusterMeta is what identifies a cluster
type ClusterMeta struct {
	Name   string `json:"name" jsonschema:"required"`
	Region string `json:"region" jsonschema:"required"`
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
	CloudFormationRoleARN string

	Region      string
	Profile     string
	WaitTimeout time.Duration
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfig is a simple config, to be replaced with Cluster API
type ClusterConfig struct {
	metav1.TypeMeta

	Metadata *ClusterMeta `json:"metadata" jsonschema:"required"`

	// +optional
	IAM *ClusterIAM `json:"iam,omitempty"`

	// +optional
	VPC *ClusterVPC `json:"vpc,omitempty"`

	// +optional
	NodeGroups []*NodeGroup `json:"nodeGroups,omitempty"`

	// +optional
	ManagedNodeGroups []*ManagedNodeGroup `json:"managedNodeGroups,omitempty"`

	// +optional
	FargateProfiles []*FargateProfile `json:"fargateProfiles,omitempty"`

	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

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
		IAM: &ClusterIAM{},
		VPC: NewClusterVPC(),
		CloudWatch: &ClusterCloudWatch{
			ClusterLogging: &ClusterCloudWatchLogging{},
		},
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
		ClusterEndpoints: ClusterEndpointAccessDefaults(),
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

// NewNodeGroup creates new nodegroup, and returns pointer to it
func NewNodeGroup() *NodeGroup {
	return &NodeGroup{
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
				ImageBuilder:   Disabled(),
				AutoScaler:     Disabled(),
				ExternalDNS:    Disabled(),
				CertManager:    Disabled(),
				AppMesh:        Disabled(),
				AppMeshPreview: Disabled(),
				EBS:            Disabled(),
				FSX:            Disabled(),
				EFS:            Disabled(),
				ALBIngress:     Disabled(),
				XRay:           Disabled(),
				CloudWatch:     Disabled(),
			},
		},
		SSH: &NodeGroupSSH{
			Allow:         Disabled(),
			PublicKeyPath: &DefaultNodeSSHPublicKeyPath,
		},
	}
}

// NewManagedNodeGroup creates a new ManagedNodeGroup
func NewManagedNodeGroup() *ManagedNodeGroup {
	var (
		publicKey  = DefaultNodeSSHPublicKeyPath
		volumeSize = DefaultNodeVolumeSize
	)
	return &ManagedNodeGroup{
		VolumeSize:    &volumeSize,
		ScalingConfig: &ScalingConfig{},
		SSH: &NodeGroupSSH{
			Allow:         Disabled(),
			PublicKeyName: &publicKey,
		},
		IAM: &NodeGroupIAM{
			WithAddonPolicies: NodeGroupIAMAddonPolicies{
				ImageBuilder:   Disabled(),
				AutoScaler:     Disabled(),
				ExternalDNS:    Disabled(),
				CertManager:    Disabled(),
				AppMesh:        Disabled(),
				AppMeshPreview: Disabled(),
				EBS:            Disabled(),
				FSX:            Disabled(),
				EFS:            Disabled(),
				ALBIngress:     Disabled(),
				XRay:           Disabled(),
				CloudWatch:     Disabled(),
			},
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

// NodeGroup holds all configuration attributes that are
// specific to a nodegroup
type NodeGroup struct {
	Name string `json:"name" jsonschema:"required"`
	// +optional
	AMI string `json:"ami,omitempty"`
	// +optional
	AMIFamily string `json:"amiFamily,omitempty"`
	// +optional
	InstanceType string `json:"instanceType,omitempty"`
	//+optional
	InstancesDistribution *NodeGroupInstancesDistribution `json:"instancesDistribution,omitempty"`
	// +optional
	InstancePrefix string `json:"instancePrefix,omitempty"`
	// +optional
	InstanceName string `json:"instanceName,omitempty"`
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
	ASGMetricsCollection []MetricsCollection `json:"asgMetricsCollection,omitempty"`

	// +optional
	EBSOptimized *bool `json:"ebsOptimized,omitempty"`

	// +optional
	VolumeSize *int `json:"volumeSize,omitempty"`
	// +optional
	VolumeType *string `json:"volumeType,omitempty"`
	// +optional
	VolumeName *string `json:"volumeName,omitempty"`
	// +optional
	VolumeEncrypted *bool `json:"volumeEncrypted,omitempty"`
	// +optional
	VolumeKmsKeyID *string `json:"volumeKmsKeyID,omitempty"`
	// +optional
	VolumeIOPS *int `json:"volumeIOPS,omitempty"`

	// +optional
	MaxPodsPerNode int `json:"maxPodsPerNode,omitempty"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// +optional
	Taints map[string]string `json:"taints,omitempty"`

	// +optional
	ClassicLoadBalancerNames []string `json:"classicLoadBalancerNames,omitempty"`

	// +optional
	TargetGroupARNs []string `json:"targetGroupARNs,omitempty"`

	// +optional
	SSH *NodeGroupSSH `json:"ssh,omitempty"`

	// +optional
	IAM *NodeGroupIAM `json:"iam,omitempty"`

	// +optional
	Bottlerocket *NodeGroupBottlerocket `json:"bottlerocket,omitempty"`

	// +optional
	PreBootstrapCommands []string `json:"preBootstrapCommands,omitempty"`

	// +optional
	OverrideBootstrapCommand *string `json:"overrideBootstrapCommand,omitempty"`

	// +optional
	ClusterDNS string `json:"clusterDNS,omitempty"`

	// +optional
	KubeletExtraConfig *InlineDocument `json:"kubeletExtraConfig,omitempty"`
}

// Git groups all configuration options related to enabling GitOps on a
// cluster and linking it to a Git repository.
type Git struct {
	Repo *Repo `json:"repo,omitempty"`
	// +optional
	Operator Operator `json:"operator,omitempty"`
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
	URL string `json:"url,omitempty"` // the Git SSH URL to the repository which will contain the cluster configuration, e.g. git@github.com:org/repo
	// +optional
	Branch string `json:"branch,omitempty"` // the branch under which cluster configuration files will be committed & pushed, e.g. master
	// +optional
	Paths []string `json:"paths,omitempty"` // relative paths within the Git repository which the GitOps operator will monitor to find Kubernetes manifests to apply, e.g. ["kube-system", "base"]
	// +optional
	FluxPath string `json:"fluxPath,omitempty"` // the directory under which Flux configuration files will be written, e.g. flux/
	// +optional
	User  string `json:"user,omitempty"`  // Git user which will be used to commit changes
	Email string `json:"email,omitempty"` // Git email which will be used to commit changes
	// +optional
	PrivateSSHKeyPath string `json:"privateSSHKeyPath,omitempty"` // path to the private SSH key to use to authenticate
}

// Operator groups all configuration options related to the operator used to
// keep the cluster and the Git repository in sync.
type Operator struct {
	// +optional
	CommitOperatorManifests *bool `json:"commitOperatorManifests,omitempty"` // Commit and push Flux manifests to the Git Repo on install
	// +optional
	Label string `json:"label,omitempty"` // e.g. flux
	// +optional
	Namespace string `json:"namespace,omitempty"` // e.g. flux
	// +optional
	WithHelm *bool `json:"withHelm,omitempty"` // whether to install the Flux Helm Operator or not
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"` // Instruct Flux to read-only mode and create the deploy key as read-only
}

// Profile groups all details on a quickstart profile to enable on the cluster
// and add to the Git repository.
type Profile struct {
	Source string `json:"source,omitempty"` // e.g. app-dev
	// +optional
	Revision string `json:"revision,omitempty"` // branch, tag or commit hash
	// +optional
	OutputPath string `json:"outputPath,omitempty"` // output directory for processed profile templates (generate profile command)
}

// HasBootstrapProfile returns true if there is a profile with a source specified
func (c *ClusterConfig) HasBootstrapProfile() bool {
	return c.Git != nil && c.Git.BootstrapProfile != nil && c.Git.BootstrapProfile.Source != ""
}

// HasGitopsRepoConfigured returns true if there is a profile with a source specified
func (c *ClusterConfig) HasGitopsRepoConfigured() bool {
	return c.Git != nil && c.Git.Repo != nil && c.Git.Repo.URL != ""
}

// ListOptions returns metav1.ListOptions with label selector for the nodegroup
func (n *NodeGroup) ListOptions() metav1.ListOptions {
	return makeListOptions(n.Name)
}

// NameString returns common name string
func (n *NodeGroup) NameString() string {
	return n.Name
}

// Size returns the minimum nodegroup size
func (n *NodeGroup) Size() int {
	if n.MinSize == nil {
		return 0
	}
	return *n.MinSize
}

// GetAMIFamily returns the AMI family
func (n *NodeGroup) GetAMIFamily() string {
	return n.AMIFamily
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
		InstanceRolePermissionsBoundary string `json:"instanceRolePermissionsBoundary,omitempty"`
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
		AppMeshPreview *bool `json:"appMeshPreview"`
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
		// +optional
		SourceSecurityGroupIDs []string `json:"sourceSecurityGroupIds,omitempty"`
	}

	// NodeGroupInstancesDistribution holds the configuration for spot instances
	NodeGroupInstancesDistribution struct {
		//+required
		InstanceTypes []string `json:"instanceTypes,omitempty" jsonschema:"required"`
		// +optional
		MaxPrice *float64 `json:"maxPrice,omitempty"`
		//+optional
		OnDemandBaseCapacity *int `json:"onDemandBaseCapacity,omitempty"`
		//+optional
		OnDemandPercentageAboveBaseCapacity *int `json:"onDemandPercentageAboveBaseCapacity,omitempty"`
		//+optional
		SpotInstancePools *int `json:"spotInstancePools,omitempty"`
		//+optional
		SpotAllocationStrategy *string `json:"spotAllocationStrategy,omitempty"`
	}

	// NodeGroupBottlerocket holds the configuration for Bottlerocket based
	// NodeGroups.
	NodeGroupBottlerocket struct {
		// +optional
		EnableAdminContainer *bool `json:"enableAdminContainer,omitempty"`
		// +optional
		Settings *InlineDocument `json:"settings,omitempty"`
	}
)

// MetricsCollection used by the scaling config,
// see [cloudformation
// docs](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-as-metricscollection.html)
type MetricsCollection struct {
	// +required
	Granularity string `json:"granularity" jsonschema:"required"`
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

// ManagedNodeGroup defines an EKS-managed nodegroup
// TODO Validate for unmapped fields and throw an error
type ManagedNodeGroup struct {
	Name string `json:"name" jsonschema:"required"`

	// +optional
	AMIFamily string `json:"amiFamily,omitempty"`
	// +optional
	InstanceType string `json:"instanceType,omitempty"`
	// +optional
	*ScalingConfig
	// +optional
	VolumeSize *int `json:"volumeSize,omitempty"`
	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`
	// +optional
	SSH *NodeGroupSSH `json:"ssh,omitempty"`

	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// +optional
	PrivateNetworking bool `json:"privateNetworking"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// +optional
	IAM *NodeGroupIAM `json:"iam,omitempty"`
}

// ListOptions returns metav1.ListOptions with label selector for the managed nodegroup
func (n *ManagedNodeGroup) ListOptions() metav1.ListOptions {
	return makeListOptions(n.Name)
}

// NameString returns the nodegroup name
func (n *ManagedNodeGroup) NameString() string {
	return n.Name
}

// Size returns the minimum nodegroup size
func (n *ManagedNodeGroup) Size() int {
	if n.MinSize == nil {
		return 0
	}
	return *n.MinSize
}

// GetAMIFamily returns the AMI family
func (n *ManagedNodeGroup) GetAMIFamily() string {
	return n.AMIFamily
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
	Name string `json:"name" jsonschema:"required"`
	// PodExecutionRoleARN is the IAM role's ARN to use to run pods onto Fargate.
	PodExecutionRoleARN string `json:"podExecutionRoleARN,omitempty"`
	// Selectors define the rules to select workload to schedule onto Fargate.
	Selectors []FargateProfileSelector `json:"selectors"`
	// +optional
	// Subnets which Fargate should use to do network placement of the selected workload.
	// If none provided, all subnets for the cluster will be used.
	Subnets []string `json:"subnets,omitempty"`

	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// FargateProfileSelector defines rules to select workload to schedule onto Fargate.
type FargateProfileSelector struct {
	// Namespace is the Kubernetes namespace from which to select workload.
	Namespace string `json:"namespace" jsonschema:"required"`
	// +optional
	// Labels are the Kubernetes label selectors to use to select workload.
	Labels map[string]string `json:"labels,omitempty"`
}

// SecretsEncryption defines the configuration for KMS encryption provider
type SecretsEncryption struct {
	KeyARN *string `json:"keyARN,omitempty" jsonschema:"required"`
}
