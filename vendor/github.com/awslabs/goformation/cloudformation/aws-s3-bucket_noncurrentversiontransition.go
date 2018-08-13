package cloudformation

import (
	"encoding/json"
)

// AWSS3Bucket_NoncurrentVersionTransition AWS CloudFormation Resource (AWS::S3::Bucket.NoncurrentVersionTransition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-lifecycleconfig-rule-noncurrentversiontransition.html
type AWSS3Bucket_NoncurrentVersionTransition struct {

	// StorageClass AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-lifecycleconfig-rule-noncurrentversiontransition.html#cfn-s3-bucket-lifecycleconfig-rule-noncurrentversiontransition-storageclass
	StorageClass *Value `json:"StorageClass,omitempty"`

	// TransitionInDays AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-lifecycleconfig-rule-noncurrentversiontransition.html#cfn-s3-bucket-lifecycleconfig-rule-noncurrentversiontransition-transitionindays
	TransitionInDays *Value `json:"TransitionInDays,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_NoncurrentVersionTransition) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.NoncurrentVersionTransition"
}

func (r *AWSS3Bucket_NoncurrentVersionTransition) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
