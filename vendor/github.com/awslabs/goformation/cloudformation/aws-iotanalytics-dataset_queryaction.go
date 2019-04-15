package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_QueryAction AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.QueryAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-queryaction.html
type AWSIoTAnalyticsDataset_QueryAction struct {

	// Filters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-queryaction.html#cfn-iotanalytics-dataset-queryaction-filters
	Filters []AWSIoTAnalyticsDataset_Filter `json:"Filters,omitempty"`

	// SqlQuery AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-queryaction.html#cfn-iotanalytics-dataset-queryaction-sqlquery
	SqlQuery *Value `json:"SqlQuery,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_QueryAction) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.QueryAction"
}

func (r *AWSIoTAnalyticsDataset_QueryAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
