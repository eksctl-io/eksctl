package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_ResourceConfiguration AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.ResourceConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-resourceconfiguration.html
type AWSIoTAnalyticsDataset_ResourceConfiguration struct {

	// ComputeType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-resourceconfiguration.html#cfn-iotanalytics-dataset-resourceconfiguration-computetype
	ComputeType *Value `json:"ComputeType,omitempty"`

	// VolumeSizeInGB AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-resourceconfiguration.html#cfn-iotanalytics-dataset-resourceconfiguration-volumesizeingb
	VolumeSizeInGB *Value `json:"VolumeSizeInGB,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_ResourceConfiguration) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.ResourceConfiguration"
}

func (r *AWSIoTAnalyticsDataset_ResourceConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
