package iotanalytics

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Dataset_Variable AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset.Variable)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html
type Dataset_Variable struct {

	// DatasetContentVersionValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-datasetcontentversionvalue
	DatasetContentVersionValue *Dataset_DatasetContentVersionValue `json:"DatasetContentVersionValue,omitempty"`

	// DoubleValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-doublevalue
	DoubleValue *types.Value `json:"DoubleValue,omitempty"`

	// OutputFileUriValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-outputfileurivalue
	OutputFileUriValue *Dataset_OutputFileUriValue `json:"OutputFileUriValue,omitempty"`

	// StringValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-stringvalue
	StringValue *types.Value `json:"StringValue,omitempty"`

	// VariableName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotanalytics-dataset-variable.html#cfn-iotanalytics-dataset-variable-variablename
	VariableName *types.Value `json:"VariableName,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Dataset_Variable) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset.Variable"
}
