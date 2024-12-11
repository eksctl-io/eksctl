package wafv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WebACL_JsonMatchPattern AWS CloudFormation Resource (AWS::WAFv2::WebACL.JsonMatchPattern)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-jsonmatchpattern.html
type WebACL_JsonMatchPattern struct {

	// All AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-jsonmatchpattern.html#cfn-wafv2-webacl-jsonmatchpattern-all
	All interface{} `json:"All,omitempty"`

	// IncludedPaths AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafv2-webacl-jsonmatchpattern.html#cfn-wafv2-webacl-jsonmatchpattern-includedpaths
	IncludedPaths *types.Value `json:"IncludedPaths,omitempty"`

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
func (r *WebACL_JsonMatchPattern) AWSCloudFormationType() string {
	return "AWS::WAFv2::WebACL.JsonMatchPattern"
}
