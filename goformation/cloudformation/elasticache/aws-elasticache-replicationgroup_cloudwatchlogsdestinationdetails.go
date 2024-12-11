package elasticache

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ReplicationGroup_CloudWatchLogsDestinationDetails AWS CloudFormation Resource (AWS::ElastiCache::ReplicationGroup.CloudWatchLogsDestinationDetails)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-replicationgroup-cloudwatchlogsdestinationdetails.html
type ReplicationGroup_CloudWatchLogsDestinationDetails struct {

	// LogGroup AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-replicationgroup-cloudwatchlogsdestinationdetails.html#cfn-elasticache-replicationgroup-cloudwatchlogsdestinationdetails-loggroup
	LogGroup *types.Value `json:"LogGroup,omitempty"`

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
func (r *ReplicationGroup_CloudWatchLogsDestinationDetails) AWSCloudFormationType() string {
	return "AWS::ElastiCache::ReplicationGroup.CloudWatchLogsDestinationDetails"
}
