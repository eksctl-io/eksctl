package cmdutils

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type FilterLogFunc func()

func ApplyFilter(clusterConfig *api.ClusterConfig, ngFilter *NodeGroupFilter) (FilterLogFunc, FilterLogFunc) {
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

	return makeLogFunc(ngFilter, filteredNodeGroups, filteredManagedNodeGroups)
}

func makeLogFunc(ngFilter *NodeGroupFilter, nodeGroups []*api.NodeGroup, managedNodeGroups []*api.ManagedNodeGroup) (FilterLogFunc, FilterLogFunc) {

	ngLog := func() {
		var names []string
		for _, ng := range nodeGroups {
			names = append(names, ng.NameString())
		}
		ngFilter.doLogInfo("nodegroups", names)

	}
	mNgLog := func() {
		var names []string
		for _, ng := range managedNodeGroups {
			names = append(names, ng.NameString())
		}
		ngFilter.doLogInfo("managed nodegroups", names)
	}

	return ngLog, mNgLog
}
