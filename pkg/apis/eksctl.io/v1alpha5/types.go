package v1alpha5

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/taints"
)

// Values for `KubernetesVersion`
// All valid values should go in this block
const (
	Version1_23 = "1.23"

	Version1_24 = "1.24"

	Version1_25 = "1.25"

	Version1_26 = "1.26"

	Version1_27 = "1.27"

	Version1_28 = "1.28"

	// DefaultVersion (default)
	DefaultVersion = Version1_27

	LatestVersion = Version1_28

	DockershimDeprecationVersion = Version1_24
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

	// Version1_14 represents Kubernetes version 1.14.x
	Version1_14 = "1.14"

	// Version1_15 represents Kubernetes version 1.15.x
	Version1_15 = "1.15"

	// Version1_16 represents Kubernetes version 1.16.x
	Version1_16 = "1.16"

	// Version1_17 represents Kubernetes version 1.17.x
	Version1_17 = "1.17"

	// Version1_18 represents Kubernetes version 1.18.x
	Version1_18 = "1.18"

	// Version1_19 represents Kubernetes version 1.19.x
	Version1_19 = "1.19"

	// Version1_20 represents Kubernetes version 1.20.x
	Version1_20 = "1.20"

	// Version1_21 represents Kubernetes version 1.21.x
	Version1_21 = "1.21"

	// Version1_22 represents Kubernetes version 1.22.x
	Version1_22 = "1.22"
)

// Not yet supported versions
const (
	// Version1_29 represents Kubernetes version 1.29.x
	Version1_29 = "1.29"
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

	// RegionCAWest1 represents the Canada West region Calgary.
	RegionCAWest1 = "ca-west-1"

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

	// RegionEUCentral2 represents the EU Central Region Zurich.
	RegionEUCentral2 = "eu-central-2"

	// RegionEUSouth1 represents the Eu South Region Milan
	RegionEUSouth1 = "eu-south-1"

	// RegionEUSouth2 represents the Eu South Region Spain
	RegionEUSouth2 = "eu-south-2"

	// RegionAPNorthEast1 represents the Asia-Pacific North East Region Tokyo
	RegionAPNorthEast1 = "ap-northeast-1"

	// RegionAPNorthEast2 represents the Asia-Pacific North East Region Seoul
	RegionAPNorthEast2 = "ap-northeast-2"

	// RegionAPNorthEast3 represents the Asia-Pacific North East region Osaka
	RegionAPNorthEast3 = "ap-northeast-3"

	// RegionAPSouthEast1 represents the Asia-Pacific South East Region Singapore
	RegionAPSouthEast1 = "ap-southeast-1"

	// RegionAPSouthEast2 represents the Asia-Pacific South East Region Sydney
	RegionAPSouthEast2 = "ap-southeast-2"

	// RegionAPSouthEast3 represents the Asia-Pacific South East Region Jakarta
	RegionAPSouthEast3 = "ap-southeast-3"

	// RegionAPSouthEast4 represents the Asia-Pacific South East Region Melbourne
	RegionAPSouthEast4 = "ap-southeast-4"

	// RegionAPSouth1 represents the Asia-Pacific South Region Mumbai
	RegionAPSouth1 = "ap-south-1"

	// RegionAPSouth2 represents the Asia-Pacific South Region Hyderabad
	RegionAPSouth2 = "ap-south-2"

	// RegionAPEast1 represents the Asia Pacific Region Hong Kong
	RegionAPEast1 = "ap-east-1"

	// RegionMECentral1 represents the Middle East Region Dubai
	RegionMECentral1 = "me-central-1"

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

	// RegionILCentral1 represents the Israel region Tel Aviv
	RegionILCentral1 = "il-central-1"

	// RegionUSGovWest1 represents the region GovCloud (US-West)
	RegionUSGovWest1 = "us-gov-west-1"

	// RegionUSGovEast1 represents the region GovCloud (US-East)
	RegionUSGovEast1 = "us-gov-east-1"

	// RegionUSISOEast1 represents the region US ISO East.
	RegionUSISOEast1 = "us-iso-east-1"

	// RegionUSISOBEast1 represents the region US ISOB East (Ohio).
	RegionUSISOBEast1 = "us-isob-east-1"

	// RegionUSISOWest1 represents the region US ISOB West.
	RegionUSISOWest1 = "us-iso-west-1"

	// DefaultRegion defines the default region, where to deploy the EKS cluster
	DefaultRegion = RegionUSWest2
)

// Values for `NodeAMIFamily`
// All valid values of supported families should go in this block
const (
	// DefaultNodeImageFamily (default)
	DefaultNodeImageFamily      = NodeImageFamilyAmazonLinux2
	NodeImageFamilyAmazonLinux2 = "AmazonLinux2"
	NodeImageFamilyUbuntu2004   = "Ubuntu2004"
	NodeImageFamilyUbuntu1804   = "Ubuntu1804"
	NodeImageFamilyBottlerocket = "Bottlerocket"

	NodeImageFamilyWindowsServer2019CoreContainer = "WindowsServer2019CoreContainer"
	NodeImageFamilyWindowsServer2019FullContainer = "WindowsServer2019FullContainer"

	NodeImageFamilyWindowsServer2022CoreContainer = "WindowsServer2022CoreContainer"
	NodeImageFamilyWindowsServer2022FullContainer = "WindowsServer2022FullContainer"
)

// Deprecated `NodeAMIFamily`
const (
	NodeImageFamilyWindowsServer2004CoreContainer = "WindowsServer2004CoreContainer"
	NodeImageFamilyWindowsServer20H2CoreContainer = "WindowsServer20H2CoreContainer"
)

// Container runtime values.
const (
	ContainerRuntimeContainerD       = "containerd"
	ContainerRuntimeDockerD          = "dockerd"
	ContainerRuntimeDockerForWindows = "docker"
)

