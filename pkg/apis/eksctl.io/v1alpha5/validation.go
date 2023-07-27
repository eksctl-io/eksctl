package v1alpha5

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/go-version"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"

	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/taints"

	"k8s.io/apimachinery/pkg/util/validation"
	kubeletapis "k8s.io/kubelet/pkg/apis"
)

// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-blockdevicemapping-ebs.html
const (
	MinThroughput = DefaultNodeVolumeThroughput
	MaxThroughput = 1000
	MinIO1Iops    = DefaultNodeVolumeIO1IOPS
	MaxIO1Iops    = 64000
	MinGP3Iops    = DefaultNodeVolumeGP3IOPS
	MaxGP3Iops    = 16000
	OneDay        = 86400
)

var (
	// ErrClusterEndpointNoAccess indicates the config prevents API access
	ErrClusterEndpointNoAccess = errors.New("Kubernetes API access must have one of public or private clusterEndpoints enabled")

	// ErrClusterEndpointPrivateOnly warns private-only access requires changes
	// to AWS resource configuration in order to effectively use clients in the VPC
	ErrClusterEndpointPrivateOnly = errors.New("warning, having public access disallowed will subsequently interfere with some " +
		"features of eksctl. This will require running subsequent eksctl (and Kubernetes) " +
		"commands/API calls from within the VPC.  Running these in the VPC requires making " +
		"updates to some AWS resources.  See: " +
		"https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html " +
		"for more details")
)

// NOTE: we don't use k8s.io/apimachinery/pkg/util/sets here to keep API package free of dependencies
type nameSet map[string]struct{}

func (s nameSet) checkUnique(path, name string) (bool, error) {
	if _, notUnique := s[name]; notUnique {
		return false, fmt.Errorf("%s %q is not unique", path, name)
	}
	s[name] = struct{}{}
	return true, nil
}

func setNonEmpty(field string) error {
	return fmt.Errorf("%s must be set and non-empty", field)
}

// ValidateClusterConfig checks compatible fields of a given ClusterConfig
func ValidateClusterConfig(cfg *ClusterConfig) error {
	if IsDisabled(cfg.IAM.WithOIDC) && len(cfg.IAM.ServiceAccounts) > 0 {
		return fmt.Errorf("iam.withOIDC must be enabled explicitly for iam.serviceAccounts to be created")
	}

	saNames := nameSet{}
	for i, sa := range cfg.IAM.ServiceAccounts {
		path := fmt.Sprintf("iam.serviceAccounts[%d]", i)
		if sa.Name == "" {
			return fmt.Errorf("%s.name must be set", path)
		}
		if ok, err := saNames.checkUnique("<namespace>/<name> of "+path, sa.NameString()); !ok {
			return err
		}
		if !sa.WellKnownPolicies.HasPolicy() && len(sa.AttachPolicyARNs) == 0 && sa.AttachPolicy == nil && sa.AttachRoleARN == "" {
			return fmt.Errorf("%[1]s.wellKnownPolicies, %[1]s.attachPolicyARNs,%[1]s.attachRoleARN  or %[1]s.attachPolicy must be set", path)
		}
	}

	if err := cfg.validateKubernetesNetworkConfig(); err != nil {
		return err
	}

	// names must be unique across both managed and unmanaged nodegroups
	ngNames := nameSet{}
	validateNg := func(ng *NodeGroupBase, path string) error {
		if ng.Name == "" {
			return fmt.Errorf("%s.name must be set", path)
		}
		if _, err := ngNames.checkUnique(path+".name", ng.Name); err != nil {
			return err
		}
		if cfg.PrivateCluster.Enabled && !ng.PrivateNetworking {
			return fmt.Errorf("%s.privateNetworking must be enabled for a fully-private cluster", path)
		}
		return nil
	}

	if err := validateIdentityProviders(cfg.IdentityProviders); err != nil {
		return err
	}

	var ngOutpostARN string
	for i, ng := range cfg.NodeGroups {
		path := fmt.Sprintf("nodeGroups[%d]", i)
		if err := validateNg(ng.NodeGroupBase, path); err != nil {
			return err
		}
		if ng.OutpostARN != "" {
			if ngOutpostARN != "" && ng.OutpostARN != ngOutpostARN {
				return fmt.Errorf("cannot create nodegroups in two different Outposts; got Outpost ARN %q and %q", ngOutpostARN, ng.OutpostARN)
			}
			ngOutpostARN = ng.OutpostARN
		}
	}

	for i, ng := range cfg.ManagedNodeGroups {
		path := fmt.Sprintf("managedNodeGroups[%d]", i)
		if err := validateNg(ng.NodeGroupBase, path); err != nil {
			return err
		}
	}

	if err := validateCloudWatchLogging(cfg); err != nil {
		return err
	}

	if err := validateAvailabilityZones(cfg.AvailabilityZones); err != nil {
		return err
	}

	if cfg.Outpost != nil {
		if cfg.Outpost.ControlPlaneOutpostARN == "" {
			return errors.New("outpost.controlPlaneOutpostARN is required for Outposts")
		}
		if err := validateOutpostARN(cfg.Outpost.ControlPlaneOutpostARN); err != nil {
			return err
		}

		if cfg.IPv6Enabled() {
			return errors.New("IPv6 is not supported on Outposts")
		}
		if len(cfg.Addons) > 0 {
			return errors.New("Addons are not supported on Outposts")
		}
		if len(cfg.IdentityProviders) > 0 {
			return errors.New("Identity Providers are not supported on Outposts")
		}
		if len(cfg.FargateProfiles) > 0 {
			return errors.New("Fargate is not supported on Outposts")
		}
		if cfg.Karpenter != nil {
			return errors.New("Karpenter is not supported on Outposts")
		}
		if cfg.SecretsEncryption != nil && cfg.SecretsEncryption.KeyARN != "" {
			return errors.New("KMS encryption is not supported on Outposts")
		}
		const zonesErr = "cannot specify %s on Outposts; the AZ defaults to the Outpost AZ"
		if len(cfg.AvailabilityZones) > 0 {
			return fmt.Errorf(zonesErr, "availabilityZones")
		}
		if len(cfg.LocalZones) > 0 {
			return fmt.Errorf(zonesErr, "localZones")
		}
		if cfg.GitOps != nil {
			return errors.New("GitOps is not supported on Outposts")
		}
		if cfg.IAM != nil && IsEnabled(cfg.IAM.WithOIDC) {
			return errors.New("iam.withOIDC is not supported on Outposts")
		}
		if cfg.VPC != nil {
			if IsEnabled(cfg.VPC.AutoAllocateIPv6) {
				return errors.New("autoAllocateIPv6 is not supported on Outposts")
			}
			if len(cfg.VPC.PublicAccessCIDRs) > 0 {
				return errors.New("publicAccessCIDRs is not supported on Outposts")
			}
		}
	} else if ngOutpostARN != "" && cfg.IsFullyPrivate() {
		return errors.New("nodeGroup.outpostARN is not supported on a fully-private cluster (privateCluster.enabled)")
	}

	if err := cfg.ValidateVPCConfig(); err != nil {
		return err
	}

	if err := ValidateSecretsEncryption(cfg); err != nil {
		return err
	}

	if err := validateIAMIdentityMappings(cfg); err != nil {
		return err
	}

	if err := validateKarpenterConfig(cfg); err != nil {
		return fmt.Errorf("failed to validate Karpenter config: %w", err)
	}

	return nil
}

