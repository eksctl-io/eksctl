package builder

import (
	gfneks "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/eks"
	gfnt "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// CapabilityResourceSet is a resource set for capability.
type CapabilityResourceSet struct {
	*resourceSet
	clusterName string
	capability  api.Capability
}

// NewCapabilityResourceSet creates and returns a new CapabilityResourceSet.
func NewCapabilityResourceSet(clusterName string, capability api.Capability) *CapabilityResourceSet {
	return &CapabilityResourceSet{
		resourceSet: newResourceSet(),
		clusterName: clusterName,
		capability:  capability,
	}
}

// AddAllResources adds all resources required for creating a capability.
func (c *CapabilityResourceSet) AddAllResources() error {
	var roleArn *gfnt.Value
	// Only include RoleArn if IAM-related fields are specified
	if c.capability.RoleARN != "" || len(c.capability.AttachPolicyARNs) > 0 || c.capability.AttachPolicy != nil {
		if c.capability.RoleARN != "" {
			roleArn = gfnt.NewString(c.capability.RoleARN)
		}
		// If no explicit role but policies are specified, roleArn will be nil
		// which means CloudFormation/EKS will handle role creation
	}

	var deletePropagationPolicy *gfnt.Value
	if c.capability.DeletePropagationPolicy != "" {
		deletePropagationPolicy = gfnt.NewString(string(c.capability.DeletePropagationPolicy))
	}

	capability := &gfneks.Capability{
		CapabilityName:           gfnt.NewString(c.capability.Name),
		Type:                     gfnt.NewString(string(c.capability.Type)),
		ClusterName:              gfnt.NewString(c.clusterName),
		RoleArn:                  roleArn,
		DeletePropagationPolicy:  deletePropagationPolicy,
	}

	// Only set Configuration if it's not nil to avoid null values in CloudFormation
	if c.capability.Configuration != nil {
		capability.Configuration = c.capability.Configuration
	}

	c.newResource("Capability", capability)
	return nil
}

// RenderJSON implements the ResourceSet interface.
func (c *CapabilityResourceSet) RenderJSON() ([]byte, error) {
	return c.renderJSON()
}

// WithIAM implements the ResourceSet interface.
func (*CapabilityResourceSet) WithIAM() bool {
	return false
}

// WithNamedIAM implements the ResourceSet interface.
func (*CapabilityResourceSet) WithNamedIAM() bool {
	return false
}