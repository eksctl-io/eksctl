package cloudfront

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CachePolicy_ParametersInCacheKeyAndForwardedToOrigin AWS CloudFormation Resource (AWS::CloudFront::CachePolicy.ParametersInCacheKeyAndForwardedToOrigin)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin.html
type CachePolicy_ParametersInCacheKeyAndForwardedToOrigin struct {

	// CookiesConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin.html#cfn-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin-cookiesconfig
	CookiesConfig *CachePolicy_CookiesConfig `json:"CookiesConfig,omitempty"`

	// EnableAcceptEncodingBrotli AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin.html#cfn-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin-enableacceptencodingbrotli
	EnableAcceptEncodingBrotli *types.Value `json:"EnableAcceptEncodingBrotli,omitempty"`

	// EnableAcceptEncodingGzip AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin.html#cfn-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin-enableacceptencodinggzip
	EnableAcceptEncodingGzip *types.Value `json:"EnableAcceptEncodingGzip"`

	// HeadersConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin.html#cfn-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin-headersconfig
	HeadersConfig *CachePolicy_HeadersConfig `json:"HeadersConfig,omitempty"`

	// QueryStringsConfig AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin.html#cfn-cloudfront-cachepolicy-parametersincachekeyandforwardedtoorigin-querystringsconfig
	QueryStringsConfig *CachePolicy_QueryStringsConfig `json:"QueryStringsConfig,omitempty"`

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
func (r *CachePolicy_ParametersInCacheKeyAndForwardedToOrigin) AWSCloudFormationType() string {
	return "AWS::CloudFront::CachePolicy.ParametersInCacheKeyAndForwardedToOrigin"
}