// ValidateClusterVersion validates the cluster version.
func ValidateClusterVersion(clusterConfig *ClusterConfig) error {
	if clusterVersion := clusterConfig.Metadata.Version; clusterVersion != "" && clusterVersion != DefaultVersion && !IsSupportedVersion(clusterVersion) {
		if IsDeprecatedVersion(clusterVersion) {
			return fmt.Errorf("invalid version, %s is no longer supported, supported values: %s\nsee also: https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html", clusterVersion, strings.Join(SupportedVersions(), ", "))
		}
		return fmt.Errorf("invalid version, supported values: %s", strings.Join(SupportedVersions(), ", "))
	}
	return nil
}

func validateKarpenterConfig(cfg *ClusterConfig) error {
	if cfg.Karpenter == nil {
		return nil
	}
	if cfg.Karpenter.Version == "" {
		return errors.New("version field is required if installing Karpenter is enabled")
	}

	v, err := version.NewVersion(cfg.Karpenter.Version)
	if err != nil {
		return fmt.Errorf("failed to parse Karpenter version %q: %w", cfg.Karpenter.Version, err)
	}

	supportedVersion, err := version.NewVersion(supportedKarpenterVersion)
	if err != nil {
		return fmt.Errorf("failed to parse supported Karpenter version %s: %w", supportedKarpenterVersion, err)
	}

	if v.LessThan(supportedVersion) {
		return fmt.Errorf("minimum supported version is %s", supportedKarpenterVersion)
	}

	if IsDisabled(cfg.IAM.WithOIDC) {
		return errors.New("iam.withOIDC must be enabled with Karpenter")
	}
	return nil
}

func validateCloudWatchLogging(clusterConfig *ClusterConfig) error {
	if !clusterConfig.HasClusterCloudWatchLogging() {
		if clusterConfig.CloudWatch != nil &&
			clusterConfig.CloudWatch.ClusterLogging != nil &&
			clusterConfig.CloudWatch.ClusterLogging.LogRetentionInDays != 0 {
			return errors.New("cannot set cloudWatch.clusterLogging.logRetentionInDays without enabling log types")
		}
		return nil
	}

	for i, logType := range clusterConfig.CloudWatch.ClusterLogging.EnableTypes {
		isUnknown := true
		for _, knownLogType := range SupportedCloudWatchClusterLogTypes() {
			if logType == knownLogType {
				isUnknown = false
			}
		}
		if isUnknown {
			return errors.Errorf("log type %q (cloudWatch.clusterLogging.enableTypes[%d]) is unknown", logType, i)
		}
	}
	if logRetentionDays := clusterConfig.CloudWatch.ClusterLogging.LogRetentionInDays; logRetentionDays != 0 {
		for _, v := range LogRetentionInDaysValues {
			if v == logRetentionDays {
				return nil
			}
		}
		return errors.Errorf("invalid value %d for logRetentionInDays; supported values are %v", logRetentionDays, LogRetentionInDaysValues)
	}

	return nil
}

// ValidateVPCConfig validates the vpc setting if it is defined.
func (c *ClusterConfig) ValidateVPCConfig() error {
	if c.VPC == nil {
		return nil
	}

	if err := c.ValidatePrivateCluster(); err != nil {
		return err
	}

	if err := c.ValidateClusterEndpointConfig(); err != nil {
		return err
	}

	if len(c.VPC.ExtraCIDRs) > 0 {
		cidrs, err := validateCIDRs(c.VPC.ExtraCIDRs)
		if err != nil {
			return err
		}
		c.VPC.ExtraCIDRs = cidrs
	}
	if len(c.VPC.PublicAccessCIDRs) > 0 {
		cidrs, err := validateCIDRs(c.VPC.PublicAccessCIDRs)
		if err != nil {
			return err
		}
		c.VPC.PublicAccessCIDRs = cidrs
	}
	if len(c.VPC.ExtraIPv6CIDRs) > 0 {
		if !c.IPv6Enabled() {
			return fmt.Errorf("cannot specify vpc.extraIPv6CIDRs with an IPv4 cluster")
		}
		cidrs, err := validateCIDRs(c.VPC.ExtraIPv6CIDRs)
		if err != nil {
			return err
		}
		c.VPC.ExtraIPv6CIDRs = cidrs
	}

	if (c.VPC.IPv6Cidr != "" || c.VPC.IPv6Pool != "") && !c.IPv6Enabled() {
		return fmt.Errorf("Ipv6Cidr and Ipv6CidrPool are only supported when IPFamily is set to IPv6")
	}

	if c.IPv6Enabled() {
		if IsEnabled(c.VPC.AutoAllocateIPv6) {
			return fmt.Errorf("auto allocate ipv6 is not supported with IPv6")
		}
		if err := c.ipv6CidrsValid(); err != nil {
			return err
		}
		if c.VPC.NAT != nil {
			return fmt.Errorf("setting NAT is not supported with IPv6")
		}
		if len(c.LocalZones) > 0 {
			return errors.New("localZones are not supported with IPv6")
		}
	}

	// manageSharedNodeSecurityGroupRules cannot be disabled if using eksctl managed security groups
	if c.VPC.SharedNodeSecurityGroup == "" && IsDisabled(c.VPC.ManageSharedNodeSecurityGroupRules) {
		return errors.New("vpc.manageSharedNodeSecurityGroupRules must be enabled when using eksctl-managed security groups")
	}

	if c.VPC.HostnameType != "" {
		if c.HasAnySubnets() {
			return errors.New("vpc.hostnameType is not supported with a pre-existing VPC")
		}
		var hostnameType ec2types.HostnameType
		found := false
		for _, h := range hostnameType.Values() {
			if h == ec2types.HostnameType(c.VPC.HostnameType) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid value %q for vpc.hostnameType; supported values are %v", c.VPC.HostnameType, hostnameType.Values())
		}
	}

	if len(c.LocalZones) > 0 {
		if c.VPC.ID != "" {
			return errors.New("localZones are not supported with a pre-existing VPC")
		}
		if c.VPC.NAT != nil && c.VPC.NAT.Gateway != nil && *c.VPC.NAT.Gateway == ClusterHighlyAvailableNAT {
			return fmt.Errorf("%s NAT gateway is not supported for localZones", ClusterHighlyAvailableNAT)
		}
	}

	return nil
}

func (c *ClusterConfig) unsupportedVPCCNIAddonVersion() (bool, error) {
	for _, addon := range c.Addons {
		if addon.Name == VPCCNIAddon {
			if addon.Version == "" {
				return false, nil
			}
			if addon.Version == "latest" {
				return false, nil
			}

			return versionLessThan(addon.Version, minimumVPCCNIVersionForIPv6)
		}
	}
	return false, nil
}

func versionLessThan(v1, v2 string) (bool, error) {
	v1Version, err := parseVersion(v1)
	if err != nil {
		return false, err
	}
	v2Version, err := parseVersion(v2)
	if err != nil {
		return false, err
	}
	return v1Version.LessThan(v2Version), nil
}

func parseVersion(v string) (*version.Version, error) {
	version, err := version.NewVersion(v)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %q: %w", v, err)
	}
	return version, nil
}

func (c *ClusterConfig) ipv6CidrsValid() error {
	if c.VPC.IPv6Cidr == "" && c.VPC.IPv6Pool == "" {
		return nil
	}

	if c.VPC.IPv6Cidr != "" && c.VPC.IPv6Pool != "" {
		if c.VPC.ID != "" {
			return fmt.Errorf("cannot provide VPC.IPv6Cidr when using a pre-existing VPC.ID")
		}
		return nil
	}
	return fmt.Errorf("Ipv6Cidr and Ipv6Pool must both be configured to use a custom IPv6 CIDR and address pool")
}

// addonContainsManagedAddons finds managed addons in the config and returns those it couldn't find.
func (c *ClusterConfig) addonContainsManagedAddons(addons []string) []string {
	var missing []string
	for _, a := range addons {
		found := false
		for _, add := range c.Addons {
			if strings.ToLower(add.Name) == a {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, a)
		}
	}
	return missing
}

