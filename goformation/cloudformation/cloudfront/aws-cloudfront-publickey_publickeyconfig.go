package cloudfront

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PublicKey_PublicKeyConfig AWS CloudFormation Resource (AWS::CloudFront::PublicKey.PublicKeyConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-publickey-publickeyconfig.html
type PublicKey_PublicKeyConfig struct {

	// CallerReference AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-publickey-publickeyconfig.html#cfn-cloudfront-publickey-publickeyconfig-callerreference
	CallerReference *types.Value `json:"CallerReference,omitempty"`

	// Comment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-publickey-publickeyconfig.html#cfn-cloudfront-publickey-publickeyconfig-comment
	Comment *types.Value `json:"Comment,omitempty"`

	// EncodedKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-publickey-publickeyconfig.html#cfn-cloudfront-publickey-publickeyconfig-encodedkey
	EncodedKey *types.Value `json:"EncodedKey,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-publickey-publickeyconfig.html#cfn-cloudfront-publickey-publickeyconfig-name
	Name *types.Value `json:"Name,omitempty"`

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
func (r *PublicKey_PublicKeyConfig) AWSCloudFormationType() string {
	return "AWS::CloudFront::PublicKey.PublicKeyConfig"
}
