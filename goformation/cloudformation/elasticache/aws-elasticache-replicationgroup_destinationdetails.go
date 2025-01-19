package elasticache

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ReplicationGroup_DestinationDetails AWS CloudFormation Resource (AWS::ElastiCache::ReplicationGroup.DestinationDetails)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-replicationgroup-destinationdetails.html
type ReplicationGroup_DestinationDetails struct {

	// CloudWatchLogsDetails AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-replicationgroup-destinationdetails.html#cfn-elasticache-replicationgroup-destinationdetails-cloudwatchlogsdetails
	CloudWatchLogsDetails *ReplicationGroup_CloudWatchLogsDestinationDetails `json:"CloudWatchLogsDetails,omitempty"`

	// KinesisFirehoseDetails AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-replicationgroup-destinationdetails.html#cfn-elasticache-replicationgroup-destinationdetails-kinesisfirehosedetails
	KinesisFirehoseDetails *ReplicationGroup_KinesisFirehoseDestinationDetails `json:"KinesisFirehoseDetails,omitempty"`

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
func (r *ReplicationGroup_DestinationDetails) AWSCloudFormationType() string {
	return "AWS::ElastiCache::ReplicationGroup.DestinationDetails"
}