// ValidateClusterEndpointConfig checks the endpoint configuration for potential issues
func (c *ClusterConfig) ValidateClusterEndpointConfig() error {
	if c.VPC.ClusterEndpoints != nil {
		if !c.HasClusterEndpointAccess() {
			return ErrClusterEndpointNoAccess
		}
		endpts := c.VPC.ClusterEndpoints

		if noAccess(endpts) {
			return ErrClusterEndpointNoAccess
		}
	}
	return nil
}

// ValidatePrivateCluster validates the private cluster config
func (c *ClusterConfig) ValidatePrivateCluster() error {
	if c.PrivateCluster.Enabled {
		if c.VPC != nil && c.VPC.ID != "" && len(c.VPC.Subnets.Private) == 0 {
			return errors.New("vpc.subnets.private must be specified in a fully-private cluster when a pre-existing VPC is supplied")
		}

		if additionalEndpoints := c.PrivateCluster.AdditionalEndpointServices; len(additionalEndpoints) > 0 {
			if c.PrivateCluster.SkipEndpointCreation {
				return errors.New("privateCluster.additionalEndpointServices cannot be set when privateCluster.skipEndpointCreation is true")
			}
			if err := ValidateAdditionalEndpointServices(additionalEndpoints); err != nil {
				return fmt.Errorf("invalid value in privateCluster.additionalEndpointServices: %w", err)
			}
		}

		if c.VPC != nil && c.VPC.ClusterEndpoints == nil {
			c.VPC.ClusterEndpoints = &ClusterEndpoints{}
		}
		if len(c.LocalZones) > 0 {
			return errors.New("localZones cannot be used in a fully-private cluster")
		}
		// public access is initially enabled to allow running operations that access the Kubernetes API
		if !c.IsControlPlaneOnOutposts() {
			c.VPC.ClusterEndpoints.PublicAccess = Enabled()
			c.VPC.ClusterEndpoints.PrivateAccess = Enabled()
		}
	}
	return nil
}

// validateKubernetesNetworkConfig validates the k8s network config
func (c *ClusterConfig) validateKubernetesNetworkConfig() error {
	if c.KubernetesNetworkConfig == nil {
		return nil
	}
	if c.KubernetesNetworkConfig.ServiceIPv4CIDR != "" {
		if c.IPv6Enabled() {
			return errors.New("service IPv4 CIDR is not supported with IPv6")
		}
		serviceIP := c.KubernetesNetworkConfig.ServiceIPv4CIDR
		if _, _, err := net.ParseCIDR(serviceIP); serviceIP != "" && err != nil {
			return errors.Wrap(err, "invalid IPv4 CIDR for kubernetesNetworkConfig.serviceIPv4CIDR")
		}
	}

	switch strings.ToLower(c.KubernetesNetworkConfig.IPFamily) {
	case strings.ToLower(IPV4Family), "":
	case strings.ToLower(IPV6Family):
		if missing := c.addonContainsManagedAddons([]string{VPCCNIAddon, CoreDNSAddon, KubeProxyAddon}); len(missing) != 0 {
			return fmt.Errorf("the default core addons must be defined for IPv6; missing addon(s): %s", strings.Join(missing, ", "))
		}

		unsupportedVersion, err := c.unsupportedVPCCNIAddonVersion()
		if err != nil {
			return err
		}

		if unsupportedVersion {
			return fmt.Errorf("%s version must be at least version %s for IPv6", VPCCNIAddon, minimumVPCCNIVersionForIPv6)
		}

		if c.IAM == nil || c.IAM != nil && IsDisabled(c.IAM.WithOIDC) {
			return fmt.Errorf("oidc needs to be enabled if IPv6 is set")
		}

		if version, err := utils.CompareVersions(c.Metadata.Version, Version1_21); err != nil {
			return fmt.Errorf("failed to convert %s cluster version to semver: %w", c.Metadata.Version, err)
		} else if err == nil && version == -1 {
			return fmt.Errorf("cluster version must be >= %s", Version1_21)
		}
	default:
		return fmt.Errorf("invalid value %q for ipFamily; allowed are %s and %s", c.KubernetesNetworkConfig.IPFamily, IPV4Family, IPV6Family)
	}

	return nil
}

// NoAccess returns true if neither public are private cluster endpoint access is enabled and false otherwise
func noAccess(ces *ClusterEndpoints) bool {
	return !(IsEnabled(ces.PublicAccess) || IsEnabled(ces.PrivateAccess))
}

// PrivateOnly returns true if public cluster endpoint access is disabled and private cluster endpoint access is enabled, and false otherwise
func PrivateOnly(ces *ClusterEndpoints) bool {
	return !*ces.PublicAccess && *ces.PrivateAccess
}

func validateNodeGroupBase(np NodePool, path string, controlPlaneOnOutposts bool) error {
	ng := np.BaseNodeGroup()
	if ng.VolumeSize == nil {
		errCantSet := func(field string) error {
			return fmt.Errorf("%s.%s cannot be set without %s.volumeSize", path, field, path)
		}
		if IsSetAndNonEmptyString(ng.VolumeName) {
			return errCantSet("volumeName")
		}
		if IsEnabled(ng.VolumeEncrypted) {
			return errCantSet("volumeEncrypted")
		}
		if IsSetAndNonEmptyString(ng.VolumeKmsKeyID) {
			return errCantSet("volumeKmsKeyID")
		}
	}

	if err := validateVolumeOpts(ng, path, controlPlaneOnOutposts); err != nil {
		return err
	}

	if ng.VolumeEncrypted == nil || IsDisabled(ng.VolumeEncrypted) {
		if IsSetAndNonEmptyString(ng.VolumeKmsKeyID) {
			return fmt.Errorf("%s.volumeKmsKeyID can not be set without %s.volumeEncrypted enabled explicitly", path, path)
		}
	}
	if ng.MaxPodsPerNode < 0 {
		return fmt.Errorf("%s.maxPodsPerNode cannot be negative", path)
	}

	if IsEnabled(ng.DisablePodIMDS) && ng.IAM != nil {
		fmtFieldConflictErr := func(_ string) error {
			return fmt.Errorf("%s.disablePodIMDS and %s.iam.withAddonPolicies cannot be set at the same time", path, path)
		}
		if err := validateNodeGroupIAMWithAddonPolicies(ng.IAM.WithAddonPolicies, fmtFieldConflictErr); err != nil {
			return err
		}
	}

	if len(ng.AvailabilityZones) > 0 && len(ng.Subnets) > 0 {
		return fmt.Errorf("only one of %[1]s.subnets or %[1]s.availabilityZones should be set", path)
	}

	if ng.Placement != nil {
		if ng.Placement.GroupName == "" {
			return fmt.Errorf("%s.placement.groupName must be set and non-empty", path)
		}
	}

	if IsEnabled(ng.EFAEnabled) {
		if len(ng.AvailabilityZones) > 1 || len(ng.Subnets) > 1 {
			return fmt.Errorf("%s.efaEnabled nodegroups must have only one subnet or one availability zone", path)
		}
	}

	if ng.AMIFamily != "" {
		if !isSupportedAMIFamily(ng.AMIFamily) {
			if ng.AMIFamily == NodeImageFamilyWindowsServer20H2CoreContainer || ng.AMIFamily == NodeImageFamilyWindowsServer2004CoreContainer {
				return fmt.Errorf("AMI Family %s is deprecated. For more information, head to the Amazon documentation on Windows AMIs (https://docs.aws.amazon.com/eks/latest/userguide/eks-optimized-windows-ami.html)", ng.AMIFamily)
			}
			return fmt.Errorf("AMI Family %s is not supported - use one of: %s", ng.AMIFamily, strings.Join(supportedAMIFamilies(), ", "))
		}
		if controlPlaneOnOutposts && ng.AMIFamily != NodeImageFamilyAmazonLinux2 {
			return fmt.Errorf("only %s is supported on local clusters", NodeImageFamilyAmazonLinux2)
		}
	}

	if ng.SSH != nil {
		if enableSSM := ng.SSH.EnableSSM; enableSSM != nil {
			if !*enableSSM {
				return errors.New("SSM agent is now built into EKS AMIs and cannot be disabled")
			}
			logger.Warning("SSM is now enabled by default; `ssh.enableSSM` is deprecated and will be removed in a future release")
		}
	}

	if instanceutils.IsNvidiaInstanceType(SelectInstanceType(np)) &&
		(ng.AMIFamily != NodeImageFamilyAmazonLinux2 && ng.AMIFamily != NodeImageFamilyBottlerocket && ng.AMIFamily != "") {
		logger.Warning("%s does not ship with NVIDIA GPU drivers installed, hence won't support running GPU-accelerated workloads out of the box", ng.AMIFamily)
	}

	if ng.AMIFamily != NodeImageFamilyAmazonLinux2 && ng.AMIFamily != "" {
		instanceType := SelectInstanceType(np)
		unsupportedErr := func(instanceTypeName string) error {
			return fmt.Errorf("%s instance types are not supported for %s", instanceTypeName, ng.AMIFamily)
		}
		// Only AL2 supports Inferentia hosts.
		if instanceutils.IsInferentiaInstanceType(instanceType) {
			return unsupportedErr("Inferentia")
		}
		// Only AL2 supports Trainium hosts.
		if instanceutils.IsTrainiumInstanceType(instanceType) {
			return unsupportedErr("Trainium")
		}
	}

	if ng.CapacityReservation != nil {
		if ng.CapacityReservation.CapacityReservationPreference != nil {
			if ng.CapacityReservation.CapacityReservationTarget != nil {
				return errors.New("only one of CapacityReservationPreference or CapacityReservationTarget may be specified at a time")
			}

			if *ng.CapacityReservation.CapacityReservationPreference != OpenCapacityReservation && *ng.CapacityReservation.CapacityReservationPreference != NoneCapacityReservation {
				return fmt.Errorf(`accepted values include "open" and "none"; got "%s"`, *ng.CapacityReservation.CapacityReservationPreference)
			}
		}

		if ng.CapacityReservation.CapacityReservationTarget != nil {
			if ng.CapacityReservation.CapacityReservationTarget.CapacityReservationID != nil && ng.CapacityReservation.CapacityReservationTarget.CapacityReservationResourceGroupARN != nil {
				return errors.New("only one of CapacityReservationID or CapacityReservationResourceGroupARN may be specified at a time")
			}
		}
	}

	return nil
}

