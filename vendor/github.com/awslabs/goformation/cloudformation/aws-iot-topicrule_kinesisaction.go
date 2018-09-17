package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_KinesisAction AWS CloudFormation Resource (AWS::IoT::TopicRule.KinesisAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-kinesisaction.html
type AWSIoTTopicRule_KinesisAction struct {

	// PartitionKey AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-kinesisaction.html#cfn-iot-topicrule-kinesisaction-partitionkey
	PartitionKey *Value `json:"PartitionKey,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-kinesisaction.html#cfn-iot-topicrule-kinesisaction-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// StreamName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-kinesisaction.html#cfn-iot-topicrule-kinesisaction-streamname
	StreamName *Value `json:"StreamName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_KinesisAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.KinesisAction"
}

func (r *AWSIoTTopicRule_KinesisAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
