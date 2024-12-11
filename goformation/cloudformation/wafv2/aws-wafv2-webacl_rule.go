package wafv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WebACL_Rule AWS CloudFormation Resource (AWS::WAFv2::WebACL.Rule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-rule.html
type WebACL_Rule struct {

	// Action AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-rule.html#cfn-wafv2-webacl-rule-action
	Action *WebACL_RuleAction `json:"Action,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-rule.html#cfn-wafv2-webacl-rule-name
	Name *types.Value `json:"Name,omitempty"`

	// OverrideAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-rule.html#cfn-wafv2-webacl-rule-overrideaction
	OverrideAction *WebACL_OverrideAction `json:"OverrideAction,omitempty"`

	// Priority AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-rule.html#cfn-wafv2-webacl-rule-priority
	Priority *types.Value `json:"Priority"`

	// RuleLabels AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-rule.html#cfn-wafv2-webacl-rule-rulelabels
	RuleLabels []WebACL_Label `json:"RuleLabels,omitempty"`

	// Statement AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-rule.html#cfn-wafv2-webacl-rule-statement
	Statement *WebACL_Statement `json:"Statement,omitempty"`

	// VisibilityConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-rule.html#cfn-wafv2-webacl-rule-visibilityconfig
	VisibilityConfig *WebACL_VisibilityConfig `json:"VisibilityConfig,omitempty"`

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
func (r *WebACL_Rule) AWSCloudFormationType() string {
	return "AWS::WAFv2::WebACL.Rule"
}
