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

var LogRetentionInDaysValues = []int{1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, 3653}

// ClusterCloudWatchLogging container config parameters related to cluster logging
type ClusterCloudWatchLogging struct {

	// Types of logging to enable (see [CloudWatch docs](/usage/cloudwatch-cluster-logging/#clusterconfig-examples)).
	// Valid entries are `CloudWatchLogging` constants
	//+optional
	EnableTypes []string `json:"enableTypes,omitempty"`
	// LogRetentionInDays sets the number of days to retain the logs for (see [CloudWatch docs](https://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_PutRetentionPolicy.html#API_PutRetentionPolicy_RequestSyntax)) .
	// Valid values are: 1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731,
	// 1827, and 3653.
	//+optional
	LogRetentionInDays int `json:"logRetentionInDays,omitempty"`
}

// SupportedCloudWatchClusterLogTypes returns all supported logging facilities
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
