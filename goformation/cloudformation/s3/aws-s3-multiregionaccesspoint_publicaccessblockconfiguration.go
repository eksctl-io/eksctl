package s3

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// MultiRegionAccessPoint_PublicAccessBlockConfiguration AWS CloudFormation Resource (AWS::S3::MultiRegionAccessPoint.PublicAccessBlockConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-multiregionaccesspoint-publicaccessblockconfiguration.html
type MultiRegionAccessPoint_PublicAccessBlockConfiguration struct {

	// BlockPublicAcls AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-multiregionaccesspoint-publicaccessblockconfiguration.html#cfn-s3-multiregionaccesspoint-publicaccessblockconfiguration-blockpublicacls
	BlockPublicAcls *types.Value `json:"BlockPublicAcls,omitempty"`

	// BlockPublicPolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-multiregionaccesspoint-publicaccessblockconfiguration.html#cfn-s3-multiregionaccesspoint-publicaccessblockconfiguration-blockpublicpolicy
	BlockPublicPolicy *types.Value `json:"BlockPublicPolicy,omitempty"`

	// IgnorePublicAcls AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-multiregionaccesspoint-publicaccessblockconfiguration.html#cfn-s3-multiregionaccesspoint-publicaccessblockconfiguration-ignorepublicacls
	IgnorePublicAcls *types.Value `json:"IgnorePublicAcls,omitempty"`

	// RestrictPublicBuckets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-multiregionaccesspoint-publicaccessblockconfiguration.html#cfn-s3-multiregionaccesspoint-publicaccessblockconfiguration-restrictpublicbuckets
	RestrictPublicBuckets *types.Value `json:"RestrictPublicBuckets,omitempty"`

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
func (r *MultiRegionAccessPoint_PublicAccessBlockConfiguration) AWSCloudFormationType() string {
	return "AWS::S3::MultiRegionAccessPoint.PublicAccessBlockConfiguration"
}
