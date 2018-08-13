package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_S3Action AWS CloudFormation Resource (AWS::IoT::TopicRule.S3Action)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-s3action.html
type AWSIoTTopicRule_S3Action struct {

	// BucketName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-s3action.html#cfn-iot-topicrule-s3action-bucketname
	BucketName *Value `json:"BucketName,omitempty"`

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-s3action.html#cfn-iot-topicrule-s3action-key
	Key *Value `json:"Key,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-s3action.html#cfn-iot-topicrule-s3action-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_S3Action) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.S3Action"
}

func (r *AWSIoTTopicRule_S3Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
