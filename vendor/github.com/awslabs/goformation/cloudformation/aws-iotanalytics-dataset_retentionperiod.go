package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_RetentionPeriod AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.RetentionPeriod)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-retentionperiod.html
type AWSIoTAnalyticsDataset_RetentionPeriod struct {

	// NumberOfDays AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-retentionperiod.html#cfn-iotanalytics-dataset-retentionperiod-numberofdays
	NumberOfDays *Value `json:"NumberOfDays,omitempty"`

	// Unlimited AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-retentionperiod.html#cfn-iotanalytics-dataset-retentionperiod-unlimited
	Unlimited *Value `json:"Unlimited,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_RetentionPeriod) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.RetentionPeriod"
}

func (r *AWSIoTAnalyticsDataset_RetentionPeriod) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
