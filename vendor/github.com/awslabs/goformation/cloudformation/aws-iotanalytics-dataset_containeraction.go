package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_ContainerAction AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.ContainerAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-containeraction.html
type AWSIoTAnalyticsDataset_ContainerAction struct {

	// ExecutionRoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-containeraction.html#cfn-iotanalytics-dataset-containeraction-executionrolearn
	ExecutionRoleArn *Value `json:"ExecutionRoleArn,omitempty"`

	// Image AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-containeraction.html#cfn-iotanalytics-dataset-containeraction-image
	Image *Value `json:"Image,omitempty"`

	// ResourceConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-containeraction.html#cfn-iotanalytics-dataset-containeraction-resourceconfiguration
	ResourceConfiguration *AWSIoTAnalyticsDataset_ResourceConfiguration `json:"ResourceConfiguration,omitempty"`

	// Variables AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-containeraction.html#cfn-iotanalytics-dataset-containeraction-variables
	Variables []AWSIoTAnalyticsDataset_Variable `json:"Variables,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_ContainerAction) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.ContainerAction"
}

func (r *AWSIoTAnalyticsDataset_ContainerAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
