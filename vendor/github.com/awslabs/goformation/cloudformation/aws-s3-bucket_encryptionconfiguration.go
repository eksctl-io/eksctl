package cloudformation

// AWSS3Bucket_EncryptionConfiguration AWS CloudFormation Resource (AWS::S3::Bucket.EncryptionConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-encryptionconfiguration.html
type AWSS3Bucket_EncryptionConfiguration struct {

	// ReplicaKmsKeyID AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-encryptionconfiguration.html#cfn-s3-bucket-encryptionconfiguration-replicakmskeyid
	ReplicaKmsKeyID *StringIntrinsic `json:"ReplicaKmsKeyID,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_EncryptionConfiguration) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.EncryptionConfiguration"
}