func validateVolumeOpts(ng *NodeGroupBase, path string, controlPlaneOnOutposts bool) error {
	if ng.VolumeType != nil {
		volumeType := *ng.VolumeType
		if ng.VolumeIOPS != nil && !(volumeType == NodeVolumeTypeIO1 || volumeType == NodeVolumeTypeGP3) {
			return fmt.Errorf("%s.volumeIOPS is only supported for %s and %s volume types", path, NodeVolumeTypeIO1, NodeVolumeTypeGP3)
		}

		if volumeType == NodeVolumeTypeIO1 {
			if ng.VolumeIOPS != nil && !(*ng.VolumeIOPS >= MinIO1Iops && *ng.VolumeIOPS <= MaxIO1Iops) {
				return fmt.Errorf("value for %s.volumeIOPS must be within range %d-%d", path, MinIO1Iops, MaxIO1Iops)
			}
		}

		if ng.VolumeThroughput != nil && volumeType != NodeVolumeTypeGP3 {
			return fmt.Errorf("%s.volumeThroughput is only supported for %s volume type", path, NodeVolumeTypeGP3)
		}

		if controlPlaneOnOutposts && volumeType != NodeVolumeTypeGP2 {
			return fmt.Errorf("cannot set %q for %s.volumeType; only %q volume types are supported on Outposts", volumeType, path, NodeVolumeTypeGP2)
		}
	}

	if ng.VolumeType == nil || *ng.VolumeType == NodeVolumeTypeGP3 {
		if ng.VolumeIOPS != nil && !(*ng.VolumeIOPS >= MinGP3Iops && *ng.VolumeIOPS <= MaxGP3Iops) {
			return fmt.Errorf("value for %s.volumeIOPS must be within range %d-%d", path, MinGP3Iops, MaxGP3Iops)
		}

		if ng.VolumeThroughput != nil && !(*ng.VolumeThroughput >= MinThroughput && *ng.VolumeThroughput <= MaxThroughput) {
			return fmt.Errorf("value for %s.volumeThroughput must be within range %d-%d", path, MinThroughput, MaxThroughput)
		}
	}

	return nil
}

func validateIdentityProvider(idP IdentityProvider) error {
	switch idP := (idP.Inner).(type) {
	case *OIDCIdentityProvider:
		if idP.Name == "" {
			return setNonEmpty("name")
		}
		if idP.ClientID == "" {
			return setNonEmpty("clientID")
		}
		if idP.IssuerURL == "" {
			return setNonEmpty("issuerURL")
		}
	}
	return nil
}

func validateIdentityProviders(idPs []IdentityProvider) error {
	for k, idP := range idPs {
		if err := validateIdentityProvider(idP); err != nil {
			return errors.Wrapf(err, "identityProviders[%d] is invalid", k)
		}
	}
	return nil
}

type unsupportedFieldError struct {
	ng    *NodeGroupBase
	path  string
	field string
}

func (ue *unsupportedFieldError) Error() string {
	return fmt.Sprintf("%s is not supported for %s nodegroups (path=%s.%s)", ue.field, ue.ng.AMIFamily, ue.path, ue.field)
}

// IsInvalidNameArg checks whether the name contains invalid characters
func IsInvalidNameArg(name string) bool {
	re := regexp.MustCompile(`[^a-zA-Z0-9\-]+`)
	return re.MatchString(name)
}

// errInvalidName error when invalid characters for a name is provided
func ErrInvalidName(name string) error {
	return fmt.Errorf("validation for %s failed, name must satisfy regular expression pattern: [a-zA-Z][-a-zA-Z0-9]*", name)
}

func validateNodeGroupName(name string) error {
	if name != "" && IsInvalidNameArg(name) {
		return ErrInvalidName(name)
	}

	return nil
}

