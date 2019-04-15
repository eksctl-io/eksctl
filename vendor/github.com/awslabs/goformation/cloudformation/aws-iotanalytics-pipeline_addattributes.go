package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_AddAttributes AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.AddAttributes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-addattributes.html
type AWSIoTAnalyticsPipeline_AddAttributes struct {

	// Attributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-addattributes.html#cfn-iotanalytics-pipeline-addattributes-attributes
	Attributes interface{} `json:"Attributes,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-addattributes.html#cfn-iotanalytics-pipeline-addattributes-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-addattributes.html#cfn-iotanalytics-pipeline-addattributes-next
	Next *Value `json:"Next,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_AddAttributes) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.AddAttributes"
}

func (r *AWSIoTAnalyticsPipeline_AddAttributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
