package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_OutputFileUriValue AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.OutputFileUriValue)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable-outputfileurivalue.html
type AWSIoTAnalyticsDataset_OutputFileUriValue struct {

	// FileName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable-outputfileurivalue.html#cfn-iotanalytics-dataset-variable-outputfileurivalue-filename
	FileName *Value `json:"FileName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_OutputFileUriValue) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.OutputFileUriValue"
}

func (r *AWSIoTAnalyticsDataset_OutputFileUriValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
