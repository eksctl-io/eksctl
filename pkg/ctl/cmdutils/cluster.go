package cmdutils

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
	"github.com/weaveworks/eksctl/pkg/eks"
)

// ApplyFilter applies nodegroup filters and returns a log function
func ApplyFilter(clusterConfig *api.ClusterConfig, ngFilter filter.NodegroupFilter) func() {
	var (
		filteredNodeGroups        []*api.NodeGroup
		filteredManagedNodeGroups []*api.ManagedNodeGroup
	)

	for _, ng := range clusterConfig.NodeGroups {
		if ngFilter.Match(ng.NameString()) {
			filteredNodeGroups = append(filteredNodeGroups, ng)
		}
	}

	for _, ng := range clusterConfig.ManagedNodeGroups {
		if ngFilter.Match(ng.NameString()) {
			filteredManagedNodeGroups = append(filteredManagedNodeGroups, ng)
		}
	}

	clusterConfig.NodeGroups, clusterConfig.ManagedNodeGroups = filteredNodeGroups, filteredManagedNodeGroups

	return func() {
		ngFilter.LogInfo(clusterConfig)
	}
}

// ToKubeNodeGroups combines managed and unmanaged nodegroups and returns a slice of eks.KubeNodeGroup containing
// both types of nodegroups
func ToKubeNodeGroups(unmanagedNodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup) []eks.KubeNodeGroup {
	var kubeNodeGroups []eks.KubeNodeGroup
	for _, ng := range unmanagedNodeGroups {
		kubeNodeGroups = append(kubeNodeGroups, ng)
	}
	for _, ng := range managedNodeGroups {
		kubeNodeGroups = append(kubeNodeGroups, ng)
	}
	return kubeNodeGroups
}
