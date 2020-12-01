package v1alpha5

import (
	"fmt"
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
	// AttachPolicy holds a policy document to attach to this service account
	// +optional
	AttachPolicy InlineDocument `json:"attachPolicy,omitempty"`
	// Force applies the add-on to overwrite an existing add-on
	Force bool
}

func (a Addon) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("name required")
	}

	if err := a.checkOnlyOnePolicyProviderIsSet(); err != nil {
		return err
	}

	return nil
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

	if setPolicyProviders > 1 {
		return fmt.Errorf("at most one of serviceAccountRoleARN, attachPolicyARNs and attachPolicy can be specified")
	}
	return nil
}
