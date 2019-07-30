package v1alpha5

// ClusterCloudWatch contains config parameters related to CloudWatch
type ClusterCloudWatch struct {
	//+optional
	ClusterLogging *ClusterCloudWatchLogging `json:"clusterLogging,omitempty"`
}

// ClusterCloudWatchLogging container config parameters related to cluster logging
type ClusterCloudWatchLogging struct {
	//+optional
	EnableTypes []string `json:"enableTypes,omitempty"`
}

// SupportedCloudWatchClusterLogTypes retuls all supported logging facilities
func SupportedCloudWatchClusterLogTypes() []string {
	return []string{"api", "audit", "authenticator", "controllerManager", "scheduler"}
}

// HasClusterCloudWatchLogging determines if cluster logging was enabled or not
func (c *ClusterConfig) HasClusterCloudWatchLogging() bool {
	if c.CloudWatch == nil {
		return false
	}
	if c.CloudWatch != nil {
		if c.CloudWatch.ClusterLogging == nil {
			return false
		}
		if len(c.CloudWatch.ClusterLogging.EnableTypes) == 0 {
			return false
		}
	}
	return true
}

// AppendClusterCloudWatchLogTypes will append given log types to the config structure
func (c *ClusterConfig) AppendClusterCloudWatchLogTypes(types ...string) {
	c.CloudWatch.ClusterLogging.EnableTypes = append(c.CloudWatch.ClusterLogging.EnableTypes, types...)
}
