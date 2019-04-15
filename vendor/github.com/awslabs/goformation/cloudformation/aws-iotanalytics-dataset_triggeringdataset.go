package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_TriggeringDataset AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.TriggeringDataset)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-triggeringdataset.html
type AWSIoTAnalyticsDataset_TriggeringDataset struct {

	// DatasetName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-triggeringdataset.html#cfn-iotanalytics-dataset-triggeringdataset-datasetname
	DatasetName *Value `json:"DatasetName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_TriggeringDataset) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.TriggeringDataset"
}

func (r *AWSIoTAnalyticsDataset_TriggeringDataset) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
