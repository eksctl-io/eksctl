package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_Filter AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.Filter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-filter.html
type AWSIoTAnalyticsDataset_Filter struct {

	// DeltaTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-filter.html#cfn-iotanalytics-dataset-filter-deltatime
	DeltaTime *AWSIoTAnalyticsDataset_DeltaTime `json:"DeltaTime,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_Filter) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.Filter"
}

func (r *AWSIoTAnalyticsDataset_Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
