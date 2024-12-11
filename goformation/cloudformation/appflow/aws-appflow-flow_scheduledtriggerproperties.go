package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Flow_ScheduledTriggerProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.ScheduledTriggerProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html
type Flow_ScheduledTriggerProperties struct {

	// DataPullMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-datapullmode
	DataPullMode *types.Value `json:"DataPullMode,omitempty"`

	// ScheduleEndTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-scheduleendtime
	ScheduleEndTime *types.Value `json:"ScheduleEndTime,omitempty"`

	// ScheduleExpression AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-scheduleexpression
	ScheduleExpression *types.Value `json:"ScheduleExpression,omitempty"`

	// ScheduleOffset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-scheduleoffset
	ScheduleOffset *types.Value `json:"ScheduleOffset,omitempty"`

	// ScheduleStartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-schedulestarttime
	ScheduleStartTime *types.Value `json:"ScheduleStartTime,omitempty"`

	// TimeZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-scheduledtriggerproperties.html#cfn-appflow-flow-scheduledtriggerproperties-timezone
	TimeZone *types.Value `json:"TimeZone,omitempty"`

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
func (r *Flow_ScheduledTriggerProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.ScheduledTriggerProperties"
}
