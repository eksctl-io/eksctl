package cloudfront

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// OriginRequestPolicy_QueryStringsConfig AWS CloudFormation Resource (AWS::CloudFront::OriginRequestPolicy.QueryStringsConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-querystringsconfig.html
type OriginRequestPolicy_QueryStringsConfig struct {

	// QueryStringBehavior AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-querystringsconfig.html#cfn-cloudfront-originrequestpolicy-querystringsconfig-querystringbehavior
	QueryStringBehavior *types.Value `json:"QueryStringBehavior,omitempty"`

	// QueryStrings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-querystringsconfig.html#cfn-cloudfront-originrequestpolicy-querystringsconfig-querystrings
	QueryStrings *types.Value `json:"QueryStrings,omitempty"`

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
func (r *OriginRequestPolicy_QueryStringsConfig) AWSCloudFormationType() string {
	return "AWS::CloudFront::OriginRequestPolicy.QueryStringsConfig"
}
