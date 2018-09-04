package cloudformation

// AWSS3Bucket_AccessControlTranslation AWS CloudFormation Resource (AWS::S3::Bucket.AccessControlTranslation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-accesscontroltranslation.html
type AWSS3Bucket_AccessControlTranslation struct {

	// Owner AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket-accesscontroltranslation.html#cfn-s3-bucket-accesscontroltranslation-owner
	Owner *StringIntrinsic `json:"Owner,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSS3Bucket_AccessControlTranslation) AWSCloudFormationType() string {
	return "AWS::S3::Bucket.AccessControlTranslation"
}
