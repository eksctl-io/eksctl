package cloudformation

import (
	"encoding/json"
)

// AWSCodeBuildProject_FilterGroup AWS CloudFormation Resource (AWS::CodeBuild::Project.FilterGroup)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-filtergroup.html
type AWSCodeBuildProject_FilterGroup struct {
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_FilterGroup) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.FilterGroup"
}

func (r *AWSCodeBuildProject_FilterGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