const (
	// DefaultNodeType is the default instance type to use for nodes
	DefaultNodeType = "m5.large"

	// DefaultNodeCount defines the default number of nodes to be created
	DefaultNodeCount = 2

	// DefaultMaxSize defines the default maximum number of nodes inside the ASG
	DefaultMaxSize = 1

	// NodeImageResolverAuto represents auto AMI resolver (see ami package)
	NodeImageResolverAuto = "auto"
	// NodeImageResolverAutoSSM is used to indicate that the latest EKS AMIs should be used for the nodes. The AMI is selected
	// using an SSM GetParameter query
	NodeImageResolverAutoSSM = "auto-ssm"

	// EksctlVersionTag defines the version of eksctl which is used to provision or update EKS cluster
	EksctlVersionTag = "alpha.eksctl.io/eksctl-version"

	// ClusterNameTag defines the tag of the cluster name
	ClusterNameTag = "alpha.eksctl.io/cluster-name"

	// ClusterOIDCEnabledTag determines whether OIDC is enabled or not.
	ClusterOIDCEnabledTag = "alpha.eksctl.io/cluster-oidc-enabled"

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

	// PodIdentityAssociationNameTag defines the tag of Pod Identity Association name
	PodIdentityAssociationNameTag = "alpha.eksctl.io/podidentityassociation-name"

	// AddonNameTag defines the tag of the IAM service account name
	AddonNameTag = "alpha.eksctl.io/addon-name"

	// ClusterNameLabel defines the tag of the cluster name
	ClusterNameLabel = "alpha.eksctl.io/cluster-name"

	// NodeGroupNameLabel defines the label of the nodegroup name
	NodeGroupNameLabel = "alpha.eksctl.io/nodegroup-name"

	// KarpenterNameTag defines the tag of the Karpenter stack name
	KarpenterNameTag = "alpha.eksctl.io/karpenter-name"

	// KarpenterVersionTag defines the tag for Karpenter's version
	KarpenterVersionTag = "alpha.eksctl.io/karpenter-version"

	EKSNodeGroupNameLabel = "eks.amazonaws.com/nodegroup"

	// SpotAllocationStrategyLowestPrice defines the ASG spot allocation strategy of lowest-price
	SpotAllocationStrategyLowestPrice = "lowest-price"

	// SpotAllocationStrategyCapacityOptimized defines the ASG spot allocation strategy of capacity-optimized
	SpotAllocationStrategyCapacityOptimized = "capacity-optimized"

	// SpotAllocationStrategyCapacityOptimizedPrioritized defines the ASG spot allocation strategy of capacity-optimized-prioritized
	// Use the capacity-optimized-prioritized allocation strategy and then set the order of instance types in
	// the list of launch template overrides from highest to lowest priority (first to last in the list).
	// Amazon EC2 Auto Scaling honors the instance type priorities on a best-effort basis but optimizes
	// for capacity first. This is a good option for workloads where the possibility of disruption must be
	// minimized, but also the preference for certain instance types matters.
	// https://docs.aws.amazon.com/autoscaling/ec2/userguide/asg-purchase-options.html#asg-spot-strategy
	SpotAllocationStrategyCapacityOptimizedPrioritized = "capacity-optimized-prioritized"

	// eksResourceAccountStandard defines the AWS EKS account ID that provides node resources in default regions
	// for standard AWS partition
	eksResourceAccountStandard = "602401143452"

	// eksResourceAccountAPEast1 defines the AWS EKS account ID that provides node resources in ap-east-1 region
	eksResourceAccountAPEast1 = "800184023465"

	// eksResourceAccountCAWest1 defines the AWS EKS account ID that provides node resources in ca-west-1 region
	eksResourceAccountCAWest1 = "761377655185"

	// eksResourceAccountMECentral1 defines the AWS EKS account ID that provides node resources in me-central-1 region
	eksResourceAccountMECentral1 = "759879836304"

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

	// eksResourceAccountEUSouth2 defines the AWS EKS account ID that provides node resources in eu-south-2
	eksResourceAccountEUSouth2 = "455263428931"

	// eksResourceAccountEUCentral2 defines the AWS EKS account ID that provides node resources in eu-central-2.
	eksResourceAccountEUCentral2 = "900612956339"

	// eksResourceAccountUSGovWest1 defines the AWS EKS account ID that provides node resources in us-gov-west-1
	eksResourceAccountUSGovWest1 = "013241004608"

	// eksResourceAccountUSGovEast1 defines the AWS EKS account ID that provides node resources in us-gov-east-1
	eksResourceAccountUSGovEast1 = "151742754352"

	// eksResourceAccountAPSouth2 defines the AWS EKS account ID that provides node resources in ap-south-2
	eksResourceAccountAPSouth2 = "900889452093"

	// eksResourceAccountAPSouthEast3 defines the AWS EKS account ID that provides node resources in ap-southeast-3
	eksResourceAccountAPSouthEast3 = "296578399912"

	// eksResourceAccountILCentral1 defines the AWS EKS account ID that provides node resources in il-central-1
	eksResourceAccountILCentral1 = "066635153087"

	// eksResourceAccountAPSouthEast4 defines the AWS EKS account ID that provides node resources in ap-southeast-4
	eksResourceAccountAPSouthEast4 = "491585149902"
	// eksResourceAccountUSISOEast1 defines the AWS EKS account ID that provides node resources in us-iso-east-1
	eksResourceAccountUSISOEast1 = "725322719131"

	// eksResourceAccountUSISOBEast1 defines the AWS EKS account ID that provides node resources in us-isob-east-1
	eksResourceAccountUSISOBEast1 = "187977181151"

	// eksResourceAccountUSISOWest1 defines the AWS EKS account ID that provides node resources in us-iso-west-1
	eksResourceAccountUSISOWest1 = "608367168043"
)

