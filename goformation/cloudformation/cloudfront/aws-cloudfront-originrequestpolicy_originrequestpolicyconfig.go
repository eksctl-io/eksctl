package cloudfront

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// OriginRequestPolicy_OriginRequestPolicyConfig AWS CloudFormation Resource (AWS::CloudFront::OriginRequestPolicy.OriginRequestPolicyConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-originrequestpolicyconfig.html
type OriginRequestPolicy_OriginRequestPolicyConfig struct {

	// Comment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-originrequestpolicyconfig.html#cfn-cloudfront-originrequestpolicy-originrequestpolicyconfig-comment
	Comment *types.Value `json:"Comment,omitempty"`

	// CookiesConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-originrequestpolicyconfig.html#cfn-cloudfront-originrequestpolicy-originrequestpolicyconfig-cookiesconfig
	CookiesConfig *OriginRequestPolicy_CookiesConfig `json:"CookiesConfig,omitempty"`

	// HeadersConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-originrequestpolicyconfig.html#cfn-cloudfront-originrequestpolicy-originrequestpolicyconfig-headersconfig
	HeadersConfig *OriginRequestPolicy_HeadersConfig `json:"HeadersConfig,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-originrequestpolicyconfig.html#cfn-cloudfront-originrequestpolicy-originrequestpolicyconfig-name
	Name *types.Value `json:"Name,omitempty"`

	// QueryStringsConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-originrequestpolicy-originrequestpolicyconfig.html#cfn-cloudfront-originrequestpolicy-originrequestpolicyconfig-querystringsconfig
	QueryStringsConfig *OriginRequestPolicy_QueryStringsConfig `json:"QueryStringsConfig,omitempty"`

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
func (r *OriginRequestPolicy_OriginRequestPolicyConfig) AWSCloudFormationType() string {
	return "AWS::CloudFront::OriginRequestPolicy.OriginRequestPolicyConfig"
}
