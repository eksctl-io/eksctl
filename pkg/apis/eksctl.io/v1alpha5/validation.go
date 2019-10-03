package v1alpha5

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/util/validation"
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
		"https://docs.aws.amazon.com/eks/latest/userguide/cluster-endpoint.html#private-access " +
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
		if len(sa.AttachPolicyARNs) == 0 && sa.AttachPolicy == nil {
			return fmt.Errorf("%s.attachPolicyARNs or %s.attachPolicy must be set", path, path)
		}
	}

	ngNames := nameSet{}
	for i, ng := range cfg.NodeGroups {
		path := fmt.Sprintf("nodeGroups[%d]", i)
		if ng.Name == "" {
			return fmt.Errorf("%s.name must be set", path)
		}
		if ok, err := ngNames.checkUnique(path+".name", ng.NameString()); !ok {
			return err
		}
	}

	if cfg.HasClusterCloudWatchLogging() {
		for i, logType := range cfg.CloudWatch.ClusterLogging.EnableTypes {
			isUnknown := true
			for _, knownLogType := range SupportedCloudWatchClusterLogTypes() {
				if logType == knownLogType {
					isUnknown = false
				}
			}
			if isUnknown {
				return fmt.Errorf("log type %q (cloudWatch.clusterLogging.enableTypes[%d]) is unknown", logType, i)
			}
		}
	}

	if !cfg.HasClusterEndpointAccess() {
		return ErrClusterEndpointNoAccess
	}
	return nil
}

// PrivateOnlyUseUtilsMsg returns a message that indicates that the operation must be done using
// eksctl utils update-cluster-endpoints
func PrivateOnlyUseUtilsMsg() string {
	return "eksctl cannot join worker nodes to the EKS cluster when public access isn't allowed. " +
		"use 'eksctl utils update-cluster-endpoints ...' after creating cluster with default access"
}

//ValidateClusterEndpointConfig checks the endpoint configuration for potential issues
func (c *ClusterConfig) ValidateClusterEndpointConfig() error {
	endpts := c.VPC.ClusterEndpoints
	if NoAccess(endpts) {
		return ErrClusterEndpointNoAccess
	}
	if PrivateOnly(endpts) {
		return ErrClusterEndpointPrivateOnly
	}
	return nil
}

//NoAccess returns true if neither public are private cluster endpoint access is enabled and false otherwise
func NoAccess(ces *ClusterEndpoints) bool {
	return !(*ces.PublicAccess || *ces.PrivateAccess)
}

//PrivateOnly returns true if public cluster endpoint access is disabled and private cluster endpoint access is enabled, and false otherwise
func PrivateOnly(ces *ClusterEndpoints) bool {
	return !*ces.PublicAccess && *ces.PrivateAccess
}

