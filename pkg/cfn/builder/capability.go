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

	var deletePropagationPolicy *gfnt.Value
	if c.capability.DeletePropagationPolicy != "" {
		deletePropagationPolicy = gfnt.NewString(string(c.capability.DeletePropagationPolicy))
	}

	capability := &gfneks.Capability{
		CapabilityName:          gfnt.NewString(c.capability.Name),
		Type:                    gfnt.NewString(string(c.capability.Type)),
		ClusterName:             gfnt.NewString(c.clusterName),
		RoleArn:                 gfnt.NewString(c.capability.RoleARN),
		DeletePropagationPolicy: deletePropagationPolicy,
	}

	// Only set Configuration if it's not nil to avoid null values in CloudFormation
	if c.capability.Configuration != nil {

		config, err := convertConfiguration(c.capability.Configuration)
		if err != nil {
			return err
		}
		capability.Configuration = config
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

func convertConfiguration(config *api.CapabilityConfiguration) (*gfneks.CapabilityConfiguration, error) {

	req := &gfneks.CapabilityConfiguration{
		ArgoCd: gfneks.ArgoCDConfiguration{},
	}

	if config.ArgoCD.Namespace != "" {
		req.ArgoCd.Namespace = gfnt.NewString(config.ArgoCD.Namespace)
	}

	if config.ArgoCD.NetworkAccess != nil && len(config.ArgoCD.NetworkAccess.VPCEIDs) > 0 {
		req.ArgoCd.NetworkAccess = &gfneks.ArgoCDNetworkAccess{
			VPCEIDs: gfnt.NewStringSlice(config.ArgoCD.NetworkAccess.VPCEIDs...),
		}
	}

	if len(config.ArgoCD.RBACRoleMappings) > 0 {
		for _, mapping := range config.ArgoCD.RBACRoleMappings {
			var identities []gfneks.SSOIdentity
			for _, identity := range mapping.Identities {
				identities = append(identities, gfneks.SSOIdentity{
					ID:   gfnt.NewString(identity.ID),
					Type: gfnt.NewString(identity.Type),
				})
			}
			req.ArgoCd.RBACRoleMappings = append(req.ArgoCd.RBACRoleMappings, gfneks.ArgoCDRoleMapping{
				Role:       gfnt.NewString(mapping.Role),
				Identities: identities,
			})
		}
	}

	req.ArgoCd.AWSIDC = &gfneks.ArgoCDAWSIDC{
		IDCInstanceARN: gfnt.NewString(config.ArgoCD.AWSIDC.IDCInstanceARN),
	}
	if config.ArgoCD.AWSIDC.IDCRegion != "" {
		req.ArgoCd.AWSIDC.IDCRegion = gfnt.NewString(config.ArgoCD.AWSIDC.IDCRegion)
	}

	return req, nil
}
