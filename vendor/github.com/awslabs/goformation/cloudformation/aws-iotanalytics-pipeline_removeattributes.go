package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_RemoveAttributes AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.RemoveAttributes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-removeattributes.html
type AWSIoTAnalyticsPipeline_RemoveAttributes struct {

	// Attributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-removeattributes.html#cfn-iotanalytics-pipeline-removeattributes-attributes
	Attributes []*Value `json:"Attributes,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-removeattributes.html#cfn-iotanalytics-pipeline-removeattributes-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-removeattributes.html#cfn-iotanalytics-pipeline-removeattributes-next
	Next *Value `json:"Next,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_RemoveAttributes) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.RemoveAttributes"
}

func (r *AWSIoTAnalyticsPipeline_RemoveAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
