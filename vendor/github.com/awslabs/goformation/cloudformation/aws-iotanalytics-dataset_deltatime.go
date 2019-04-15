package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_DeltaTime AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.DeltaTime)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-deltatime.html
type AWSIoTAnalyticsDataset_DeltaTime struct {

	// OffsetSeconds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-deltatime.html#cfn-iotanalytics-dataset-deltatime-offsetseconds
	OffsetSeconds *Value `json:"OffsetSeconds,omitempty"`

	// TimeExpression AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-deltatime.html#cfn-iotanalytics-dataset-deltatime-timeexpression
	TimeExpression *Value `json:"TimeExpression,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_DeltaTime) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.DeltaTime"
}

func (r *AWSIoTAnalyticsDataset_DeltaTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
