package v1alpha5

import (
	"errors"
	"slices"
)

// validateBareCluster validates a cluster for unsupported fields if VPC CNI is disabled.
func validateBareCluster(clusterConfig *ClusterConfig) error {
	if !clusterConfig.AddonsConfig.DisableDefaultAddons || slices.ContainsFunc(clusterConfig.Addons, func(addon *Addon) bool {
		return addon.Name == VPCCNIAddon
	}) {
		return nil
	}
	if clusterConfig.HasNodes() || clusterConfig.IsFargateEnabled() || clusterConfig.Karpenter != nil || clusterConfig.HasGitOpsFluxConfigured() ||
		(clusterConfig.IAM != nil && (len(clusterConfig.IAM.ServiceAccounts) > 0) || len(clusterConfig.IAM.PodIdentityAssociations) > 0) {
		return errors.New("fields nodeGroups, managedNodeGroups, fargateProfiles, karpenter, gitops, iam.serviceAccounts, " +
			"and iam.podIdentityAssociations are not supported during cluster creation in a cluster without VPC CNI; please remove these fields " +
			"and add them back after cluster creation is successful")
	}
	return nil
}
