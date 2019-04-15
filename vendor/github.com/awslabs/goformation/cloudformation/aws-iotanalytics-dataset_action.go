package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_Action AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.Action)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-action.html
type AWSIoTAnalyticsDataset_Action struct {

	// ActionName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-action.html#cfn-iotanalytics-dataset-action-actionname
	ActionName *Value `json:"ActionName,omitempty"`

	// ContainerAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-action.html#cfn-iotanalytics-dataset-action-containeraction
	ContainerAction *AWSIoTAnalyticsDataset_ContainerAction `json:"ContainerAction,omitempty"`

	// QueryAction AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-action.html#cfn-iotanalytics-dataset-action-queryaction
	QueryAction *AWSIoTAnalyticsDataset_QueryAction `json:"QueryAction,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_Action) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.Action"
}

func (r *AWSIoTAnalyticsDataset_Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