// ValidateNodeGroup checks compatible fields of a given nodegroup
func ValidateNodeGroup(i int, ng *NodeGroup) error {
	path := fmt.Sprintf("nodeGroups[%d]", i)

	if ng.VolumeSize == nil {
		errCantSet := func(field string) error {
			return fmt.Errorf("%s.%s cannot be set without %s.volumeSize", path, field, path)
		}
		if IsSetAndNonEmptyString(ng.VolumeType) {
			return errCantSet("volumeType")
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

	if ng.VolumeType != nil && *ng.VolumeType == NodeVolumeTypeIO1 {
		if ng.VolumeIOPS == nil {
			return fmt.Errorf("%s.volumeIOPS is required for %s volume type", path, NodeVolumeTypeIO1)
		}
	} else if ng.VolumeIOPS != nil {
		return fmt.Errorf("%s.volumeIOPS is only supported for %s volume type", path, NodeVolumeTypeIO1)
	}

	if ng.VolumeEncrypted == nil || IsDisabled(ng.VolumeEncrypted) {
		if IsSetAndNonEmptyString(ng.VolumeKmsKeyID) {
			return fmt.Errorf("%s.VolumeKmsKeyID can not be set without %s.VolumeEncrypted enabled explicitly", path, path)
		}
	}

	if ng.IAM != nil {
		if err := validateNodeGroupIAM(i, ng, ng.IAM.InstanceProfileARN, "instanceProfileARN", path); err != nil {
			return err
		}
		if err := validateNodeGroupIAM(i, ng, ng.IAM.InstanceRoleARN, "instanceRoleARN", path); err != nil {
			return err
		}

		if err := ValidateNodeGroupLabels(ng); err != nil {
			return err
		}

		if err := validateNodeGroupSSH(ng.SSH); err != nil {
			return fmt.Errorf("only one ssh public key can be specified per node-group")
		}
	}

	if ng.IsWindows() {
		fieldNotSupported := func(field string) error {
			return fmt.Errorf("%s is not supported for Windows node groups (path=%s.%s)", field, path, field)
		}
		if ng.KubeletExtraConfig != nil {
			return fieldNotSupported("kubeletExtraConfig")
		}
		if ng.PreBootstrapCommands != nil {
			return fieldNotSupported("preBootstrapCommands")
		}
		if ng.OverrideBootstrapCommand != nil {
			return fieldNotSupported("overrideBootstrapCommand")
		}

	} else if err := validateNodeGroupKubeletExtraConfig(ng.KubeletExtraConfig); err != nil {
		return err
	}

	if err := validateInstancesDistribution(ng); err != nil {
		return err
	}

	return nil
}

// ValidateNodeGroupLabels uses proper Kubernetes label validation,
// it's designed to make sure users don't pass weird labels to the
// nodes, which would prevent kubelets to startup properly
func ValidateNodeGroupLabels(ng *NodeGroup) error {
	// compact version based on:
	// - https://github.com/kubernetes/kubernetes/blob/v1.13.2/cmd/kubelet/app/options/options.go#L257-L267
	// - https://github.com/kubernetes/kubernetes/blob/v1.13.2/pkg/kubelet/apis/well_known_labels.go
	// we cannot import those packages because they break other dependencies

	unknownKubernetesLabels := []string{}

	for l := range ng.Labels {
		labelParts := strings.Split(l, "/")

		if len(labelParts) > 2 {
			return fmt.Errorf("node label key %q is of invalid format, can only use one '/' separator", l)
		}

		if errs := validation.IsQualifiedName(l); len(errs) > 0 {
			return fmt.Errorf("label %q is invalid - %v", l, errs)
		}
		if errs := validation.IsValidLabelValue(ng.Labels[l]); len(errs) > 0 {
			return fmt.Errorf("label %q has invalid value %q - %v", l, ng.Labels[l], errs)
		}

		isKubernetesLabel := false
		allowedKubeletNamespace := false
		if len(labelParts) == 2 {
			ns := labelParts[0]

			for _, domain := range []string{"kubernetes.io", "k8s.io"} {
				if ns == domain || strings.HasSuffix(ns, "."+domain) {
					isKubernetesLabel = true
				}
			}

			for _, domain := range []string{"kubelet.kubernetes.io", "node.kubernetes.io", "node-role.kubernetes.io"} {
				if ns == domain || strings.HasSuffix(ns, "."+domain) {
					allowedKubeletNamespace = true
				}
			}

			if isKubernetesLabel && !allowedKubeletNamespace {
				switch l {
				case
					"kubernetes.io/hostname",
					"kubernetes.io/instance-type",
					"kubernetes.io/os",
					"kubernetes.io/arch",
					"beta.kubernetes.io/instance-type",
					"beta.kubernetes.io/os",
					"beta.kubernetes.io/arch",
					"failure-domain.beta.kubernetes.io/zone",
					"failure-domain.beta.kubernetes.io/region",
					"failure-domain.kubernetes.io/zone",
					"failure-domain.kubernetes.io/region":
				default:
					unknownKubernetesLabels = append(unknownKubernetesLabels, l)
				}
			}
		}
	}

	if len(unknownKubernetesLabels) > 0 {
		return fmt.Errorf("unknown 'kubernetes.io' or 'k8s.io' labels were specified: %v", unknownKubernetesLabels)
	}
	return nil
}

func validateNodeGroupIAM(i int, ng *NodeGroup, value, fieldName, path string) error {
	if value != "" {
		fmtFieldConflictErr := func(conflictingField string) error {
			return fmt.Errorf("%s.iam.%s and %s.iam.%s cannot be set at the same time", path, fieldName, path, conflictingField)
		}

		if ng.IAM.InstanceRoleName != "" {
			return fmtFieldConflictErr("instanceRoleName")
		}
		if len(ng.IAM.AttachPolicyARNs) != 0 {
			return fmtFieldConflictErr("attachPolicyARNs")
		}
		prefix := "withAddonPolicies."
		if IsEnabled(ng.IAM.WithAddonPolicies.AutoScaler) {
			return fmtFieldConflictErr(prefix + "autoScaler")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.ExternalDNS) {
			return fmtFieldConflictErr(prefix + "externalDNS")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.CertManager) {
			return fmtFieldConflictErr(prefix + "certManager")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.ImageBuilder) {
			return fmtFieldConflictErr(prefix + "imageBuilder")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.AppMesh) {
			return fmtFieldConflictErr(prefix + "appMesh")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.EBS) {
			return fmtFieldConflictErr(prefix + "ebs")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.FSX) {
			return fmtFieldConflictErr(prefix + "fsx")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.EFS) {
			return fmtFieldConflictErr(prefix + "efs")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.ALBIngress) {
			return fmtFieldConflictErr(prefix + "albIngress")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.XRay) {
			return fmtFieldConflictErr(prefix + "xRay")
		}
		if IsEnabled(ng.IAM.WithAddonPolicies.CloudWatch) {
			return fmtFieldConflictErr(prefix + "cloudWatch")
		}
	}
	return nil
}

func validateInstancesDistribution(ng *NodeGroup) error {
	if ng.InstancesDistribution == nil {
		return nil
	}

	if ng.InstanceType != "" && ng.InstanceType != "mixed" {
		return fmt.Errorf("instanceType should be \"mixed\" or unset when using the mixed instances feature")
	}

	distribution := ng.InstancesDistribution
	if distribution.InstanceTypes == nil || len(distribution.InstanceTypes) == 0 {
		return fmt.Errorf("at least two instance types have to be specified for mixed nodegroups")
	}

	allInstanceTypes := make(map[string]bool)
	for _, instanceType := range distribution.InstanceTypes {
		allInstanceTypes[instanceType] = true
	}

	if len(allInstanceTypes) < 2 || len(allInstanceTypes) > 20 {
		return fmt.Errorf("mixed nodegroups should have between 2 and 20 different instance types")
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

	return nil
}

func validateNodeGroupSSH(SSH *NodeGroupSSH) error {
	if SSH == nil {
		return nil
	}
	numSSHFlagsEnabled := countEnabledFields(
		SSH.PublicKeyPath,
		SSH.PublicKey,
		SSH.PublicKeyName)

	if numSSHFlagsEnabled > 1 {
		return fmt.Errorf("only one ssh public key can be specified per node-group")
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
