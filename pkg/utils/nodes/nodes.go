package nodes

import api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

// ToNodePools combines managed and self-managed nodegroups and returns a slice of api.NodePool
func ToNodePools(clusterConfig *api.ClusterConfig) []api.NodePool {
	var nodePools []api.NodePool
	for _, ng := range clusterConfig.NodeGroups {
		nodePools = append(nodePools, ng)
	}
	for _, ng := range clusterConfig.ManagedNodeGroups {
		nodePools = append(nodePools, ng)
	}
	return nodePools
}
