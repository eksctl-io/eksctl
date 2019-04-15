package cloudformation

import (
	"encoding/json"
)

// AWSIoTAnalyticsDataset_Variable AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.Variable)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html
type AWSIoTAnalyticsDataset_Variable struct {

	// DatasetContentVersionValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-datasetcontentversionvalue
	DatasetContentVersionValue *AWSIoTAnalyticsDataset_DatasetContentVersionValue `json:"DatasetContentVersionValue,omitempty"`

	// DoubleValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-doublevalue
	DoubleValue *Value `json:"DoubleValue,omitempty"`

	// OutputFileUriValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-outputfileurivalue
	OutputFileUriValue *AWSIoTAnalyticsDataset_OutputFileUriValue `json:"OutputFileUriValue,omitempty"`

	// StringValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-stringvalue
	StringValue *Value `json:"StringValue,omitempty"`

	// VariableName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-variablename
	VariableName *Value `json:"VariableName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset_Variable) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.Variable"
}

func (r *AWSIoTAnalyticsDataset_Variable) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
