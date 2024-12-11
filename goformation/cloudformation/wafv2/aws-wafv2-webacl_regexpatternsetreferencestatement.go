package wafv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WebACL_RegexPatternSetReferenceStatement AWS CloudFormation Resource (AWS::WAFv2::WebACL.RegexPatternSetReferenceStatement)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-regexpatternsetreferencestatement.html
type WebACL_RegexPatternSetReferenceStatement struct {

	// Arn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-regexpatternsetreferencestatement.html#cfn-wafv2-webacl-regexpatternsetreferencestatement-arn
	Arn *types.Value `json:"Arn,omitempty"`

	// FieldToMatch AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-regexpatternsetreferencestatement.html#cfn-wafv2-webacl-regexpatternsetreferencestatement-fieldtomatch
	FieldToMatch *WebACL_FieldToMatch `json:"FieldToMatch,omitempty"`

	// TextTransformations AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-regexpatternsetreferencestatement.html#cfn-wafv2-webacl-regexpatternsetreferencestatement-texttransformations
	TextTransformations []WebACL_TextTransformation `json:"TextTransformations,omitempty"`

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
func (r *WebACL_RegexPatternSetReferenceStatement) AWSCloudFormationType() string {
	return "AWS::WAFv2::WebACL.RegexPatternSetReferenceStatement"
}
