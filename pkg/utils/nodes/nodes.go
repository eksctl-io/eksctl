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

// CollectUniqueInstanceTypes returns a list of unique insance types collected from given NodePools
func CollectUniqueInstanceTypes(pools []api.NodePool) (instances []string) {
	alreadyCollected := make(map[string]struct{}, 0)
	for _, ng := range pools {
		for _, i := range ng.InstanceTypeList() {
			if _, ok := alreadyCollected[i]; !ok {
				instances = append(instances, i)
				alreadyCollected[i] = struct{}{}
			}
		}
	}
	return instances
}

// IsManaged returns whether a nodegroup is managed
func IsManaged(np api.NodePool) bool {
	_, managed := np.(*api.ManagedNodeGroup)
	return managed
}
