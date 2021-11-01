package eks

import (
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
)

const clusterInfoCacheTTL = 15 * time.Second

type ClusterInfo struct {
	timestamp time.Time
	Cluster   *awseks.Cluster
}

func (c *ClusterProviderImpl) clusterInfoNeedsUpdate() bool {
	if c.status.ClusterInfo == nil {
		return true
	}
	if time.Since(c.status.ClusterInfo.timestamp) > clusterInfoCacheTTL {
		return true
	}
	return false
}

func (c *ClusterProviderImpl) setClusterInfo(cluster *awseks.Cluster) {
	c.status.ClusterInfo = &ClusterInfo{
		timestamp: time.Now(),
		Cluster:   cluster,
	}
}
