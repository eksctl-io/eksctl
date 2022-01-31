package v1alpha5

import (
	"fmt"
	"strings"
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
	// Force applies the add-on to overwrite an existing add-on
	Force bool `json:"-"`
}

func (a Addon) CanonicalName() string {
	return strings.ToLower(a.Name)
}

func (a Addon) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name required")
	}

	return a.checkOnlyOnePolicyProviderIsSet()
}

func (a Addon) checkOnlyOnePolicyProviderIsSet() error {
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
