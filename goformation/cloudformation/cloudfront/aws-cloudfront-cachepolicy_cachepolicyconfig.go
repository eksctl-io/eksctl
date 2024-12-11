package cloudfront

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CachePolicy_CachePolicyConfig AWS CloudFormation Resource (AWS::CloudFront::CachePolicy.CachePolicyConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-cachepolicyconfig.html
type CachePolicy_CachePolicyConfig struct {

	// Comment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-cachepolicyconfig.html#cfn-cloudfront-cachepolicy-cachepolicyconfig-comment
	Comment *types.Value `json:"Comment,omitempty"`

	// DefaultTTL AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-cachepolicyconfig.html#cfn-cloudfront-cachepolicy-cachepolicyconfig-defaultttl
	DefaultTTL *types.Value `json:"DefaultTTL"`

	// MaxTTL AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-cachepolicyconfig.html#cfn-cloudfront-cachepolicy-cachepolicyconfig-maxttl
	MaxTTL *types.Value `json:"MaxTTL"`

	// MinTTL AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-cachepolicyconfig.html#cfn-cloudfront-cachepolicy-cachepolicyconfig-minttl
	MinTTL *types.Value `json:"MinTTL"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-cachepolicyconfig.html#cfn-cloudfront-cachepolicy-cachepolicyconfig-name
	Name *types.Value `json:"Name,omitempty"`

	// ParametersInCacheKeyAndForwardedToOrigin AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-cachepolicy-cachepolicyconfig.html#cfn-cloudfront-cachepolicy-cachepolicyconfig-parametersincachekeyandforwardedtoorigin
	ParametersInCacheKeyAndForwardedToOrigin *CachePolicy_ParametersInCacheKeyAndForwardedToOrigin `json:"ParametersInCacheKeyAndForwardedToOrigin,omitempty"`

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
func (r *CachePolicy_CachePolicyConfig) AWSCloudFormationType() string {
	return "AWS::CloudFront::CachePolicy.CachePolicyConfig"
}
