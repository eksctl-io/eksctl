package builder

import (
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

func makeClusterLogging(clusterConfig *api.ClusterConfig) *gfneks.Cluster_Logging {
	var enabledTypes []gfneks.Cluster_LoggingTypeConfig
	for _, t := range clusterConfig.CloudWatch.ClusterLogging.EnableTypes {
		enabledTypes = append(enabledTypes, gfneks.Cluster_LoggingTypeConfig{
			Type: gfnt.NewString(t),
		})
	}

	return &gfneks.Cluster_Logging{
		ClusterLogging: &gfneks.Cluster_ClusterLogging{
			EnabledTypes: enabledTypes,
		},
	}
}
