package codestarnotifications

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// NotificationRule_Target AWS CloudFormation Resource (AWS::CodeStarNotifications::NotificationRule.Target)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codestarnotifications-notificationrule-target.html
type NotificationRule_Target struct {

	// TargetAddress AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codestarnotifications-notificationrule-target.html#cfn-codestarnotifications-notificationrule-target-targetaddress
	TargetAddress *types.Value `json:"TargetAddress,omitempty"`

	// TargetType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codestarnotifications-notificationrule-target.html#cfn-codestarnotifications-notificationrule-target-targettype
	TargetType *types.Value `json:"TargetType,omitempty"`

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
func (r *NotificationRule_Target) AWSCloudFormationType() string {
	return "AWS::CodeStarNotifications::NotificationRule.Target"
}
