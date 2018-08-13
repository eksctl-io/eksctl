package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_SqsAction AWS CloudFormation Resource (AWS::IoT::TopicRule.SqsAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-sqsaction.html
type AWSIoTTopicRule_SqsAction struct {

	// QueueUrl AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-sqsaction.html#cfn-iot-topicrule-sqsaction-queueurl
	QueueUrl *Value `json:"QueueUrl,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-sqsaction.html#cfn-iot-topicrule-sqsaction-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// UseBase64 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-sqsaction.html#cfn-iot-topicrule-sqsaction-usebase64
	UseBase64 *Value `json:"UseBase64,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_SqsAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.SqsAction"
}

func (r *AWSIoTTopicRule_SqsAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
