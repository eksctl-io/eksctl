package builder

import (
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func makeClusterLogging(clusterConfig *api.ClusterConfig) *gfneks.Cluster_Logging {
	if !clusterConfig.HasClusterCloudWatchLogging() {
		return nil
	}

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
