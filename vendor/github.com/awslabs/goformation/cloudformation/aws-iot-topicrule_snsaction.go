package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_SnsAction AWS CloudFormation Resource (AWS::IoT::TopicRule.SnsAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-snsaction.html
type AWSIoTTopicRule_SnsAction struct {

	// MessageFormat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-snsaction.html#cfn-iot-topicrule-snsaction-messageformat
	MessageFormat *Value `json:"MessageFormat,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-snsaction.html#cfn-iot-topicrule-snsaction-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// TargetArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-snsaction.html#cfn-iot-topicrule-snsaction-targetarn
	TargetArn *Value `json:"TargetArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_SnsAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.SnsAction"
}

func (r *AWSIoTTopicRule_SnsAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
