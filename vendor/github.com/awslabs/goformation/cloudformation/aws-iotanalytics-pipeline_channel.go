package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_Channel AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.Channel)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-channel.html
type AWSIoTAnalyticsPipeline_Channel struct {

	// ChannelName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-channel.html#cfn-iotanalytics-pipeline-channel-channelname
	ChannelName *Value `json:"ChannelName,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-channel.html#cfn-iotanalytics-pipeline-channel-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-channel.html#cfn-iotanalytics-pipeline-channel-next
	Next *Value `json:"Next,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_Channel) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.Channel"
}

func (r *AWSIoTAnalyticsPipeline_Channel) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
