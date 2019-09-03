package eks

import (
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
)

const clusterInfoCacheTTL = 15 * time.Second

type clusterInfo struct {
	timestamp time.Time
	cluster   *awseks.Cluster
}

func (c *ClusterProvider) clusterInfoNeedsUpdate() bool {
	if c.Status.clusterInfo == nil {
		return true
	}
	if time.Since(c.Status.clusterInfo.timestamp) > clusterInfoCacheTTL {
		return true
	}
	return false
}

func (c *ClusterProvider) setClusterInfo(cluster *awseks.Cluster) {
	c.Status.clusterInfo = &clusterInfo{
		timestamp: time.Now(),
		cluster:   cluster,
	}
}
