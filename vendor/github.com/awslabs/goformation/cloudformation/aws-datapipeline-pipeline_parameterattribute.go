package cloudformation

import (
	"encoding/json"
)

// AWSDataPipelinePipeline_ParameterAttribute AWS CloudFormation Resource (AWS::DataPipeline::Pipeline.ParameterAttribute)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datapipeline-pipeline-parameterobjects-attributes.html
type AWSDataPipelinePipeline_ParameterAttribute struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datapipeline-pipeline-parameterobjects-attributes.html#cfn-datapipeline-pipeline-parameterobjects-attribtues-key
	Key *Value `json:"Key,omitempty"`

	// StringValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datapipeline-pipeline-parameterobjects-attributes.html#cfn-datapipeline-pipeline-parameterobjects-attribtues-stringvalue
	StringValue *Value `json:"StringValue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDataPipelinePipeline_ParameterAttribute) AWSCloudFormationType() string {
	return "AWS::DataPipeline::Pipeline.ParameterAttribute"
}

func (r *AWSDataPipelinePipeline_ParameterAttribute) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