// ValidateNodeGroup checks compatible fields of a given nodegroup
func ValidateNodeGroup(i int, ng *NodeGroup, cfg *ClusterConfig) error {
	normalizeAMIFamily(ng.BaseNodeGroup())
	path := fmt.Sprintf("nodeGroups[%d]", i)
	if err := validateNodeGroupBase(ng, path, cfg.IsControlPlaneOnOutposts()); err != nil {
		return err
	}

	if err := validateNodeGroupName(ng.Name); err != nil {
		return err
	}

	if ng.IAM != nil {
		if err := validateNodeGroupIAM(ng.IAM, ng.IAM.InstanceProfileARN, "instanceProfileARN", path); err != nil {
			return err
		}
		if err := validateNodeGroupIAM(ng.IAM, ng.IAM.InstanceRoleARN, "instanceRoleARN", path); err != nil {
			return err
		}
		if attachPolicyARNs := ng.IAM.AttachPolicyARNs; len(attachPolicyARNs) > 0 {
			for _, policyARN := range attachPolicyARNs {
				if _, err := arn.Parse(policyARN); err != nil {
					return errors.Wrapf(err, "invalid ARN %q in %s.iam.attachPolicyARNs", policyARN, path)
				}

			}
		}
		if err := validateDeprecatedIAMFields(ng.IAM); err != nil {
			return err
		}
	}

	if ng.AMI != "" && ng.AMIFamily == "" {
		return errors.Errorf("when using a custom AMI, amiFamily needs to be explicitly set via config file or via --node-ami-family flag")
	}

	if ng.Bottlerocket != nil && ng.AMIFamily != NodeImageFamilyBottlerocket {
		return fmt.Errorf(`bottlerocket config can only be used with amiFamily "Bottlerocket" but found "%s" (path=%s.bottlerocket)`,
			ng.AMIFamily, path)
	}

	if ng.AMI != "" && ng.OverrideBootstrapCommand == nil && ng.AMIFamily != NodeImageFamilyBottlerocket && !IsWindowsImage(ng.AMIFamily) {
		return errors.Errorf("%[1]s.overrideBootstrapCommand is required when using a custom AMI (%[1]s.ami)", path)
	}

	if err := validateTaints(ng.Taints); err != nil {
		return err
	}

	if err := validateNodeGroupLabels(ng.Labels); err != nil {
		return err
	}

	if ng.SSH != nil {
		if err := validateNodeGroupSSH(ng.SSH); err != nil {
			return err
		}
	}

	fieldNotSupported := func(field string) error {
		return &unsupportedFieldError{
			ng:    ng.NodeGroupBase,
			path:  path,
			field: field,
		}
	}

	if IsWindowsImage(ng.AMIFamily) {
		if ng.KubeletExtraConfig != nil {
			return fieldNotSupported("kubeletExtraConfig")
		}
	} else if ng.AMIFamily == NodeImageFamilyBottlerocket {
		if ng.KubeletExtraConfig != nil {
			return fieldNotSupported("kubeletExtraConfig")
		}
		if ng.PreBootstrapCommands != nil {
			return fieldNotSupported("preBootstrapCommands")
		}
		if ng.OverrideBootstrapCommand != nil {
			return fieldNotSupported("overrideBootstrapCommand")
		}
		if ng.Bottlerocket != nil {
			if err := checkBottlerocketSettings(ng.Bottlerocket.Settings, path); err != nil {
				return err
			}
		}
	} else if err := validateNodeGroupKubeletExtraConfig(ng.KubeletExtraConfig); err != nil {
		return err
	}

	if instanceutils.IsARMGPUInstanceType(SelectInstanceType(ng)) && ng.AMIFamily != NodeImageFamilyBottlerocket {
		return fmt.Errorf("ARM GPU instance types are not supported for unmanaged nodegroups with AMIFamily %s", ng.AMIFamily)
	}

	if err := validateInstancesDistribution(ng); err != nil {
		return err
	}

	if err := validateCPUCredits(ng); err != nil {
		return err
	}

	if err := validateASGSuspendProcesses(ng); err != nil {
		return err
	}

	if ng.ContainerRuntime != nil {
		if *ng.ContainerRuntime != ContainerRuntimeDockerD && *ng.ContainerRuntime != ContainerRuntimeContainerD && *ng.ContainerRuntime != ContainerRuntimeDockerForWindows {
			return fmt.Errorf("only %s, %s and %s are supported for container runtime", ContainerRuntimeContainerD, ContainerRuntimeDockerD, ContainerRuntimeDockerForWindows)
		}
		if clusterVersion := cfg.Metadata.Version; clusterVersion != "" {
			isDockershimDeprecated, err := utils.IsMinVersion(DockershimDeprecationVersion, clusterVersion)
			if err != nil {
				return err
			}
			if *ng.ContainerRuntime != ContainerRuntimeContainerD && isDockershimDeprecated {
				return fmt.Errorf("only %s is supported for container runtime, starting with EKS version %s", ContainerRuntimeContainerD, Version1_24)
			}
		}
		if ng.OverrideBootstrapCommand != nil {
			return fmt.Errorf("overrideBootstrapCommand overwrites container runtime setting; please use --container-runtime in the bootsrap script instead")
		}
	}

	if ng.MaxInstanceLifetime != nil {
		if *ng.MaxInstanceLifetime < OneDay {
			return fmt.Errorf("maximum instance lifetime must have a minimum value of 86,400 seconds (one day), but was: %d", *ng.MaxInstanceLifetime)
		}
	}

	if len(ng.LocalZones) > 0 && len(ng.AvailabilityZones) > 0 {
		return errors.New("cannot specify both localZones and availabilityZones")
	}

	if ng.OutpostARN != "" {
		if err := validateOutpostARN(ng.OutpostARN); err != nil {
			return err
		}
		if cfg.IsControlPlaneOnOutposts() && ng.OutpostARN != cfg.GetOutpost().ControlPlaneOutpostARN {
			return fmt.Errorf("nodeGroup.outpostARN must either be empty or match the control plane's Outpost ARN (%q != %q)", ng.OutpostARN, cfg.GetOutpost().ControlPlaneOutpostARN)
		}
	}

	if cfg.IsControlPlaneOnOutposts() || ng.OutpostARN != "" {
		if ng.InstanceSelector != nil && !ng.InstanceSelector.IsZero() {
			return errors.New("cannot specify instanceSelector for a nodegroup on Outposts")
		}
		const msg = "%s cannot be specified for a nodegroup on Outposts; the AZ defaults to the Outpost AZ"
		if len(ng.AvailabilityZones) > 0 {
			return fmt.Errorf(msg, "availabilityZones")
		}
		if len(ng.LocalZones) > 0 {
			return fmt.Errorf(msg, "localZones")
		}
	}

	return nil
}

func validateOutpostARN(val string) error {
	parsed, err := arn.Parse(val)
	if err != nil {
		return fmt.Errorf("invalid Outpost ARN: %w", err)
	}
	if parsed.Service != "outposts" {
		return fmt.Errorf("invalid Outpost ARN: %q", val)
	}
	return nil
}

// validateNodeGroupLabels uses proper Kubernetes label validation,
// it's designed to make sure users don't pass weird labels to the
// nodes, which would prevent kubelets to startup properly
func validateNodeGroupLabels(labels map[string]string) error {
	// compact version based on:
	// - https://github.com/kubernetes/kubernetes/blob/v1.13.2/cmd/kubelet/app/options/options.go#L257-L267
	// - https://github.com/kubernetes/kubernetes/blob/v1.13.2/pkg/kubelet/apis/well_known_labels.go
	// we cannot import those packages because they break other dependencies

	unknownKubernetesLabels := []string{}

	for label := range labels {
		labelParts := strings.Split(label, "/")

		if len(labelParts) > 2 {
			return fmt.Errorf("node label key %q is of invalid format, can only use one '/' separator", label)
		}

		if errs := validation.IsQualifiedName(label); len(errs) > 0 {
			return fmt.Errorf("label %q is invalid - %v", label, errs)
		}
		if errs := validation.IsValidLabelValue(labels[label]); len(errs) > 0 {
			return fmt.Errorf("label %q has invalid value %q - %v", label, labels[label], errs)
		}

		if len(labelParts) == 2 {
			namespace := labelParts[0]
			if isKubernetesLabel(namespace) && !kubeletapis.IsKubeletLabel(label) {
				unknownKubernetesLabels = append(unknownKubernetesLabels, label)
			}
		}
	}

	if len(unknownKubernetesLabels) > 0 {
		return fmt.Errorf("unknown 'kubernetes.io' or 'k8s.io' labels were specified: %v", unknownKubernetesLabels)
	}
	return nil
}

func isKubernetesLabel(namespace string) bool {
	for _, domain := range []string{"kubernetes.io", "k8s.io"} {
		if namespace == domain || strings.HasSuffix(namespace, "."+domain) {
			return true
		}
	}
	return false
}

