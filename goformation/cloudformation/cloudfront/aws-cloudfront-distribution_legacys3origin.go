package cloudfront

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Distribution_LegacyS3Origin AWS CloudFormation Resource (AWS::CloudFront::Distribution.LegacyS3Origin)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-legacys3origin.html
type Distribution_LegacyS3Origin struct {

	// DNSName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-legacys3origin.html#cfn-cloudfront-distribution-legacys3origin-dnsname
	DNSName *types.Value `json:"DNSName,omitempty"`

	// OriginAccessIdentity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-legacys3origin.html#cfn-cloudfront-distribution-legacys3origin-originaccessidentity
	OriginAccessIdentity *types.Value `json:"OriginAccessIdentity,omitempty"`

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
func (r *Distribution_LegacyS3Origin) AWSCloudFormationType() string {
	return "AWS::CloudFront::Distribution.LegacyS3Origin"
}
