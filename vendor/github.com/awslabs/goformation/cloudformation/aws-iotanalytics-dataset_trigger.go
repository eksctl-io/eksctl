package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_Trigger AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.Trigger)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-trigger.html
type AWSIoTAnalyticsDataset_Trigger struct {

	// Schedule AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-trigger.html#cfn-iotanalytics-dataset-trigger-schedule
	Schedule *AWSIoTAnalyticsDataset_Schedule `json:"Schedule,omitempty"`

	// TriggeringDataset AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-trigger.html#cfn-iotanalytics-dataset-trigger-triggeringdataset
	TriggeringDataset *AWSIoTAnalyticsDataset_TriggeringDataset `json:"TriggeringDataset,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_Trigger) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.Trigger"
}

func (r *AWSIoTAnalyticsDataset_Trigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