// Values for `VolumeType`
const (
	// NodeVolumeTypeGP2 is General Purpose SSD
	NodeVolumeTypeGP2 = "gp2"
	// NodeVolumeTypeGP3 is General Purpose SSD which can be optimised for high throughput (default)
	NodeVolumeTypeGP3 = "gp3"
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
	// NodeGroupTypeUnowned defines an unowned managed nodegroup
	NodeGroupTypeUnowned NodeGroupType = "unowned"
	// DefaultNodeVolumeThroughput defines the default throughput for gp3 volumes, set to the min value
	DefaultNodeVolumeThroughput = 125
	// DefaultNodeVolumeIO1IOPS defines the default throughput for io1 volumes, set to the min value
	DefaultNodeVolumeIO1IOPS = 100
	// DefaultNodeVolumeGP3IOPS defines the default throughput for gp3, set to the min value
	DefaultNodeVolumeGP3IOPS = 3000
)

// Values for `IPFamily`
const (
	// IPV4Family defines an IP family of v4 to be used when creating a new VPC and cluster.
	IPV4Family = "IPv4"
	// IPV6Family defines an IP family of v6 to be used when creating a new VPC and cluster.
	IPV6Family = "IPv6"
)

// Values for core addons
const (
	minimumVPCCNIVersionForIPv6 = "1.10.0"
	VPCCNIAddon                 = "vpc-cni"
	KubeProxyAddon              = "kube-proxy"
	CoreDNSAddon                = "coredns"
	PodIdentityAgentAddon       = "eks-pod-identity-agent"
	AWSEBSCSIDriverAddon        = "aws-ebs-csi-driver"
	AWSEFSCSIDriverAddon        = "aws-efs-csi-driver"
)

// supported version of Karpenter
const (
	supportedKarpenterVersion = "v0.20.0"
)

// Values for Capacity Reservation Preference
const (
	OpenCapacityReservation = "open"
	NoneCapacityReservation = "none"
)

var (
	// DefaultIPFamily defines the default IP family to use when creating a new VPC and cluster.
	DefaultIPFamily = IPV4Family
)

var (
	// DefaultWaitTimeout defines the default wait timeout
	DefaultWaitTimeout = 25 * time.Minute

	// DefaultNodeSSHPublicKeyPath is the default path to SSH public key
	DefaultNodeSSHPublicKeyPath = "~/.ssh/id_rsa.pub"

	// DefaultNodeVolumeType defines the default root volume type to use for
	// non-Outpost clusters.
	DefaultNodeVolumeType = NodeVolumeTypeGP3

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

// IsEmpty will only return true if s is not nil and not empty
func IsEmpty(s *string) bool { return !IsSetAndNonEmptyString(s) }

// SupportedRegions are the regions where EKS is available
func SupportedRegions() []string {
	return []string{
		RegionUSWest1,
		RegionUSWest2,
		RegionUSEast1,
		RegionUSEast2,
		RegionCACentral1,
		RegionCAWest1,
		RegionEUWest1,
		RegionEUWest2,
		RegionEUWest3,
		RegionEUNorth1,
		RegionEUCentral1,
		RegionEUCentral2,
		RegionEUSouth1,
		RegionEUSouth2,
		RegionAPNorthEast1,
		RegionAPNorthEast2,
		RegionAPNorthEast3,
		RegionAPSouthEast1,
		RegionAPSouthEast2,
		RegionAPSouthEast3,
		RegionAPSouthEast4,
		RegionAPSouth1,
		RegionAPSouth2,
		RegionAPEast1,
		RegionMECentral1,
		RegionMESouth1,
		RegionSAEast1,
		RegionAFSouth1,
		RegionCNNorthwest1,
		RegionCNNorth1,
		RegionILCentral1,
		RegionUSGovWest1,
		RegionUSGovEast1,
		RegionUSISOEast1,
		RegionUSISOBEast1,
		RegionUSISOWest1,
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
		Version1_15,
		Version1_16,
		Version1_17,
		Version1_18,
		Version1_19,
		Version1_20,
		Version1_21,
		Version1_22,
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
		Version1_23,
		Version1_24,
		Version1_25,
		Version1_26,
		Version1_27,
		Version1_28,
	}
}

// IsSupportedVersion returns true if the given Kubernetes version is supported by eksctl and EKS
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
		NodeVolumeTypeGP3,
		NodeVolumeTypeIO1,
		NodeVolumeTypeSC1,
		NodeVolumeTypeST1,
	}
}

// supportedAMIFamilies are the AMI families supported by EKS
func supportedAMIFamilies() []string {
	return []string{
		NodeImageFamilyAmazonLinux2,
		NodeImageFamilyUbuntu2004,
		NodeImageFamilyUbuntu1804,
		NodeImageFamilyBottlerocket,
		NodeImageFamilyWindowsServer2019CoreContainer,
		NodeImageFamilyWindowsServer2019FullContainer,
		NodeImageFamilyWindowsServer2022CoreContainer,
		NodeImageFamilyWindowsServer2022FullContainer,
	}
}

// validateSpotAllocationStrategy validates that the specified spot allocation strategy is supported.
func validateSpotAllocationStrategy(allocationStrategy string) error {
	var strategy ec2types.SpotAllocationStrategy
	for _, s := range strategy.Values() {
		if string(s) == allocationStrategy {
			return nil
		}
	}
	return fmt.Errorf("spotAllocationStrategy should be one of: %v", strategy.Values())
}

