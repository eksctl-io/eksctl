package events

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Rule_RetryPolicy AWS CloudFormation Resource (AWS::Events::Rule.RetryPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-retrypolicy.html
type Rule_RetryPolicy struct {

	// MaximumEventAgeInSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-retrypolicy.html#cfn-events-rule-retrypolicy-maximumeventageinseconds
	MaximumEventAgeInSeconds *types.Value `json:"MaximumEventAgeInSeconds,omitempty"`

	// MaximumRetryAttempts AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-retrypolicy.html#cfn-events-rule-retrypolicy-maximumretryattempts
	MaximumRetryAttempts *types.Value `json:"MaximumRetryAttempts,omitempty"`

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
func (r *Rule_RetryPolicy) AWSCloudFormationType() string {
	return "AWS::Events::Rule.RetryPolicy"
}
