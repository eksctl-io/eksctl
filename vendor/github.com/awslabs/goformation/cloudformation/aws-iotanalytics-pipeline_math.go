package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_Math AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.Math)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-math.html
type AWSIoTAnalyticsPipeline_Math struct {

	// Attribute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-math.html#cfn-iotanalytics-pipeline-math-attribute
	Attribute *Value `json:"Attribute,omitempty"`

	// Math AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-math.html#cfn-iotanalytics-pipeline-math-math
	Math *Value `json:"Math,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-math.html#cfn-iotanalytics-pipeline-math-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-math.html#cfn-iotanalytics-pipeline-math-next
	Next *Value `json:"Next,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_Math) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.Math"
}

func (r *AWSIoTAnalyticsPipeline_Math) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
