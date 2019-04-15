package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_SelectAttributes AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.SelectAttributes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-selectattributes.html
type AWSIoTAnalyticsPipeline_SelectAttributes struct {

	// Attributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-selectattributes.html#cfn-iotanalytics-pipeline-selectattributes-attributes
	Attributes []*Value `json:"Attributes,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-selectattributes.html#cfn-iotanalytics-pipeline-selectattributes-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-selectattributes.html#cfn-iotanalytics-pipeline-selectattributes-next
	Next *Value `json:"Next,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_SelectAttributes) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.SelectAttributes"
}

func (r *AWSIoTAnalyticsPipeline_SelectAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
