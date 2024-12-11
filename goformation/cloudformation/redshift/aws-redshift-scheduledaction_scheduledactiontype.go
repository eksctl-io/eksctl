package redshift

import (
	"goformation/v4/cloudformation/policies"
)

// ScheduledAction_ScheduledActionType AWS CloudFormation Resource (AWS::Redshift::ScheduledAction.ScheduledActionType)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-scheduledactiontype.html
type ScheduledAction_ScheduledActionType struct {

	// PauseCluster AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-scheduledactiontype.html#cfn-redshift-scheduledaction-scheduledactiontype-pausecluster
	PauseCluster *ScheduledAction_PauseClusterMessage `json:"PauseCluster,omitempty"`

	// ResizeCluster AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-scheduledactiontype.html#cfn-redshift-scheduledaction-scheduledactiontype-resizecluster
	ResizeCluster *ScheduledAction_ResizeClusterMessage `json:"ResizeCluster,omitempty"`

	// ResumeCluster AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-scheduledactiontype.html#cfn-redshift-scheduledaction-scheduledactiontype-resumecluster
	ResumeCluster *ScheduledAction_ResumeClusterMessage `json:"ResumeCluster,omitempty"`

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
func (r *ScheduledAction_ScheduledActionType) AWSCloudFormationType() string {
	return "AWS::Redshift::ScheduledAction.ScheduledActionType"
}
