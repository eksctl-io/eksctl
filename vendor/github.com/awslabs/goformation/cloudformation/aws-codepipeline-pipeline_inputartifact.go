package cloudformation

import (
	"encoding/json"
)

// AWSCodePipelinePipeline_InputArtifact AWS CloudFormation Resource (AWS::CodePipeline::Pipeline.InputArtifact)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-stages-actions-inputartifacts.html
type AWSCodePipelinePipeline_InputArtifact struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-stages-actions-inputartifacts.html#cfn-codepipeline-pipeline-stages-actions-inputartifacts-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodePipelinePipeline_InputArtifact) AWSCloudFormationType() string {
	return "AWS::CodePipeline::Pipeline.InputArtifact"
}

func (r *AWSCodePipelinePipeline_InputArtifact) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
