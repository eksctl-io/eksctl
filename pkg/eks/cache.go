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

func (c *ClusterProvider) clusterInfoNeedsUpdate() bool {
	if c.Status.ClusterInfo == nil {
		return true
	}
	if time.Since(c.Status.ClusterInfo.timestamp) > clusterInfoCacheTTL {
		return true
	}
	return false
}

func (c *ClusterProvider) setClusterInfo(cluster *awseks.Cluster) {
	c.Status.ClusterInfo = &ClusterInfo{
		timestamp: time.Now(),
		Cluster:   cluster,
	}
}
