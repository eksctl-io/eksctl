package cloudformation

import (
	"encoding/json"
)

// AWSSESReceiptRule_S3Action AWS CloudFormation Resource (AWS::SES::ReceiptRule.S3Action)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-s3action.html
type AWSSESReceiptRule_S3Action struct {

	// BucketName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-s3action.html#cfn-ses-receiptrule-s3action-bucketname
	BucketName *Value `json:"BucketName,omitempty"`

	// KmsKeyArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-s3action.html#cfn-ses-receiptrule-s3action-kmskeyarn
	KmsKeyArn *Value `json:"KmsKeyArn,omitempty"`

	// ObjectKeyPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-s3action.html#cfn-ses-receiptrule-s3action-objectkeyprefix
	ObjectKeyPrefix *Value `json:"ObjectKeyPrefix,omitempty"`

	// TopicArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ses-receiptrule-s3action.html#cfn-ses-receiptrule-s3action-topicarn
	TopicArn *Value `json:"TopicArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESReceiptRule_S3Action) AWSCloudFormationType() string {
	return "AWS::SES::ReceiptRule.S3Action"
}

func (r *AWSSESReceiptRule_S3Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
