package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_Schedule AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.Schedule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-trigger-schedule.html
type AWSIoTAnalyticsDataset_Schedule struct {

	// ScheduleExpression AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-trigger-schedule.html#cfn-iotanalytics-dataset-trigger-schedule-scheduleexpression
	ScheduleExpression *Value `json:"ScheduleExpression,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_Schedule) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.Schedule"
}

func (r *AWSIoTAnalyticsDataset_Schedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
