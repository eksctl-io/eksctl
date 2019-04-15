package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_DatasetContentVersionValue AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.DatasetContentVersionValue)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable-datasetcontentversionvalue.html
type AWSIoTAnalyticsDataset_DatasetContentVersionValue struct {

	// DatasetName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable-datasetcontentversionvalue.html#cfn-iotanalytics-dataset-variable-datasetcontentversionvalue-datasetname
	DatasetName *Value `json:"DatasetName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_DatasetContentVersionValue) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.DatasetContentVersionValue"
}

func (r *AWSIoTAnalyticsDataset_DatasetContentVersionValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