func validateNodeGroupIAMWithAddonPolicies(
	policies NodeGroupIAMAddonPolicies,
	fmtFieldConflictErr func(conflictingField string) error,
) error {
	prefix := "withAddonPolicies."
	if IsEnabled(policies.AutoScaler) {
		return fmtFieldConflictErr(prefix + "autoScaler")
	}
	if IsEnabled(policies.ExternalDNS) {
		return fmtFieldConflictErr(prefix + "externalDNS")
	}
	if IsEnabled(policies.CertManager) {
		return fmtFieldConflictErr(prefix + "certManager")
	}
	if IsEnabled(policies.ImageBuilder) {
		return fmtFieldConflictErr(prefix + "imageBuilder")
	}
	if IsEnabled(policies.AppMesh) {
		return fmtFieldConflictErr(prefix + "appMesh")
	}
	if IsEnabled(policies.AppMeshPreview) {
		return fmtFieldConflictErr(prefix + "appMeshPreview")
	}
	if IsEnabled(policies.EBS) {
		return fmtFieldConflictErr(prefix + "ebs")
	}
	if IsEnabled(policies.FSX) {
		return fmtFieldConflictErr(prefix + "fsx")
	}
	if IsEnabled(policies.EFS) {
		return fmtFieldConflictErr(prefix + "efs")
	}
	if IsEnabled(policies.AWSLoadBalancerController) {
		return fmtFieldConflictErr(prefix + "awsLoadBalancerController")
	}
	if IsEnabled(policies.DeprecatedALBIngress) {
		return fmtFieldConflictErr(prefix + "albIngress")
	}
	if IsEnabled(policies.XRay) {
		return fmtFieldConflictErr(prefix + "xRay")
	}
	if IsEnabled(policies.CloudWatch) {
		return fmtFieldConflictErr(prefix + "cloudWatch")
	}
	return nil
}

func validateDeprecatedIAMFields(iam *NodeGroupIAM) error {
	if IsEnabled(iam.WithAddonPolicies.DeprecatedALBIngress) {
		if IsEnabled(iam.WithAddonPolicies.AWSLoadBalancerController) {
			return fmt.Errorf(`"awsLoadBalancerController" and "albIngress" cannot both be configured, ` +
				`please use "awsLoadBalancerController" as "albIngress" is deprecated`)
		}
		logger.Warning("nodegroup.iam.withAddonPolicies.albIngress field is deprecated, please use awsLoadBalancerController instead")
	}

	return nil
}

func validateNodeGroupIAM(iam *NodeGroupIAM, value, fieldName, path string) error {
	if value != "" {
		fmtFieldConflictErr := func(conflictingField string) error {
			return fmt.Errorf("%s.iam.%s and %s.iam.%s cannot be set at the same time", path, fieldName, path, conflictingField)
		}

		if iam.InstanceRoleName != "" {
			return fmtFieldConflictErr("instanceRoleName")
		}
		if iam.AttachPolicy != nil {
			return fmtFieldConflictErr("attachPolicy")
		}
		if len(iam.AttachPolicyARNs) != 0 {
			return fmtFieldConflictErr("attachPolicyARNs")
		}
		if iam.InstanceRolePermissionsBoundary != "" {
			return fmtFieldConflictErr("instanceRolePermissionsBoundary")
		}
		if err := validateNodeGroupIAMWithAddonPolicies(iam.WithAddonPolicies, fmtFieldConflictErr); err != nil {
			return err
		}
	}
	return nil
}

