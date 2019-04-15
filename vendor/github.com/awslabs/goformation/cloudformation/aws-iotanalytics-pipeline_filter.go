package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsPipeline_Filter AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline.Filter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-filter.html
type AWSIoTAnalyticsPipeline_Filter struct {

	// Filter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-filter.html#cfn-iotanalytics-pipeline-filter-filter
	Filter *Value `json:"Filter,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-filter.html#cfn-iotanalytics-pipeline-filter-name
	Name *Value `json:"Name,omitempty"`

	// Next AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-pipeline-filter.html#cfn-iotanalytics-pipeline-filter-next
	Next *Value `json:"Next,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline_Filter) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline.Filter"
}

func (r *AWSIoTAnalyticsPipeline_Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
