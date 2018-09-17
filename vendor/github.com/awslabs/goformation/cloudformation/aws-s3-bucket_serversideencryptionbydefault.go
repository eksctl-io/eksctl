package cloudformation

import (
	"encoding/json"
)

// AWSS3Bucket_ServerSideEncryptionByDefault AWS CloudFormation Resource (AWS::S3::Bucket.ServerSideEncryptionByDefault)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-serversideencryptionbydefault.html
type AWSS3Bucket_ServerSideEncryptionByDefault struct {

	// KMSMasterKeyID AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-serversideencryptionbydefault.html#cfn-s3-bucket-serversideencryptionbydefault-kmsmasterkeyid
	KMSMasterKeyID *Value `json:"KMSMasterKeyID,omitempty"`

	// SSEAlgorithm AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-serversideencryptionbydefault.html#cfn-s3-bucket-serversideencryptionbydefault-ssealgorithm
	SSEAlgorithm *Value `json:"SSEAlgorithm,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_ServerSideEncryptionByDefault) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.ServerSideEncryptionByDefault"
}

func (r *AWSS3Bucket_ServerSideEncryptionByDefault) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
