package timestream

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ScheduledQuery_ScheduleConfiguration AWS CloudFormation Resource (AWS::Timestream::ScheduledQuery.ScheduleConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-scheduleconfiguration.html
type ScheduledQuery_ScheduleConfiguration struct {

	// ScheduleExpression AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-scheduleconfiguration.html#cfn-timestream-scheduledquery-scheduleconfiguration-scheduleexpression
	ScheduleExpression *types.Value `json:"ScheduleExpression,omitempty"`

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
func (r *ScheduledQuery_ScheduleConfiguration) AWSCloudFormationType() string {
	return "AWS::Timestream::ScheduledQuery.ScheduleConfiguration"
}
