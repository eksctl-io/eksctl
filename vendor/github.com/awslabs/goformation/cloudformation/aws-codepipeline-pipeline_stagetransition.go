package cloudformation

import (
	"encoding/json"
)

// AWSCodePipelinePipeline_StageTransition AWS CloudFormation Resource (AWS::CodePipeline::Pipeline.StageTransition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-disableinboundstagetransitions.html
type AWSCodePipelinePipeline_StageTransition struct {

	// Reason AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-disableinboundstagetransitions.html#cfn-codepipeline-pipeline-disableinboundstagetransitions-reason
	Reason *Value `json:"Reason,omitempty"`

	// StageName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-disableinboundstagetransitions.html#cfn-codepipeline-pipeline-disableinboundstagetransitions-stagename
	StageName *Value `json:"StageName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodePipelinePipeline_StageTransition) AWSCloudFormationType() string {
	return "AWS::CodePipeline::Pipeline.StageTransition"
}

func (r *AWSCodePipelinePipeline_StageTransition) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
