package cloudfront

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ResponseHeadersPolicy_StrictTransportSecurity AWS CloudFormation Resource (AWS::CloudFront::ResponseHeadersPolicy.StrictTransportSecurity)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-stricttransportsecurity.html
type ResponseHeadersPolicy_StrictTransportSecurity struct {

	// AccessControlMaxAgeSec AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-stricttransportsecurity.html#cfn-cloudfront-responseheaderspolicy-stricttransportsecurity-accesscontrolmaxagesec
	AccessControlMaxAgeSec *types.Value `json:"AccessControlMaxAgeSec"`

	// IncludeSubdomains AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-stricttransportsecurity.html#cfn-cloudfront-responseheaderspolicy-stricttransportsecurity-includesubdomains
	IncludeSubdomains *types.Value `json:"IncludeSubdomains,omitempty"`

	// Override AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-stricttransportsecurity.html#cfn-cloudfront-responseheaderspolicy-stricttransportsecurity-override
	Override *types.Value `json:"Override"`

	// Preload AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-responseheaderspolicy-stricttransportsecurity.html#cfn-cloudfront-responseheaderspolicy-stricttransportsecurity-preload
	Preload *types.Value `json:"Preload,omitempty"`

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
func (r *ResponseHeadersPolicy_StrictTransportSecurity) AWSCloudFormationType() string {
	return "AWS::CloudFront::ResponseHeadersPolicy.StrictTransportSecurity"
}
