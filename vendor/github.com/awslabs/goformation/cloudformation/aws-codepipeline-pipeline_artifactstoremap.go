package cloudformation

import (
	"encoding/json"
)

// AWSCodePipelinePipeline_ArtifactStoreMap AWS CloudFormation Resource (AWS::CodePipeline::Pipeline.ArtifactStoreMap)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-artifactstoremap.html
type AWSCodePipelinePipeline_ArtifactStoreMap struct {

	// ArtifactStore AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-artifactstoremap.html#cfn-codepipeline-pipeline-artifactstoremap-artifactstore
	ArtifactStore *AWSCodePipelinePipeline_ArtifactStore `json:"ArtifactStore,omitempty"`

	// Region AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codepipeline-pipeline-artifactstoremap.html#cfn-codepipeline-pipeline-artifactstoremap-region
	Region *Value `json:"Region,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodePipelinePipeline_ArtifactStoreMap) AWSCloudFormationType() string {
	return "AWS::CodePipeline::Pipeline.ArtifactStoreMap"
}

func (r *AWSCodePipelinePipeline_ArtifactStoreMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
