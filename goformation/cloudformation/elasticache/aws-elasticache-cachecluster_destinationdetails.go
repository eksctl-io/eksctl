package elasticache

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// CacheCluster_DestinationDetails AWS CloudFormation Resource (AWS::ElastiCache::CacheCluster.DestinationDetails)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-cachecluster-destinationdetails.html
type CacheCluster_DestinationDetails struct {

	// CloudWatchLogsDetails AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-cachecluster-destinationdetails.html#cfn-elasticache-cachecluster-destinationdetails-cloudwatchlogsdetails
	CloudWatchLogsDetails *CacheCluster_CloudWatchLogsDestinationDetails `json:"CloudWatchLogsDetails,omitempty"`

	// KinesisFirehoseDetails AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-cachecluster-destinationdetails.html#cfn-elasticache-cachecluster-destinationdetails-kinesisfirehosedetails
	KinesisFirehoseDetails *CacheCluster_KinesisFirehoseDestinationDetails `json:"KinesisFirehoseDetails,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *CacheCluster_DestinationDetails) AWSCloudFormationType() string {
	return "AWS::ElastiCache::CacheCluster.DestinationDetails"
}
