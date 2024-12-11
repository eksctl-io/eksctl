package wafv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WebACL_IPSetForwardedIPConfiguration AWS CloudFormation Resource (AWS::WAFv2::WebACL.IPSetForwardedIPConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-ipsetforwardedipconfiguration.html
type WebACL_IPSetForwardedIPConfiguration struct {

	// FallbackBehavior AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-ipsetforwardedipconfiguration.html#cfn-wafv2-webacl-ipsetforwardedipconfiguration-fallbackbehavior
	FallbackBehavior *types.Value `json:"FallbackBehavior,omitempty"`

	// HeaderName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-ipsetforwardedipconfiguration.html#cfn-wafv2-webacl-ipsetforwardedipconfiguration-headername
	HeaderName *types.Value `json:"HeaderName,omitempty"`

	// Position AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-ipsetforwardedipconfiguration.html#cfn-wafv2-webacl-ipsetforwardedipconfiguration-position
	Position *types.Value `json:"Position,omitempty"`

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
func (r *WebACL_IPSetForwardedIPConfiguration) AWSCloudFormationType() string {
	return "AWS::WAFv2::WebACL.IPSetForwardedIPConfiguration"
}
