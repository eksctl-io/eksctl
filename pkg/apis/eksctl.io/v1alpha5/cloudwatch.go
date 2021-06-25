package v1alpha5

// ClusterCloudWatch contains config parameters related to CloudWatch
type ClusterCloudWatch struct {
	//+optional
	ClusterLogging *ClusterCloudWatchLogging `json:"clusterLogging,omitempty"`
}

// Values for `CloudWatchLogging`
const (
	apiLogging               = "api"
	auditLogging             = "audit"
	authenticatorLogging     = "authenticator"
	controllerManagerLogging = "controllerManager"
	schedulerLogging         = "scheduler"
	allLogging               = "all"
	wildcardLogging          = "*"
)

// ClusterCloudWatchLogging container config parameters related to cluster logging
type ClusterCloudWatchLogging struct {

	// Types of logging to enable (see [CloudWatch docs](/usage/cloudwatch-cluster-logging/#clusterconfig-examples)).
	// Valid entries are `CloudWatchLogging` constants
	//+optional
	EnableTypes []string `json:"enableTypes,omitempty"`
}

// SupportedCloudWatchClusterLogTypes retuls all supported logging facilities
func SupportedCloudWatchClusterLogTypes() []string {
	return []string{apiLogging, auditLogging, authenticatorLogging, controllerManagerLogging, schedulerLogging}
}

// HasClusterCloudWatchLogging determines if cluster logging was enabled or not
func (c *ClusterConfig) HasClusterCloudWatchLogging() bool {
	return c.CloudWatch != nil && c.CloudWatch.ClusterLogging != nil && len(c.CloudWatch.ClusterLogging.EnableTypes) > 0
}

func (c *ClusterConfig) ContainsWildcardCloudWatchLogging() bool {
	for _, v := range c.CloudWatch.ClusterLogging.EnableTypes {
		if v == allLogging || v == wildcardLogging {
			return true
		}
	}
	return false
}

// AppendClusterCloudWatchLogTypes will append given log types to the config structure
func (c *ClusterConfig) AppendClusterCloudWatchLogTypes(types ...string) {
	c.CloudWatch.ClusterLogging.EnableTypes = append(c.CloudWatch.ClusterLogging.EnableTypes, types...)
}
