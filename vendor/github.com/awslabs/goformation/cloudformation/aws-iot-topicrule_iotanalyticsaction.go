package cloudformation

import (
	"encoding/json"
)

// AWSIoTTopicRule_IotAnalyticsAction AWS CloudFormation Resource (AWS::IoT::TopicRule.IotAnalyticsAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-iotanalyticsaction.html
type AWSIoTTopicRule_IotAnalyticsAction struct {

	// ChannelName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-iotanalyticsaction.html#cfn-iot-topicrule-iotanalyticsaction-channelname
	ChannelName *Value `json:"ChannelName,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-topicrule-iotanalyticsaction.html#cfn-iot-topicrule-iotanalyticsaction-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTTopicRule_IotAnalyticsAction) AWSCloudFormationType() string {
	return "AWS::IoT::TopicRule.IotAnalyticsAction"
}

func (r *AWSIoTTopicRule_IotAnalyticsAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
