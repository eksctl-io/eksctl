package cloudformation

import (
	"encoding/json"
)

// AWSS3Bucket_VersioningConfiguration AWS CloudFormation Resource (AWS::S3::Bucket.VersioningConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-versioningconfig.html
type AWSS3Bucket_VersioningConfiguration struct {

	// Status AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-versioningconfig.html#cfn-s3-bucket-versioningconfig-status
	Status *Value `json:"Status,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_VersioningConfiguration) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.VersioningConfiguration"
}

func (r *AWSS3Bucket_VersioningConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
