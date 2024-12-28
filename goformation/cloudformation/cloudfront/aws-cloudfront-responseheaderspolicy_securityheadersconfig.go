package cloudfront

import (
	"goformation/v4/cloudformation/policies"
)

// ResponseHeadersPolicy_SecurityHeadersConfig AWS CloudFormation Resource (AWS::CloudFront::ResponseHeadersPolicy.SecurityHeadersConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-securityheadersconfig.html
type ResponseHeadersPolicy_SecurityHeadersConfig struct {

	// ContentSecurityPolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-securityheadersconfig.html#cfn-cloudfront-responseheaderspolicy-securityheadersconfig-contentsecuritypolicy
	ContentSecurityPolicy *ResponseHeadersPolicy_ContentSecurityPolicy `json:"ContentSecurityPolicy,omitempty"`

	// ContentTypeOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-securityheadersconfig.html#cfn-cloudfront-responseheaderspolicy-securityheadersconfig-contenttypeoptions
	ContentTypeOptions *ResponseHeadersPolicy_ContentTypeOptions `json:"ContentTypeOptions,omitempty"`

	// FrameOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-securityheadersconfig.html#cfn-cloudfront-responseheaderspolicy-securityheadersconfig-frameoptions
	FrameOptions *ResponseHeadersPolicy_FrameOptions `json:"FrameOptions,omitempty"`

	// ReferrerPolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-securityheadersconfig.html#cfn-cloudfront-responseheaderspolicy-securityheadersconfig-referrerpolicy
	ReferrerPolicy *ResponseHeadersPolicy_ReferrerPolicy `json:"ReferrerPolicy,omitempty"`

	// StrictTransportSecurity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-securityheadersconfig.html#cfn-cloudfront-responseheaderspolicy-securityheadersconfig-stricttransportsecurity
	StrictTransportSecurity *ResponseHeadersPolicy_StrictTransportSecurity `json:"StrictTransportSecurity,omitempty"`

	// XSSProtection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-securityheadersconfig.html#cfn-cloudfront-responseheaderspolicy-securityheadersconfig-xssprotection
	XSSProtection *ResponseHeadersPolicy_XSSProtection `json:"XSSProtection,omitempty"`

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
func (r *ResponseHeadersPolicy_SecurityHeadersConfig) AWSCloudFormationType() string {
	return "AWS::CloudFront::ResponseHeadersPolicy.SecurityHeadersConfig"
}
