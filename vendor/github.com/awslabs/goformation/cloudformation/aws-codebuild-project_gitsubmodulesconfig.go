package cloudformation

import (
	"encoding/json"
)

// AWSCodeBuildProject_GitSubmodulesConfig AWS CloudFormation Resource (AWS::CodeBuild::Project.GitSubmodulesConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-gitsubmodulesconfig.html
type AWSCodeBuildProject_GitSubmodulesConfig struct {

	// FetchSubmodules AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-gitsubmodulesconfig.html#cfn-codebuild-project-gitsubmodulesconfig-fetchsubmodules
	FetchSubmodules *Value `json:"FetchSubmodules,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_GitSubmodulesConfig) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.GitSubmodulesConfig"
}

func (r *AWSCodeBuildProject_GitSubmodulesConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
