package builder

import (
	gfneks "goformation/v4/cloudformation/eks"
	gfnt "goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// AccessEntryResourceSet is a resource set for access entry.
type AccessEntryResourceSet struct {
	*resourceSet
	clusterName string
	accessEntry api.AccessEntry
}

// NewAccessEntryResourceSet creates and returns a new AccessEntryResourceSet.
func NewAccessEntryResourceSet(clusterName string, accessEntry api.AccessEntry) *AccessEntryResourceSet {
	return &AccessEntryResourceSet{
		resourceSet: newResourceSet(),
		clusterName: clusterName,
		accessEntry: accessEntry,
	}
}

// AddAllResources adds all resources required for creating an access entry.
func (a *AccessEntryResourceSet) AddAllResources() error {
	var accessPolicies []gfneks.AccessEntry_AccessPolicy
	for _, p := range a.accessEntry.AccessPolicies {
		var namespaces *gfnt.Value
		if len(p.AccessScope.Namespaces) > 0 {
			namespaces = gfnt.NewStringSlice(p.AccessScope.Namespaces...)
		}
		accessPolicies = append(accessPolicies, gfneks.AccessEntry_AccessPolicy{
			PolicyArn: gfnt.NewString(p.PolicyARN.String()),
			AccessScope: &gfneks.AccessEntry_AccessScope{
				Type:       gfnt.NewString(string(p.AccessScope.Type)),
				Namespaces: namespaces,
			},
		})
	}

	var entryType *gfnt.Value
	if a.accessEntry.Type != "" {
		entryType = gfnt.NewString(a.accessEntry.Type)
	}

	var kubernetesGroups *gfnt.Value
	if len(a.accessEntry.KubernetesGroups) > 0 {
		kubernetesGroups = gfnt.NewStringSlice(a.accessEntry.KubernetesGroups...)
	}
	var username *gfnt.Value
	if a.accessEntry.KubernetesUsername != "" {
		username = gfnt.NewString(a.accessEntry.KubernetesUsername)
	}
	a.newResource("AccessEntry", &gfneks.AccessEntry{
		PrincipalArn:     gfnt.NewString(a.accessEntry.PrincipalARN.String()),
		Type:             entryType,
		ClusterName:      gfnt.NewString(a.clusterName),
		KubernetesGroups: kubernetesGroups,
		Username:         username,
		AccessPolicies:   accessPolicies,
	})
	return nil
}

// RenderJSON implements the ResourceSet interface.
func (a *AccessEntryResourceSet) RenderJSON() ([]byte, error) {
	return a.renderJSON()
}

// WithIAM implements the ResourceSet interface.
func (*AccessEntryResourceSet) WithIAM() bool {
	return false
}

// WithNamedIAM implements the ResourceSet interface.
func (*AccessEntryResourceSet) WithNamedIAM() bool {
	return false
}
