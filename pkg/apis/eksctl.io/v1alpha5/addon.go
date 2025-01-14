package v1alpha5

import (
	"encoding/json"
	"fmt"
	"strings"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"sigs.k8s.io/yaml"
)

// Values for core addons
const (
	minimumVPCCNIVersionForIPv6 = "1.10.0"

	VPCCNIAddon           = "vpc-cni"
	KubeProxyAddon        = "kube-proxy"
	CoreDNSAddon          = "coredns"
	PodIdentityAgentAddon = "eks-pod-identity-agent"
	MetricsServerAddon    = "metrics-server"
	AWSEBSCSIDriverAddon  = "aws-ebs-csi-driver"
	AWSEFSCSIDriverAddon  = "aws-efs-csi-driver"
)

// Addon holds the EKS addon configuration
type Addon struct {
	// +required
	Name string `json:"name,omitempty"`
	// +optional
	Version string `json:"version,omitempty"`
	// +optional
	ServiceAccountRoleARN string `json:"serviceAccountRoleARN,omitempty"`
	// list of ARNs of the IAM policies to attach
	// +optional
	AttachPolicyARNs []string `json:"attachPolicyARNs,omitempty"`
	// AttachPolicy holds a policy document to attach
	// +optional
	AttachPolicy InlineDocument `json:"attachPolicy,omitempty"`
	// ARN of the permissions' boundary to associate
	// +optional
	PermissionsBoundary string `json:"permissionsBoundary,omitempty"`
	// WellKnownPolicies for attaching common IAM policies
	WellKnownPolicies WellKnownPolicies `json:"wellKnownPolicies,omitempty"`
	// The metadata to apply to the cluster to assist with categorization and organization.
	// Each tag consists of a key and an optional value, both of which you define.
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
	// ResolveConflicts determines how to resolve field value conflicts for an EKS add-on
	// if a value was changed from default
	ResolveConflicts ekstypes.ResolveConflicts `json:"resolveConflicts,omitempty"`
	// PodIdentityAssociations holds a list of associations to be configured for the addon
	// +optional
	PodIdentityAssociations *[]PodIdentityAssociation `json:"podIdentityAssociations,omitempty"`
	// UseDefaultPodIdentityAssociations uses the pod identity associations recommended by the EKS API.
	// Defaults to false.
	// +optional
	UseDefaultPodIdentityAssociations bool `json:"useDefaultPodIdentityAssociations,omitempty"`
	// ConfigurationValues defines the set of configuration properties for add-ons.
	// For now, all properties will be specified as a JSON string
	// and have to respect the schema from DescribeAddonConfiguration.
	// +optional
	ConfigurationValues string `json:"configurationValues,omitempty"`
	// Force overwrites an existing self-managed add-on with an EKS managed add-on.
	// Force is intended to be used when migrating an existing self-managed add-on to an EKS managed add-on.
	Force bool `json:"-"`
	// +optional
	Publishers []string `json:"publishers,omitempty"`
	// +optional
	Types []string `json:"types,omitempty"`
	// +optional
	Owners []string `json:"owners,omitempty"`
}

// AddonsConfig holds the addons config.
type AddonsConfig struct {
	// AutoApplyPodIdentityAssociations specifies whether to automatically apply pod identity associations
	// for supported addons that require IAM permissions.
	// +optional
	AutoApplyPodIdentityAssociations bool `json:"autoApplyPodIdentityAssociations,omitempty"`

	// DisableDefaultAddons enables or disables creation of default networking addons when the cluster
	// is created.
	// By default, all default addons are installed as EKS addons.
	// +optional
	DisableDefaultAddons bool `json:"disableDefaultAddons,omitempty"`
}

func (a Addon) CanonicalName() string {
	return strings.ToLower(a.Name)
}

