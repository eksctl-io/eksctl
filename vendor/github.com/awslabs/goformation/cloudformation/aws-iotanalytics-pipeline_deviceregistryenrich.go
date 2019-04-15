package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_DeviceRegistryEnrich AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.DeviceRegistryEnrich)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceregistryenrich.html
type AWSIoTAnalyticsPipeline_DeviceRegistryEnrich struct {

	// Attribute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceregistryenrich.html#cfn-iotanalytics-pipeline-deviceregistryenrich-attribute
	Attribute *Value `json:"Attribute,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceregistryenrich.html#cfn-iotanalytics-pipeline-deviceregistryenrich-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceregistryenrich.html#cfn-iotanalytics-pipeline-deviceregistryenrich-next
	Next *Value `json:"Next,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceregistryenrich.html#cfn-iotanalytics-pipeline-deviceregistryenrich-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`

	// ThingName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-deviceregistryenrich.html#cfn-iotanalytics-pipeline-deviceregistryenrich-thingname
	ThingName *Value `json:"ThingName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_DeviceRegistryEnrich) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.DeviceRegistryEnrich"
}

func (r *AWSIoTAnalyticsPipeline_DeviceRegistryEnrich) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