// EKSResourceAccountID provides worker node resources(ami/ecr image) in different aws account
// for different aws partitions & opt-in regions.
func EKSResourceAccountID(region string) string {
	switch region {
	case RegionAPEast1:
		return eksResourceAccountAPEast1
	case RegionCAWest1:
		return eksResourceAccountCAWest1
	case RegionMECentral1:
		return eksResourceAccountMECentral1
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
	case RegionEUSouth2:
		return eksResourceAccountEUSouth2
	case RegionEUCentral2:
		return eksResourceAccountEUCentral2
	case RegionAPSouth2:
		return eksResourceAccountAPSouth2
	case RegionAPSouthEast3:
		return eksResourceAccountAPSouthEast3
	case RegionAPSouthEast4:
		return eksResourceAccountAPSouthEast4
	case RegionILCentral1:
		return eksResourceAccountILCentral1
	case RegionUSISOEast1:
		return eksResourceAccountUSISOEast1
	case RegionUSISOBEast1:
		return eksResourceAccountUSISOBEast1
	case RegionUSISOWest1:
		return eksResourceAccountUSISOWest1
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
	// Internal fields
	// AccountID the ID of the account hosting this cluster
	AccountID string `json:"-"`
}

// KubernetesNetworkConfig contains cluster networking options
type KubernetesNetworkConfig struct {
	// Valid variants are `IPFamily` constants
	// +optional
	IPFamily string `json:"ipFamily,omitempty"`
	// ServiceIPv4CIDR is the CIDR range from where `ClusterIP`s are assigned
	ServiceIPv4CIDR string `json:"serviceIPv4CIDR,omitempty"`
}

func (k *KubernetesNetworkConfig) IPv6Enabled() bool {
	return strings.EqualFold(k.IPFamily, IPV6Family)
}

type EKSCTLCreated string

// ClusterStatus holds read-only attributes of a cluster
type ClusterStatus struct {
	Endpoint                 string                   `json:"endpoint,omitempty"`
	CertificateAuthorityData []byte                   `json:"certificateAuthorityData,omitempty"`
	ARN                      string                   `json:"arn,omitempty"`
	KubernetesNetworkConfig  *KubernetesNetworkConfig `json:"-"`
	ID                       string                   `json:"-"`
	APIServerUnreachable     bool                     `json:"-"`

	StackName     string        `json:"stackName,omitempty"`
	EKSCTLCreated EKSCTLCreated `json:"eksctlCreated,omitempty"`
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

func (c ClusterConfig) HasNodes() bool {
	for _, m := range c.ManagedNodeGroups {
		if m.GetDesiredCapacity() > 0 {
			return true
		}
	}

	for _, n := range c.NodeGroups {
		if n.GetDesiredCapacity() > 0 {
			return true
		}
	}
	return false
}

// ID returns the cluster ID.
func (c *ClusterConfig) ID() string {
	if c.Status != nil && c.Status.ID != "" {
		return c.Status.ID
	}
	return c.Metadata.Name
}

// Meta returns the cluster metadata.
func (c *ClusterConfig) Meta() *ClusterMeta {
	return c.Metadata
}

// GetStatus returns the cluster status.
func (c *ClusterConfig) GetStatus() *ClusterStatus {
	return c.Status
}

// IsFullyPrivate returns true if this is a fully-private cluster.
func (c *ClusterConfig) IsFullyPrivate() bool {
	return c.PrivateCluster != nil && c.PrivateCluster.Enabled
}

// IsControlPlaneOnOutposts returns true if the control plane is on Outposts.
func (c *ClusterConfig) IsControlPlaneOnOutposts() bool {
	return c.Outpost != nil && c.Outpost.ControlPlaneOutpostARN != ""
}

// GetOutpost returns the Outpost info.
func (c *ClusterConfig) GetOutpost() *Outpost {
	return c.Outpost
}

// FindNodeGroupOutpostARN finds nodegroups that are on Outposts and returns the Outpost ARN.
func (c *ClusterConfig) FindNodeGroupOutpostARN() (outpostARN string, found bool) {
	for _, ng := range c.NodeGroups {
		if ng.OutpostARN != "" {
			return ng.OutpostARN, true
		}
	}
	return "", false
}

// ClusterProvider is the interface to AWS APIs
type ClusterProvider interface {
	CloudFormation() awsapi.CloudFormation
	CloudFormationRoleARN() string
	CloudFormationDisableRollback() bool
	ASG() awsapi.ASG
	EKS() awsapi.EKS
	SSM() awsapi.SSM
	CloudTrail() awsapi.CloudTrail
	CloudWatchLogs() awsapi.CloudWatchLogs
	IAM() awsapi.IAM
	Region() string
	Profile() Profile
	WaitTimeout() time.Duration
	CredentialsProvider() aws.CredentialsProvider
	AWSConfig() aws.Config

	ELB() awsapi.ELB
	ELBV2() awsapi.ELBV2
	STS() awsapi.STS
	STSPresigner() STSPresigner
	EC2() awsapi.EC2
	Outposts() awsapi.Outposts
}

// STSPresigner defines the method to pre-sign GetCallerIdentity requests to add a proper header required by EKS for
// authentication from the outside.
//
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_sts_presigner.go . STSPresigner
type STSPresigner interface {
	PresignGetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// ProviderConfig holds global parameters for all interactions with AWS APIs
type ProviderConfig struct {
	CloudFormationRoleARN         string
	CloudFormationDisableRollback bool

	Region      string
	Profile     Profile
	WaitTimeout time.Duration
}

// Profile is the AWS profile to use.
type Profile struct {
	Name           string
	SourceIsEnvVar bool
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
	IAMIdentityMappings []*IAMIdentityMapping `json:"iamIdentityMappings,omitempty"`

	// +optional
	IdentityProviders []IdentityProvider `json:"identityProviders,omitempty"`

	// AccessConfig specifies the access config for a cluster.
	// +optional
	AccessConfig *AccessConfig `json:"accessConfig,omitempty"`

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

	// LocalZones specifies a list of local zones where the subnets should be created.
	// Only self-managed nodegroups can be launched in local zones. These subnets are not passed to EKS.
	// +optional
	LocalZones []string `json:"localZones,omitempty"`

	// See [CloudWatch support](/usage/cloudwatch-cluster-logging/)
	// +optional
	CloudWatch *ClusterCloudWatch `json:"cloudWatch,omitempty"`

	// +optional
	SecretsEncryption *SecretsEncryption `json:"secretsEncryption,omitempty"`

	Status *ClusterStatus `json:"-"`

	// future gitops plans, replacing the Git configuration above
	// +optional
	GitOps *GitOps `json:"gitops,omitempty"`

	// Karpenter specific configuration options.
	// +optional
	Karpenter *Karpenter `json:"karpenter,omitempty"`

	// Outpost specifies the Outpost configuration.
	// +optional
	Outpost *Outpost `json:"outpost,omitempty"`
}

// Outpost holds the Outpost configuration.
type Outpost struct {
	// ControlPlaneOutpostARN specifies the Outpost ARN in which the control plane should be created.
	ControlPlaneOutpostARN string `json:"controlPlaneOutpostARN"`
	// ControlPlaneInstanceType specifies the instance type to use for creating the control plane instances.
	ControlPlaneInstanceType string `json:"controlPlaneInstanceType"`
	// ControlPlanePlacement specifies the placement configuration for control plane instances on Outposts.
	ControlPlanePlacement *Placement `json:"controlPlanePlacement,omitempty"`
}

// GetInstanceType returns the control plane instance type.
func (o *Outpost) GetInstanceType() string {
	return o.ControlPlaneInstanceType
}

// SetInstanceType sets the control plane instance type.
func (o *Outpost) SetInstanceType(instanceType string) {
	o.ControlPlaneInstanceType = instanceType
}

// HasPlacementGroup reports whether this Outpost has a placement group.
func (o *Outpost) HasPlacementGroup() bool {
	return o.ControlPlanePlacement != nil
}

// OutpostInfo describes the Outpost info.
type OutpostInfo interface {
	// IsControlPlaneOnOutposts returns true if the control plane is on Outposts.
	IsControlPlaneOnOutposts() bool

	// GetOutpost returns the Outpost info.
	GetOutpost() *Outpost
}

// ErrUnsupportedLocalCluster is an error for when an unsupported operation is attempted on a local cluster.
var ErrUnsupportedLocalCluster = errors.New("this operation is not supported on Outposts clusters")

// Karpenter provides configuration options
type Karpenter struct {
	// Version defines the Karpenter version to install
	// +required
	Version string `json:"version"`
	// CreateServiceAccount create a service account or not.
	// +optional
	CreateServiceAccount *bool `json:"createServiceAccount,omitempty"`
	// DefaultInstanceProfile override the default IAM instance profile
	// +optional
	DefaultInstanceProfile *string `json:"defaultInstanceProfile,omitempty"`
	// WithSpotInterruptionQueue if true, adds all required policies and rules
	// for supporting Spot Interruption Queue on Karpenter deployments
	WithSpotInterruptionQueue *bool `json:"withSpotInterruptionQueue,omitempty"`
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
		KubernetesNetworkConfig: &KubernetesNetworkConfig{
			IPFamily: DefaultIPFamily,
		},
		VPC: NewClusterVPC(false),
		CloudWatch: &ClusterCloudWatch{
			ClusterLogging: &ClusterCloudWatchLogging{},
		},
		PrivateCluster: &PrivateCluster{},
		AccessConfig:   &AccessConfig{},
	}

	return cfg
}

// NewClusterVPC creates new VPC config for a cluster
func NewClusterVPC(ipv6Enabled bool) *ClusterVPC {
	cidr := DefaultCIDR()

	var nat *ClusterNAT
	if !ipv6Enabled {
		nat = DefaultClusterNAT()
	}

	return &ClusterVPC{
		Network: Network{
			CIDR: &cidr,
		},
		ManageSharedNodeSecurityGroupRules: Enabled(),
		NAT:                                nat,
		AutoAllocateIPv6:                   Disabled(),
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

func (c *ClusterConfig) IPv6Enabled() bool {
	return c.KubernetesNetworkConfig != nil && c.KubernetesNetworkConfig.IPv6Enabled()
}

// SetClusterState updates the cluster state and populates the ClusterStatus using *eks.Cluster.
func (c *ClusterConfig) SetClusterState(cluster *ekstypes.Cluster) error {
	if networkConfig := cluster.KubernetesNetworkConfig; networkConfig != nil && networkConfig.ServiceIpv4Cidr != nil {
		c.Status.KubernetesNetworkConfig = &KubernetesNetworkConfig{
			ServiceIPv4CIDR: *networkConfig.ServiceIpv4Cidr,
		}
		c.KubernetesNetworkConfig = &KubernetesNetworkConfig{
			ServiceIPv4CIDR: aws.ToString(cluster.KubernetesNetworkConfig.ServiceIpv4Cidr),
		}
	}
	data, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
	if err != nil {
		return fmt.Errorf("decoding certificate authority data: %w", err)
	}
	c.Status.Endpoint = *cluster.Endpoint
	c.Status.CertificateAuthorityData = data
	c.Status.ARN = *cluster.Arn
	if outpost := cluster.OutpostConfig; outpost != nil {
		if len(outpost.OutpostArns) != 1 {
			return fmt.Errorf("expected cluster to be associated with only one Outpost; got %v", outpost.OutpostArns)
		}
		outpostARN := outpost.OutpostArns[0]
		if c.IsControlPlaneOnOutposts() && c.Outpost.ControlPlaneOutpostARN != outpostARN {
			return fmt.Errorf("outpost.controlPlaneOutpostARN %q does not match the cluster's Outpost ARN %q", c.Outpost.ControlPlaneOutpostARN, outpostARN)
		}
		c.Outpost = &Outpost{
			ControlPlaneOutpostARN:   outpostARN,
			ControlPlaneInstanceType: *outpost.ControlPlaneInstanceType,
		}
	} else if c.IsControlPlaneOnOutposts() {
		return errors.New("outpost.controlPlaneOutpostARN is set but control plane is not on Outposts")
	}
	if cluster.Id != nil {
		c.Status.ID = *cluster.Id
	}
	return nil
}

// NewNodeGroup creates a new NodeGroup, and returns a pointer to it
func NewNodeGroup() *NodeGroup {
	return &NodeGroup{
		NodeGroupBase: &NodeGroupBase{
			PrivateNetworking: false,
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
					DeprecatedALBIngress:      Disabled(),
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
			DisableIMDSv1:    Enabled(),
			DisablePodIMDS:   Disabled(),
			InstanceSelector: &InstanceSelector{},
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
					DeprecatedALBIngress:      Disabled(),
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
// specific to an unmanaged nodegroup
type NodeGroup struct {
	*NodeGroupBase

	//+optional
	InstancesDistribution *NodeGroupInstancesDistribution `json:"instancesDistribution,omitempty"`

	// +optional
	ASGMetricsCollection []MetricsCollection `json:"asgMetricsCollection,omitempty"`

	// CPUCredits configures [T3 Unlimited](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/burstable-performance-instances-unlimited-mode.html), valid only for T-type instances
	// +optional
	CPUCredits *string `json:"cpuCredits,omitempty"`

	// Associate load balancers with auto scaling group
	// +optional
	ClassicLoadBalancerNames []string `json:"classicLoadBalancerNames,omitempty"`

	// Associate target group with auto scaling group
	// +optional
	TargetGroupARNs []string `json:"targetGroupARNs,omitempty"`

	// Taints taints to apply to the nodegroup
	// +optional
	Taints taintsWrapper `json:"taints,omitempty"`

	// UpdateConfig configures how to update NodeGroups.
	// +optional
	UpdateConfig *NodeGroupUpdateConfig `json:"updateConfig,omitempty"`

	// [Custom
	// address](/usage/vpc-networking/#custom-cluster-dns-address) used for DNS
	// lookups
	// +optional
	ClusterDNS string `json:"clusterDNS,omitempty"`

	// [Customize `kubelet` config](/usage/customizing-the-kubelet/)
	// +optional
	KubeletExtraConfig *InlineDocument `json:"kubeletExtraConfig,omitempty"`

	// ContainerRuntime defines the runtime (CRI) to use for containers on the node
	// +optional
	ContainerRuntime *string `json:"containerRuntime,omitempty"`

	// MaxInstanceLifetime defines the maximum amount of time in seconds an instance stays alive.
	// +optional
	MaxInstanceLifetime *int `json:"maxInstanceLifetime,omitempty"`

	// LocalZones specifies a list of local zones where the nodegroup should be launched.
	// The cluster should have been created with all of the local zones specified in this field.
	// +optional
	LocalZones []string `json:"localZones,omitempty"`
}

// GetContainerRuntime returns the container runtime.
func (n *NodeGroup) GetContainerRuntime() string {
	if n.ContainerRuntime != nil {
		return *n.ContainerRuntime
	}
	return ""
}

func (n *NodeGroup) InstanceTypeList() []string {
	if HasMixedInstances(n) {
		return n.InstancesDistribution.InstanceTypes
	}
	if n.InstanceType != "" {
		return []string{n.InstanceType}
	}
	return nil
}

// NGTaints implements NodePool
func (n *NodeGroup) NGTaints() []NodeGroupTaint {
	return n.Taints
}

// BaseNodeGroup implements NodePool
func (n *NodeGroup) BaseNodeGroup() *NodeGroupBase {
	return n.NodeGroupBase
}

func (n *NodeGroup) GetDesiredCapacity() int {
	if n.NodeGroupBase != nil {
		return n.NodeGroupBase.GetDesiredCapacity()
	}
	return 0
}

// GetInstanceType returns the instance type.
func (n *NodeGroup) GetInstanceType() string {
	return n.InstanceType
}

// SetInstanceType sets the instance type.
func (n *NodeGroup) SetInstanceType(instanceType string) {
	n.InstanceType = instanceType
}

// GitOps groups all configuration options related to enabling GitOps Toolkit on a
// cluster and linking it to a Git repository.
// Note: this will replace the older Git types
type GitOps struct {
	// Flux holds options to enable Flux v2 on your cluster
	Flux *Flux `json:"flux,omitempty"`
}

// Flux groups all configuration options related to a Git repository used for
// GitOps Toolkit (Flux v2).
type Flux struct {
	// The repository hosting service. Can be either Github or Gitlab.
	GitProvider string `json:"gitProvider,omitempty"`

	// Flags is an arbitrary map of string to string to pass any flags to Flux bootstrap
	// via eksctl see https://fluxcd.io/docs/ for information on all flags
	Flags FluxFlags `json:"flags,omitempty"`
}

// FluxFlags is a map of string for passing arbitrary flags to Flux bootstrap
type FluxFlags map[string]string

// HasGitOpsFluxConfigured returns true if gitops.flux configuration is not nil
func (c *ClusterConfig) HasGitOpsFluxConfigured() bool {
	return c.GitOps != nil && c.GitOps.Flux != nil
}

type (
	// NodeGroupSGs controls security groups for this nodegroup
	NodeGroupSGs struct {
		// AttachIDs attaches additional security groups to the nodegroup
		// +optional
		AttachIDs []string `json:"attachIDs,omitempty"`
		// WithShared attach the security group
		// shared among all nodegroups in the cluster
		// Not supported for managed nodegroups
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
		// AttachPolicy holds a policy document to attach
		// +optional
		AttachPolicy InlineDocument `json:"attachPolicy,omitempty"`
		// list of ARNs of the IAM policies to attach
		// +optional
		AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`
		// InstanceProfileARN holds the ARN of instance profile, not supported for Managed NodeGroups
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
		AWSLoadBalancerController *bool `json:"awsLoadBalancerController"`
		// +optional
		DeprecatedALBIngress *bool `json:"albIngress"`
		// +optional
		XRay *bool `json:"xRay"`
		// +optional
		CloudWatch *bool `json:"cloudWatch"`
	}

	// NodeGroupSSH holds all the ssh access configuration to a NodeGroup
	NodeGroupSSH struct {
		// +optional If Allow is true the SSH configuration provided is used, otherwise it is ignored. Only one of
		// PublicKeyPath, PublicKey and PublicKeyName can be configured
		Allow *bool `json:"allow"`
		// +optional The path to the SSH public key to be added to the nodes SSH keychain. If Allow is true this value
		// defaults to "~/.ssh/id_rsa.pub", otherwise the value is ignored.
		PublicKeyPath *string `json:"publicKeyPath,omitempty"`
		// +optional Public key to be added to the nodes SSH keychain. If Allow is false this value is ignored.
		PublicKey *string `json:"publicKey,omitempty"`
		// +optional Public key name in EC2 to be added to the nodes SSH keychain. If Allow is false this value
		// is ignored.
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
		// Enable [capacity
		// rebalancing](https://docs.aws.amazon.com/autoscaling/ec2/userguide/capacity-rebalance.html)
		// for spot instances
		// +optional
		CapacityRebalance bool `json:"capacityRebalance"`
	}

	// NodeGroupBottlerocket holds the configuration for Bottlerocket based
	// NodeGroups.
	NodeGroupBottlerocket struct {
		// +optional
		EnableAdminContainer *bool `json:"enableAdminContainer,omitempty"`
		// Settings contains any [bottlerocket
		// settings](https://bottlerocket.dev/en/os/latest/#/api/settings/)
		// +optional
		Settings *InlineDocument `json:"settings,omitempty"`
	}

	// NodeGroupUpdateConfig contains the configuration for updating NodeGroups.
	NodeGroupUpdateConfig struct {
		// MaxUnavailable sets the max number of nodes that can become unavailable
		// when updating a nodegroup (specified as number)
		// +optional
		MaxUnavailable *int `json:"maxUnavailable,omitempty"`

		// MaxUnavailablePercentage sets the max number of nodes that can become unavailable
		// when updating a nodegroup (specified as percentage)
		// +optional
		MaxUnavailablePercentage *int `json:"maxUnavailablePercentage,omitempty"`
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

	// NGTaints returns the taints to apply for this nodegroup
	NGTaints() []NodeGroupTaint

	// InstanceTypeList returns a list of instances that are configured for that nodegroup
	InstanceTypeList() []string
}

// VolumeMapping Additional Volume Configurations
type VolumeMapping struct {
	// +optional
	// VolumeSize gigabytes
	// Defaults to `80`
	VolumeSize *int `json:"volumeSize,omitempty"`
	// Valid variants are `VolumeType` constants
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
	VolumeThroughput *int `json:"volumeThroughput,omitempty"`
	// +optional
	SnapshotID *string `json:"snapshotID,omitempty"`
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
	// Applied to the Autoscaling Group and to the EC2 instances (unmanaged),
	// Applied to the EKS Nodegroup resource and to the EC2 instances (managed)
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
	// +optional
	VolumeEncrypted *bool `json:"volumeEncrypted,omitempty"`
	// +optional
	VolumeKmsKeyID *string `json:"volumeKmsKeyID,omitempty"`
	// +optional
	VolumeIOPS *int `json:"volumeIOPS,omitempty"`
	// +optional
	VolumeThroughput *int `json:"volumeThroughput,omitempty"`

	// Additional Volume Configurations
	// +optional
	AdditionalVolumes []*VolumeMapping `json:"additionalVolumes,omitempty"`

	// PreBootstrapCommands are executed before bootstrapping instances to the
	// cluster
	// +optional
	PreBootstrapCommands []string `json:"preBootstrapCommands,omitempty"`

	// Override `eksctl`'s bootstrapping script
	// +optional
	OverrideBootstrapCommand *string `json:"overrideBootstrapCommand,omitempty"`

	// Propagate all taints and labels to the ASG automatically.
	// +optional
	PropagateASGTags *bool `json:"propagateASGTags,omitempty"`

	// DisableIMDSv1 requires requests to the metadata service to use IMDSv2 tokens
	// Defaults to `true`
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

	// EFAEnabled creates the maximum allowed number of EFA-enabled network
	// cards on nodes in this group.
	// +optional
	EFAEnabled *bool `json:"efaEnabled,omitempty"`

	// InstanceSelector specifies options for EC2 instance selector
	InstanceSelector *InstanceSelector `json:"instanceSelector,omitempty"`

	// Internal fields
	// Some AMIs (bottlerocket) have a separate volume for the OS
	AdditionalEncryptedVolume string `json:"-"`

	// Bottlerocket specifies settings for Bottlerocket nodes
	// +optional
	Bottlerocket *NodeGroupBottlerocket `json:"bottlerocket,omitempty"`

	// Enable EC2 detailed monitoring
	// +optional
	EnableDetailedMonitoring *bool `json:"enableDetailedMonitoring,omitempty"`

	// CapacityReservation defines reservation policy for a nodegroup
	CapacityReservation *CapacityReservation `json:"capacityReservation,omitempty"`

	// OutpostARN specifies the Outpost ARN in which the nodegroup should be created.
	// +optional
	OutpostARN string `json:"outpostARN,omitempty"`
}

// CapacityReservation defines a nodegroup's Capacity Reservation targeting option
// +optional
type CapacityReservation struct {
	// CapacityReservationPreference defines a nodegroup's Capacity Reservation preferences (either 'open' or 'none')
	CapacityReservationPreference *string `json:"capacityReservationPreference,omitempty"`

	// CapacityReservationTarget defines a nodegroup's target Capacity Reservation or Capacity Reservation group (not both at the same time).
	CapacityReservationTarget *CapacityReservationTarget `json:"capacityReservationTarget,omitempty"`
}

type CapacityReservationTarget struct {
	CapacityReservationID               *string `json:"capacityReservationID,omitempty"`
	CapacityReservationResourceGroupARN *string `json:"capacityReservationResourceGroupARN,omitempty"`
}

// Placement specifies placement group information
type Placement struct {
	GroupName string `json:"groupName,omitempty"`
}

// ListOptions returns metav1.ListOptions with label selector for the nodegroup
func (n *NodeGroupBase) ListOptions() metav1.ListOptions {
	return metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", NodeGroupNameLabel, n.Name),
	}
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

// NodeGroupTaint represents a Kubernetes taint
type NodeGroupTaint struct {
	Key    string             `json:"key,omitempty"`
	Value  string             `json:"value,omitempty"`
	Effect corev1.TaintEffect `json:"effect,omitempty"`
}

// ManagedNodeGroup represents an EKS-managed nodegroup
type ManagedNodeGroup struct {
	*NodeGroupBase

	// InstanceTypes specifies a list of instance types
	InstanceTypes []string `json:"instanceTypes,omitempty"`

	// Spot creates a spot nodegroup
	Spot bool `json:"spot,omitempty"`

	// Taints taints to apply to the nodegroup
	Taints []NodeGroupTaint `json:"taints,omitempty"`

	// UpdateConfig configures how to update NodeGroups.
	// +optional
	UpdateConfig *NodeGroupUpdateConfig `json:"updateConfig,omitempty"`

	// LaunchTemplate specifies an existing launch template to use
	// for the nodegroup
	LaunchTemplate *LaunchTemplate `json:"launchTemplate,omitempty"`

	// ReleaseVersion the AMI version of the EKS optimized AMI to use
	ReleaseVersion string `json:"releaseVersion"`

	// Internal fields

	Unowned bool `json:"-"`
}

func (n *NodeGroupBase) GetDesiredCapacity() int {
	if n.ScalingConfig != nil && n.ScalingConfig.DesiredCapacity != nil {
		return *n.ScalingConfig.DesiredCapacity
	}
	return 0
}

func (m *ManagedNodeGroup) GetDesiredCapacity() int {
	if m.NodeGroupBase != nil {
		return m.NodeGroupBase.GetDesiredCapacity()
	}
	return 0
}

func (m *ManagedNodeGroup) InstanceTypeList() []string {
	if len(m.InstanceTypes) > 0 {
		return m.InstanceTypes
	}
	if m.InstanceType != "" {
		return []string{m.InstanceType}
	}
	return nil
}

func (m *ManagedNodeGroup) ListOptions() metav1.ListOptions {
	if m.Unowned {
		return metav1.ListOptions{
			LabelSelector: fmt.Sprintf("%s=%s", EKSNodeGroupNameLabel, m.NameString()),
		}
	}
	return m.NodeGroupBase.ListOptions()
}

// NGTaints implements NodePool
func (m *ManagedNodeGroup) NGTaints() []NodeGroupTaint {
	return m.Taints
}

// BaseNodeGroup implements NodePool
func (m *ManagedNodeGroup) BaseNodeGroup() *NodeGroupBase {
	return m.NodeGroupBase
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
	KeyARN string `json:"keyARN,omitempty"`
}

// PrivateCluster defines the configuration for a fully-private cluster.
type PrivateCluster struct {
	// Enabled enables creation of a fully-private cluster.
	Enabled bool `json:"enabled"`

	// SkipEndpointCreation skips the creation process for endpoints completely. This is only used in case of an already
	// provided VPC and if the user decided to set it to true.
	SkipEndpointCreation bool `json:"skipEndpointCreation"`

	// AdditionalEndpointServices specifies additional endpoint services that
	// must be enabled for private access.
	// Valid entries are "cloudformation", "autoscaling" and "logs".
	AdditionalEndpointServices []string `json:"additionalEndpointServices,omitempty"`
}

// InstanceSelector holds EC2 instance selector options
type InstanceSelector struct {
	// VCPUs specifies the number of vCPUs
	VCPUs int `json:"vCPUs,omitempty"`
	// Memory specifies the memory
	// The unit defaults to GiB
	Memory string `json:"memory,omitempty"`
	// GPUs specifies the number of GPUs.
	// It can be set to 0 to select non-GPU instance types.
	GPUs *int `json:"gpus,omitempty"`
	// CPU Architecture of the EC2 instance type.
	// Valid variants are:
	// `"x86_64"`
	// `"amd64"`
	// `"arm64"`
	CPUArchitecture string `json:"cpuArchitecture,omitempty"`
}

// IsZero returns true if all fields hold a zero value
func (is InstanceSelector) IsZero() bool {
	return is == InstanceSelector{}
}

// taintsWrapper handles unmarshalling both map[string]string and []NodeGroupTaint
type taintsWrapper []NodeGroupTaint

// UnmarshalJSON implements json.Unmarshaler
func (t *taintsWrapper) UnmarshalJSON(data []byte) error {
	taintsMap := map[string]string{}
	err := json.Unmarshal(data, &taintsMap)
	if err == nil {
		parsed := taints.Parse(taintsMap)
		for _, p := range parsed {
			*t = append(*t, NodeGroupTaint{
				Key:    p.Key,
				Value:  p.Value,
				Effect: p.Effect,
			})
		}
		return nil
	}

	var ngTaints []NodeGroupTaint
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&ngTaints); err != nil {
		return fmt.Errorf("taints must be a {string: string} or a [{key, value, effect}]: %w", err)
	}
	*t = ngTaints
	return nil
}

// AccessConfig specifies the access config for a cluster.
type AccessConfig struct {
	// AuthenticationMode specifies the authentication mode for a cluster.
	AuthenticationMode ekstypes.AuthenticationMode `json:"authenticationMode,omitempty"`

	// BootstrapClusterCreatorAdminPermissions specifies whether the cluster creator IAM principal was set as a cluster
	// admin access entry during cluster creation time.
	BootstrapClusterCreatorAdminPermissions *bool `json:"bootstrapClusterCreatorAdminPermissions,omitempty"`

	// AccessEntries specifies a list of access entries for the cluster.
	// +optional
	AccessEntries []AccessEntry `json:"accessEntries,omitempty"`
}

// UnsupportedFeatureError is an error that represents an unsupported feature
// +k8s:deepcopy-gen=false
type UnsupportedFeatureError struct {
	Message string
	Err     error
}

func (u *UnsupportedFeatureError) Error() string {
	return fmt.Sprintf("%s: %v", u.Message, u.Err)
}
