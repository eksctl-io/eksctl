package cloudformation

// AWSS3Bucket_SseKmsEncryptedObjects AWS CloudFormation Resource (AWS::S3::Bucket.SseKmsEncryptedObjects)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-ssekmsencryptedobjects.html
type AWSS3Bucket_SseKmsEncryptedObjects struct {

	// Status AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-ssekmsencryptedobjects.html#cfn-s3-bucket-ssekmsencryptedobjects-status
	Status *StringIntrinsic `json:"Status,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_SseKmsEncryptedObjects) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.SseKmsEncryptedObjects"
}
