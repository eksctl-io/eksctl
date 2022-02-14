package v1alpha5

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
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

	for i, ng := range cfg.NodeGroups {
		path := fmt.Sprintf("nodeGroups[%d]", i)
		if err := validateNg(ng.NodeGroupBase, path); err != nil {
			return err
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

	if err := cfg.ValidateVPCConfig(); err != nil {
		return err
	}

	if cfg.SecretsEncryption != nil && cfg.SecretsEncryption.KeyARN == "" {
		return errors.New("field secretsEncryption.keyARN is required for enabling secrets encryption")
	}

	if err := validateKarpenterConfig(cfg); err != nil {
		return fmt.Errorf("failed to validate karpenter config: %w", err)
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
	}

	// manageSharedNodeSecurityGroupRules cannot be disabled if using eksctl managed security groups
	if c.VPC.SharedNodeSecurityGroup == "" && IsDisabled(c.VPC.ManageSharedNodeSecurityGroupRules) {
		return errors.New("vpc.manageSharedNodeSecurityGroupRules must be enabled when using ekstcl-managed security groups")
	}
	return nil
}

func (c *ClusterConfig) unsupportedVPCCNIAddonVersion() (bool, error) {
	for _, addon := range c.Addons {
		if addon.Name == VPCCNIAddon {
			if addon.Version == "" {
				addon.Version = minimumVPCCNIVersionForIPv6
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
	if !c.HasClusterEndpointAccess() {
		return ErrClusterEndpointNoAccess
	}
	endpts := c.VPC.ClusterEndpoints
	if noAccess(endpts) {
		return ErrClusterEndpointNoAccess
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
				return errors.Wrap(err, "invalid value in privateCluster.additionalEndpointServices")
			}
		}

		if c.VPC != nil && c.VPC.ClusterEndpoints == nil {
			c.VPC.ClusterEndpoints = &ClusterEndpoints{}
		}
		// public access is initially enabled to allow running operations that access the Kubernetes API
		c.VPC.ClusterEndpoints.PublicAccess = Enabled()
		c.VPC.ClusterEndpoints.PrivateAccess = Enabled()
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
			return fmt.Errorf("service ipv4 cidr is not supported with IPv6")
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
	return !(*ces.PublicAccess || *ces.PrivateAccess)
}

// PrivateOnly returns true if public cluster endpoint access is disabled and private cluster endpoint access is enabled, and false otherwise
func PrivateOnly(ces *ClusterEndpoints) bool {
	return !*ces.PublicAccess && *ces.PrivateAccess
}

func validateNodeGroupBase(np NodePool, path string) error {
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

	if err := validateVolumeOpts(ng, path); err != nil {
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

	if ng.AMIFamily != "" && !isSupportedAMIFamily(ng.AMIFamily) {
		return fmt.Errorf("AMI Family %s is not supported - use one of: %s", ng.AMIFamily, strings.Join(supportedAMIFamilies(), ", "))
	}

	if ng.SSH != nil {
		if enableSSM := ng.SSH.EnableSSM; enableSSM != nil {
			if !*enableSSM {
				return errors.New("SSM agent is now built into EKS AMIs and cannot be disabled")
			}
			logger.Warning("SSM is now enabled by default; `ssh.enableSSM` is deprecated and will be removed in a future release")
		}
	}

	if instanceutils.IsGPUInstanceType(SelectInstanceType(np)) && (ng.AMIFamily != NodeImageFamilyAmazonLinux2 && ng.AMIFamily != "") {
		return errors.Errorf("GPU instance types are not supported for %s", ng.AMIFamily)
	}

	return nil
}

func validateVolumeOpts(ng *NodeGroupBase, path string) error {
	if ng.VolumeType != nil {
		if ng.VolumeIOPS != nil && !(*ng.VolumeType == NodeVolumeTypeIO1 || *ng.VolumeType == NodeVolumeTypeGP3) {
			return fmt.Errorf("%s.volumeIOPS is only supported for %s and %s volume types", path, NodeVolumeTypeIO1, NodeVolumeTypeGP3)
		}

		if *ng.VolumeType == NodeVolumeTypeIO1 {
			if ng.VolumeIOPS != nil && !(*ng.VolumeIOPS >= MinIO1Iops && *ng.VolumeIOPS <= MaxIO1Iops) {
				return fmt.Errorf("value for %s.volumeIOPS must be within range %d-%d", path, MinIO1Iops, MaxIO1Iops)
			}
		}

		if ng.VolumeThroughput != nil && *ng.VolumeType != NodeVolumeTypeGP3 {
			return fmt.Errorf("%s.volumeThroughput is only supported for %s volume type", path, NodeVolumeTypeGP3)
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
func ValidateNodeGroup(i int, ng *NodeGroup) error {
	path := fmt.Sprintf("nodeGroups[%d]", i)
	if err := validateNodeGroupBase(ng, path); err != nil {
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

	if ng.Bottlerocket != nil && ng.AMIFamily != NodeImageFamilyBottlerocket {
		return fmt.Errorf(`bottlerocket config can only be used with amiFamily "Bottlerocket" but found %s (path=%s.bottlerocket)`,
			ng.AMIFamily, path)
	}

	if IsWindowsImage(ng.AMIFamily) || ng.AMIFamily == NodeImageFamilyBottlerocket {
		fieldNotSupported := func(field string) error {
			return &unsupportedFieldError{
				ng:    ng.NodeGroupBase,
				path:  path,
				field: field,
			}
		}
		if ng.KubeletExtraConfig != nil {
			return fieldNotSupported("kubeletExtraConfig")
		}
		if ng.AMIFamily == NodeImageFamilyBottlerocket && ng.PreBootstrapCommands != nil {
			return fieldNotSupported("preBootstrapCommands")

		}
		if ng.OverrideBootstrapCommand != nil {
			return fieldNotSupported("overrideBootstrapCommand")
		}

	} else if err := validateNodeGroupKubeletExtraConfig(ng.KubeletExtraConfig); err != nil {
		return err
	}

	if ng.AMIFamily == NodeImageFamilyBottlerocket && ng.Bottlerocket != nil {
		err := checkBottlerocketSettings(ng.Bottlerocket.Settings, path)
		if err != nil {
			return err
		}
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
		if *ng.ContainerRuntime == ContainerRuntimeContainerD && ng.AMIFamily != NodeImageFamilyAmazonLinux2 {
			// check if it's dockerd or containerd
			return fmt.Errorf("%s as runtime is only support for AL2 ami family", ContainerRuntimeContainerD)
		}
		if *ng.ContainerRuntime != ContainerRuntimeDockerD && *ng.ContainerRuntime != ContainerRuntimeContainerD {
			return fmt.Errorf("only %s and %s are supported for container runtime", ContainerRuntimeContainerD, ContainerRuntimeDockerD)
		}
	}

	if ng.MaxInstanceLifetime != nil {
		if *ng.MaxInstanceLifetime < OneDay {
			return fmt.Errorf("maximum instance lifetime must have a minimum value of 86,400 seconds (one day), but was: %d", *ng.MaxInstanceLifetime)
		}
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
func ValidateManagedNodeGroup(ng *ManagedNodeGroup, index int) error {
	switch ng.AMIFamily {
	case NodeImageFamilyAmazonLinux2, NodeImageFamilyBottlerocket, NodeImageFamilyUbuntu1804, NodeImageFamilyUbuntu2004:

	default:
		return errors.Errorf("%q is not supported for managed nodegroups", ng.AMIFamily)
	}

	path := fmt.Sprintf("managedNodeGroups[%d]", index)

	if err := validateNodeGroupBase(ng, path); err != nil {
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
		if aws.IntValue(ng.UpdateConfig.MaxUnavailable) > aws.IntValue(ng.MaxSize) {
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
			IsEnabled(ng.DisableIMDSv1) || IsEnabled(ng.DisablePodIMDS) || ng.Placement != nil {

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
		if ng.AMIFamily != NodeImageFamilyAmazonLinux2 {
			return errors.Errorf("cannot set amiFamily to %s when using a custom AMI", ng.AMIFamily)
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
		if !isSpotAllocationStrategySupported(*distribution.SpotAllocationStrategy) {
			return fmt.Errorf("spotAllocationStrategy should be one of: %v", strings.Join(supportedSpotAllocationStrategies(), ", "))
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
		NodeImageFamilyWindowsServer2004CoreContainer,
		NodeImageFamilyWindowsServer20H2CoreContainer:
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
