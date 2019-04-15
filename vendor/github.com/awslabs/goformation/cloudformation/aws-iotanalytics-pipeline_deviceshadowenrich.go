package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_DeviceShadowEnrich AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.DeviceShadowEnrich)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceshadowenrich.html
type AWSIoTAnalyticsPipeline_DeviceShadowEnrich struct {

	// Attribute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceshadowenrich.html#cfn-iotanalytics-pipeline-deviceshadowenrich-attribute
	Attribute *Value `json:"Attribute,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceshadowenrich.html#cfn-iotanalytics-pipeline-deviceshadowenrich-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceshadowenrich.html#cfn-iotanalytics-pipeline-deviceshadowenrich-next
	Next *Value `json:"Next,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceshadowenrich.html#cfn-iotanalytics-pipeline-deviceshadowenrich-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// ThingName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceshadowenrich.html#cfn-iotanalytics-pipeline-deviceshadowenrich-thingname
	ThingName *Value `json:"ThingName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_DeviceShadowEnrich) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.DeviceShadowEnrich"
}

func (r *AWSIoTAnalyticsPipeline_DeviceShadowEnrich) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