func (a Addon) Validate() error {
	invalidAddonConfigErr := func(errorMsg string) error {
		return fmt.Errorf("invalid configuration for %q addon: %s", a.Name, errorMsg)
	}

	if a.Name == "" {
		return invalidAddonConfigErr("name is required")
	}

	if !json.Valid([]byte(a.ConfigurationValues)) {
		if err := a.convertConfigurationValuesToJSON(); err != nil {
			return invalidAddonConfigErr(fmt.Sprintf("configurationValues: %q is not valid, supported format(s) are: JSON and YAML", a.ConfigurationValues))
		}
	}

	if a.HasIRSASet() {
		if a.HasPodIDsSet() {
			return invalidAddonConfigErr("cannot set IRSA config (`addon.ServiceAccountRoleARN`, `addon.AttachPolicyARNs`, `addon.AttachPolicy`, `addon.WellKnownPolicies`) and pod identity associations at the same time")
		}
		if err := a.checkAtMostOnePolicyProviderIsSet(); err != nil {
			return invalidAddonConfigErr(err.Error())
		}
	}

	if a.HasPodIDsSet() {
		if a.CanonicalName() == PodIdentityAgentAddon {
			return invalidAddonConfigErr(fmt.Sprintf("cannot set pod identity associtations for %q addon", PodIdentityAgentAddon))
		}

		for i, pia := range *a.PodIdentityAssociations {
			path := fmt.Sprintf("podIdentityAssociations[%d]", i)
			if pia.Namespace == "" {
				return invalidAddonConfigErr(fmt.Sprintf("%s.namespace must be set", path))
			}
			if pia.ServiceAccountName == "" {
				return invalidAddonConfigErr(fmt.Sprintf("%s.serviceAccountName must be set", path))
			}

			if pia.RoleARN == "" &&
				len(pia.PermissionPolicy) == 0 &&
				len(pia.PermissionPolicyARNs) == 0 &&
				!pia.WellKnownPolicies.HasPolicy() {
				return invalidAddonConfigErr(fmt.Sprintf("at least one of the following must be specified: %[1]s.roleARN, %[1]s.permissionPolicy, %[1]s.permissionPolicyARNs, %[1]s.wellKnownPolicies", path))
			}

			if pia.RoleARN != "" {
				makeIncompatibleFieldErr := func(fieldName string) error {
					return invalidAddonConfigErr(fmt.Sprintf("%[1]s.%s cannot be specified when %[1]s.roleARN is set", path, fieldName))
				}
				if len(pia.PermissionPolicy) > 0 {
					return makeIncompatibleFieldErr("permissionPolicy")
				}
				if len(pia.PermissionPolicyARNs) > 0 {
					return makeIncompatibleFieldErr("permissionPolicyARNs")
				}
				if pia.WellKnownPolicies.HasPolicy() {
					return makeIncompatibleFieldErr("wellKnownPolicies")
				}
			}
		}
	}

	return nil
}

func (a *Addon) convertConfigurationValuesToJSON() (err error) {
	rawConfigurationValues := []byte(a.ConfigurationValues)
	var js map[string]interface{}
	if err = yaml.UnmarshalStrict(rawConfigurationValues, &js); err == nil {
		var JSONConfigurationValues []byte
		if JSONConfigurationValues, err = yaml.YAMLToJSONStrict(rawConfigurationValues); err == nil {
			a.ConfigurationValues = string(JSONConfigurationValues)
		}
	}
	return err
}

func (a Addon) checkAtMostOnePolicyProviderIsSet() error {
	setPolicyProviders := 0
	if a.AttachPolicy != nil {
		setPolicyProviders++
	}

	if a.AttachPolicyARNs != nil && len(a.AttachPolicyARNs) > 0 {
		setPolicyProviders++
	}

	if a.ServiceAccountRoleARN != "" {
		setPolicyProviders++
	}

	if a.WellKnownPolicies.HasPolicy() {
		setPolicyProviders++
	}

	if setPolicyProviders > 1 {
		return fmt.Errorf("at most one of wellKnownPolicies, serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified")
	}
	return nil
}

func (a Addon) HasIRSAPoliciesSet() bool {
	return len(a.AttachPolicyARNs) != 0 || a.WellKnownPolicies.HasPolicy() || a.AttachPolicy != nil

}

func (a Addon) HasIRSASet() bool {
	return a.ServiceAccountRoleARN != "" || a.HasIRSAPoliciesSet()
}

func (a Addon) HasPodIDsSet() bool {
	return a.PodIdentityAssociations != nil && len(*a.PodIdentityAssociations) > 0
}