// ValidateManagedNodeGroup validates a ManagedNodeGroup and sets some defaults
func ValidateManagedNodeGroup(index int, ng *ManagedNodeGroup) error {
	normalizeAMIFamily(ng.BaseNodeGroup())
	path := fmt.Sprintf("managedNodeGroups[%d]", index)

	if err := validateNodeGroupBase(ng, path, false); err != nil {
		return err
	}

	if ng.IAM != nil {
		if err := validateNodeGroupIAM(ng.IAM, ng.IAM.InstanceRoleARN, "instanceRoleARN", path); err != nil {
			return err
		}

		errNotSupported := func(field string) error {
			return fmt.Errorf("%s is not supported for Managed Nodes (%s.%s)", field, path, field)
		}

		if ng.IAM.InstanceProfileARN != "" {
			return errNotSupported("instanceProfileARN")
		}
	}

	if ng.OutpostARN != "" {
		return errors.New("Outposts is not supported for managed nodegroups")
	}

	// TODO fix error messages to not use CLI flags
	if ng.MinSize == nil {
		if ng.DesiredCapacity == nil {
			defaultNodeCount := DefaultNodeCount
			ng.MinSize = &defaultNodeCount
		} else {
			ng.MinSize = ng.DesiredCapacity
		}
	} else if ng.DesiredCapacity != nil && *ng.DesiredCapacity < *ng.MinSize {
		return fmt.Errorf("cannot use --nodes-min=%d and --nodes=%d at the same time", *ng.MinSize, *ng.DesiredCapacity)
	}

	// Ensure MaxSize is set, as it is required by the ASG CFN resource
	if ng.MaxSize == nil {
		if ng.DesiredCapacity == nil {
			ng.MaxSize = ng.MinSize
		} else {
			ng.MaxSize = ng.DesiredCapacity
		}
		// MaxSize needs to be greater or equal to 1
		if *ng.MaxSize == 0 {
			defaultMaxSize := DefaultMaxSize
			ng.MaxSize = &defaultMaxSize
		}
	} else if ng.DesiredCapacity != nil && *ng.DesiredCapacity > *ng.MaxSize {
		return fmt.Errorf("cannot use --nodes-max=%d and --nodes=%d at the same time", *ng.MaxSize, *ng.DesiredCapacity)
	} else if *ng.MaxSize < *ng.MinSize {
		return fmt.Errorf("cannot use --nodes-min=%d and --nodes-max=%d at the same time", *ng.MinSize, *ng.MaxSize)
	}

	if ng.DesiredCapacity == nil {
		ng.DesiredCapacity = ng.MinSize
	}

	if ng.UpdateConfig != nil {
		if ng.UpdateConfig.MaxUnavailable == nil && ng.UpdateConfig.MaxUnavailablePercentage == nil {
			return fmt.Errorf("invalid UpdateConfig: maxUnavailable or maxUnavailablePercentage must be defined")
		}
		if ng.UpdateConfig.MaxUnavailable != nil && ng.UpdateConfig.MaxUnavailablePercentage != nil {
			return fmt.Errorf("cannot use maxUnavailable=%d and maxUnavailablePercentage=%d at the same time", *ng.UpdateConfig.MaxUnavailable, *ng.UpdateConfig.MaxUnavailablePercentage)
		}
		if aws.ToInt(ng.UpdateConfig.MaxUnavailable) > aws.ToInt(ng.MaxSize) {
			return fmt.Errorf("maxUnavailable=%d cannot be greater than maxSize=%d", *ng.UpdateConfig.MaxUnavailable, *ng.MaxSize)
		}
	}

	if IsEnabled(ng.SecurityGroups.WithLocal) || IsEnabled(ng.SecurityGroups.WithShared) {
		return errors.Errorf("securityGroups.withLocal and securityGroups.withShared are not supported for managed nodegroups (%s.securityGroups)", path)
	}

	if ng.InstanceType != "" {
		if len(ng.InstanceTypes) > 0 {
			return errors.Errorf("only one of instanceType or instanceTypes can be specified (%s)", path)
		}
		if !ng.InstanceSelector.IsZero() {
			return errors.Errorf("cannot set instanceType when instanceSelector is specified (%s)", path)
		}
	}

	if ng.AMIFamily == NodeImageFamilyBottlerocket {
		fieldNotSupported := func(field string) error {
			return &unsupportedFieldError{
				ng:    ng.NodeGroupBase,
				path:  path,
				field: field,
			}
		}
		if ng.PreBootstrapCommands != nil {
			return fieldNotSupported("preBootstrapCommands")
		}
		if ng.OverrideBootstrapCommand != nil {
			return fieldNotSupported("overrideBootstrapCommand")
		}
	}

	// Windows doesn't use overrideBootstrapCommand, as it always uses bootstrapping script that comes with Windows AMIs
	if IsWindowsImage(ng.AMIFamily) {
		fieldNotSupported := func(field string) error {
			return &unsupportedFieldError{
				ng:    ng.NodeGroupBase,
				path:  path,
				field: field,
			}
		}
		if ng.OverrideBootstrapCommand != nil {
			return fieldNotSupported("overrideBootstrapCommand")
		}
	}

	if err := validateTaints(ng.Taints); err != nil {
		return err
	}

	switch {
	case ng.LaunchTemplate != nil:
		if ng.LaunchTemplate.ID == "" {
			return errors.Errorf("launchTemplate.id is required if launchTemplate is set (%s.%s)", path, "launchTemplate")
		}

		if ng.LaunchTemplate.Version != nil {
			// TODO support `latest` and `default`
			versionNumber, err := strconv.ParseInt(*ng.LaunchTemplate.Version, 10, 64)
			if err != nil {
				return errors.Wrap(err, "invalid launch template version")
			}
			if versionNumber < 1 {
				return errors.Errorf("launchTemplate.version must be >= 1 (%s.%s)", path, "launchTemplate.version")
			}
		}

		if ng.InstanceType != "" || ng.AMI != "" || IsEnabled(ng.SSH.Allow) || IsEnabled(ng.SSH.EnableSSM) || len(ng.SSH.SourceSecurityGroupIDs) > 0 ||
			ng.VolumeSize != nil || len(ng.PreBootstrapCommands) > 0 || ng.OverrideBootstrapCommand != nil ||
			len(ng.SecurityGroups.AttachIDs) > 0 || ng.InstanceName != "" || ng.InstancePrefix != "" || ng.MaxPodsPerNode != 0 ||
			IsDisabled(ng.DisableIMDSv1) || IsEnabled(ng.DisablePodIMDS) || ng.Placement != nil {

			incompatibleFields := []string{
				"instanceType", "ami", "ssh.allow", "ssh.enableSSM", "ssh.sourceSecurityGroupIds", "securityGroups",
				"volumeSize", "instanceName", "instancePrefix", "maxPodsPerNode", "disableIMDSv1",
				"disablePodIMDS", "preBootstrapCommands", "overrideBootstrapCommand", "placement",
			}
			return errors.Errorf("cannot set %s in managedNodeGroup when a launch template is supplied", strings.Join(incompatibleFields, ", "))
		}

	case ng.AMI != "":
		if !IsAMI(ng.AMI) {
			return errors.Errorf("invalid AMI %q (%s.%s)", ng.AMI, path, "ami")
		}
		if ng.AMIFamily == "" {
			return errors.Errorf("when using a custom AMI, amiFamily needs to be explicitly set via config file or via --node-ami-family flag")
		}
		if ng.AMIFamily != NodeImageFamilyAmazonLinux2 && ng.AMIFamily != NodeImageFamilyUbuntu1804 && ng.AMIFamily != NodeImageFamilyUbuntu2004 {
			return errors.Errorf("cannot set amiFamily to %s when using a custom AMI for managed nodes, only %s, %s and %s are supported", ng.AMIFamily, NodeImageFamilyAmazonLinux2, NodeImageFamilyUbuntu1804, NodeImageFamilyUbuntu2004)
		}
		if ng.OverrideBootstrapCommand == nil {
			return errors.Errorf("%s.overrideBootstrapCommand is required when using a custom AMI (%s.ami)", path, path)
		}
		notSupportedWithCustomAMIErr := func(field string) error {
			return errors.Errorf("%s.%s is not supported when using a custom AMI (%s.ami)", path, field, path)
		}
		if ng.MaxPodsPerNode != 0 {
			return notSupportedWithCustomAMIErr("maxPodsPerNode")
		}
		if ng.SSH != nil && IsEnabled(ng.SSH.EnableSSM) {
			return notSupportedWithCustomAMIErr("enableSSM")
		}
		if ng.ReleaseVersion != "" {
			return notSupportedWithCustomAMIErr("releaseVersion")
		}

	case ng.OverrideBootstrapCommand != nil:
		return errors.Errorf("%s.overrideBootstrapCommand can only be set when a custom AMI (%s.ami) is specified", path, path)
	}

	return nil
}

func normalizeAMIFamily(ng *NodeGroupBase) {
	for _, family := range supportedAMIFamilies() {
		if strings.EqualFold(ng.AMIFamily, family) {
			ng.AMIFamily = family
			return
		}
	}
}

func validateInstancesDistribution(ng *NodeGroup) error {
	hasInstanceSelector := ng.InstanceSelector != nil && !ng.InstanceSelector.IsZero()
	if ng.InstancesDistribution == nil && !hasInstanceSelector {
		return nil
	}

	if ng.InstanceType != "" && ng.InstanceType != "mixed" {
		makeError := func(featureStr string) error {
			return errors.Errorf(`instanceType should be "mixed" or unset when using the %s feature`, featureStr)
		}
		if ng.InstancesDistribution != nil {
			return makeError("instances distribution")
		}
		return makeError("instance selector")
	}

	if ng.InstancesDistribution == nil {
		return nil
	}

	distribution := ng.InstancesDistribution
	if len(distribution.InstanceTypes) == 0 && !hasInstanceSelector {
		return fmt.Errorf("at least two instance types have to be specified for mixed nodegroups")
	}

	if !hasInstanceSelector {
		uniqueInstanceTypes := make(map[string]struct{})
		for _, instanceType := range distribution.InstanceTypes {
			uniqueInstanceTypes[instanceType] = struct{}{}
		}

		if len(uniqueInstanceTypes) > 20 {
			return fmt.Errorf("mixed nodegroups should have between 1 and 20 different instance types")
		}
	}

	if distribution.OnDemandBaseCapacity != nil && *distribution.OnDemandBaseCapacity < 0 {
		return fmt.Errorf("onDemandBaseCapacity should be 0 or more")
	}

	if distribution.OnDemandPercentageAboveBaseCapacity != nil && (*distribution.OnDemandPercentageAboveBaseCapacity < 0 || *distribution.OnDemandPercentageAboveBaseCapacity > 100) {
		return fmt.Errorf("percentageAboveBase should be between 0 and 100")
	}

	if distribution.SpotInstancePools != nil && (*distribution.SpotInstancePools < 1 || *distribution.SpotInstancePools > 20) {
		return fmt.Errorf("spotInstancePools should be between 1 and 20")
	}

	if distribution.SpotInstancePools != nil && distribution.SpotAllocationStrategy != nil && *distribution.SpotAllocationStrategy == SpotAllocationStrategyCapacityOptimized {
		return fmt.Errorf("spotInstancePools cannot be specified when also specifying spotAllocationStrategy: %s", SpotAllocationStrategyCapacityOptimized)
	}

	if distribution.SpotInstancePools != nil && distribution.SpotAllocationStrategy != nil && (*distribution.SpotAllocationStrategy == SpotAllocationStrategyCapacityOptimized || *distribution.SpotAllocationStrategy == SpotAllocationStrategyCapacityOptimizedPrioritized) {
		return fmt.Errorf("spotInstancePools cannot be specified when also specifying spotAllocationStrategy: %s", SpotAllocationStrategyCapacityOptimizedPrioritized)
	}

	if distribution.SpotAllocationStrategy != nil {
		if err := validateSpotAllocationStrategy(*distribution.SpotAllocationStrategy); err != nil {
			return err
		}
	}

	return nil
}

