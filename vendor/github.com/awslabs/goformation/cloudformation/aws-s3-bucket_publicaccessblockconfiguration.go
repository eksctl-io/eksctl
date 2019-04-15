package cloudformation

import (
	"encoding/json"
)

// AWSS3Bucket_PublicAccessBlockConfiguration AWS CloudFormation Resource (AWS::S3::Bucket.PublicAccessBlockConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-publicaccessblockconfiguration.html
type AWSS3Bucket_PublicAccessBlockConfiguration struct {

	// BlockPublicAcls AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-publicaccessblockconfiguration.html#cfn-s3-bucket-publicaccessblockconfiguration-blockpublicacls
	BlockPublicAcls *Value `json:"BlockPublicAcls,omitempty"`

	// BlockPublicPolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-publicaccessblockconfiguration.html#cfn-s3-bucket-publicaccessblockconfiguration-blockpublicpolicy
	BlockPublicPolicy *Value `json:"BlockPublicPolicy,omitempty"`

	// IgnorePublicAcls AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-publicaccessblockconfiguration.html#cfn-s3-bucket-publicaccessblockconfiguration-ignorepublicacls
	IgnorePublicAcls *Value `json:"IgnorePublicAcls,omitempty"`

	// RestrictPublicBuckets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-publicaccessblockconfiguration.html#cfn-s3-bucket-publicaccessblockconfiguration-restrictpublicbuckets
	RestrictPublicBuckets *Value `json:"RestrictPublicBuckets,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_PublicAccessBlockConfiguration) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.PublicAccessBlockConfiguration"
}

func (r *AWSS3Bucket_PublicAccessBlockConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