func validateCPUCredits(ng *NodeGroup) error {
	isTInstance := false
	instanceTypes := []string{ng.InstanceType}

	if ng.CPUCredits == nil {
		return nil
	}

	if ng.InstancesDistribution != nil {
		instanceTypes = ng.InstancesDistribution.InstanceTypes
	}

	for _, instanceType := range instanceTypes {
		if strings.HasPrefix(instanceType, "t") {
			isTInstance = true
		}
	}

	if !isTInstance {
		return fmt.Errorf("cpuCredits option set for nodegroup, but it has no t2/t3 instance types")
	}

	if strings.ToLower(*ng.CPUCredits) != "unlimited" && strings.ToLower(*ng.CPUCredits) != "standard" {
		return fmt.Errorf("cpuCredits option accepts only one of 'standard' or 'unlimited'")
	}

	return nil
}

func validateASGSuspendProcesses(ng *NodeGroup) error {
	// Processes list taken from here: https://docs.aws.amazon.com/autoscaling/ec2/APIReference/API_SuspendProcesses.html
	for _, proc := range ng.ASGSuspendProcesses {
		switch proc {
		case
			"Launch",
			"Terminate",
			"AddToLoadBalancer",
			"AlarmNotification",
			"AZRebalance",
			"HealthCheck",
			"InstanceRefresh",
			"ReplaceUnhealthy",
			"ScheduledActions":
			continue
		default:
			return fmt.Errorf("asgSuspendProcesses contains invalid process name '%s'", proc)
		}
	}
	return nil
}

func validateNodeGroupSSH(SSH *NodeGroupSSH) error {
	numSSHFlagsEnabled := countEnabledFields(
		SSH.PublicKeyPath,
		SSH.PublicKey,
		SSH.PublicKeyName)

	if numSSHFlagsEnabled > 1 {
		return errors.New("only one of publicKeyName, publicKeyPath or publicKey can be specified for SSH per node-group")
	}
	return nil
}

func countEnabledFields(fields ...*string) int {
	count := 0
	for _, flag := range fields {
		if flag != nil && *flag != "" {
			count++
		}
	}
	return count
}

func validateNodeGroupKubeletExtraConfig(kubeletConfig *InlineDocument) error {
	if kubeletConfig == nil {
		return nil
	}

	var kubeletForbiddenFields = map[string]struct{}{
		"kind":               {},
		"apiVersion":         {},
		"address":            {},
		"clusterDomain":      {},
		"authentication":     {},
		"authorization":      {},
		"serverTLSBootstrap": {},
	}

	for k := range *kubeletConfig {
		if _, exists := kubeletForbiddenFields[k]; exists {
			return fmt.Errorf("cannot override %q in kubelet config, as it's critical to eksctl functionality", k)
		}
	}
	return nil
}

func isSupportedAMIFamily(imageFamily string) bool {
	for _, image := range supportedAMIFamilies() {
		if imageFamily == image {
			return true
		}
	}
	return false
}

// IsWindowsImage reports whether the AMI family is for Windows
func IsWindowsImage(imageFamily string) bool {
	switch imageFamily {
	case NodeImageFamilyWindowsServer2019CoreContainer,
		NodeImageFamilyWindowsServer2019FullContainer,
		NodeImageFamilyWindowsServer2022CoreContainer,
		NodeImageFamilyWindowsServer2022FullContainer:
		return true

	default:
		return false
	}
}

func validateCIDRs(cidrs []string) ([]string, error) {
	var validCIDRs []string
	for _, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		validCIDRs = append(validCIDRs, ipNet.String())
	}
	return validCIDRs, nil
}

func validateTaints(ngTaints []NodeGroupTaint) error {
	for _, t := range ngTaints {
		if err := taints.Validate(corev1.Taint{
			Key:    t.Key,
			Value:  t.Value,
			Effect: t.Effect,
		}); err != nil {
			return err
		}
	}
	return nil
}

// ReservedProfileNamePrefix defines the Fargate profile name prefix reserved
// for AWS, and which therefore, cannot be used by users. AWS' API should
// reject the creation of profiles starting with this prefix, but we eagerly
// validate this client-side.
const ReservedProfileNamePrefix = "eks-"

// Validate validates this FargateProfile object.
func (fp FargateProfile) Validate() error {
	if fp.Name == "" {
		return errors.New("invalid Fargate profile: empty name")
	}
	if strings.HasPrefix(fp.Name, ReservedProfileNamePrefix) {
		return fmt.Errorf("invalid Fargate profile %q: name should NOT start with %q", fp.Name, ReservedProfileNamePrefix)
	}
	if len(fp.Selectors) == 0 {
		return fmt.Errorf("invalid Fargate profile %q: no profile selector", fp.Name)
	}
	for i, selector := range fp.Selectors {
		if err := selector.Validate(); err != nil {
			return errors.Wrapf(err, "invalid Fargate profile %q: invalid profile selector at index #%v", fp.Name, i)
		}
	}
	return nil
}

// Validate validates this FargateProfileSelector object.
func (fps FargateProfileSelector) Validate() error {
	if fps.Namespace == "" {
		return errors.New("empty namespace")
	}
	return nil
}

func checkBottlerocketSettings(doc *InlineDocument, path string) error {
	if doc == nil {
		return nil
	}

	overlapErr := func(key, ngField string) error {
		return errors.Errorf("invalid Bottlerocket setting: use %s.%s instead (path=%s)", path, ngField, key)
	}

	// Dig into kubernetes settings if provided.
	kubeVal, ok := (*doc)["kubernetes"]
	if !ok {
		return nil
	}

	kube, ok := kubeVal.(map[string]interface{})
	if !ok {
		return errors.New("invalid kubernetes settings provided: expected a map of settings")
	}

	checkMapping := map[string]string{
		"node-labels":    "labels",
		"node-taints":    "taints",
		"max-pods":       "maxPodsPerNode",
		"cluster-dns-ip": "clusterDNS",
	}

	for checkKey, shouldUse := range checkMapping {
		_, ok := kube[checkKey]
		if ok {
			return overlapErr(path+".kubernetes."+checkKey, shouldUse)
		}
	}

	return nil
}

func validateAvailabilityZones(azList []string) error {
	count := len(azList)
	switch {
	case count == 0:
		return nil
	case count < MinRequiredAvailabilityZones:
		return ErrTooFewAvailabilityZones(azList)
	default:
		return nil
	}
}

func ErrTooFewAvailabilityZones(azs []string) error {
	return fmt.Errorf("only %d zone(s) specified %v, %d are required (can be non-unique)", len(azs), azs, MinRequiredAvailabilityZones)
}

func ValidateSecretsEncryption(clusterConfig *ClusterConfig) error {
	if clusterConfig.SecretsEncryption == nil {
		return nil
	}

	if clusterConfig.SecretsEncryption.KeyARN == "" {
		return errors.New("field secretsEncryption.keyARN is required for enabling secrets encryption")
	}

	if _, err := arn.Parse(clusterConfig.SecretsEncryption.KeyARN); err != nil {
		return errors.Wrapf(err, "invalid ARN in secretsEncryption.keyARN: %q", clusterConfig.SecretsEncryption.KeyARN)
	}
	return nil
}

func validateIAMIdentityMappings(clusterConfig *ClusterConfig) error {
	for _, mapping := range clusterConfig.IAMIdentityMappings {
		if err := mapping.Validate(); err != nil {
			return err
		}
	}
	return nil
}
