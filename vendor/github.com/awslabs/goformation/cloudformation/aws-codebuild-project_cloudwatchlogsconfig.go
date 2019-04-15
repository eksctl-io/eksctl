package cloudformation

import (
	"encoding/json"
)

// AWSCodeBuildProject_CloudWatchLogsConfig AWS CloudFormation Resource (AWS::CodeBuild::Project.CloudWatchLogsConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-cloudwatchlogsconfig.html
type AWSCodeBuildProject_CloudWatchLogsConfig struct {

	// GroupName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-cloudwatchlogsconfig.html#cfn-codebuild-project-cloudwatchlogsconfig-groupname
	GroupName *Value `json:"GroupName,omitempty"`

	// Status AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-cloudwatchlogsconfig.html#cfn-codebuild-project-cloudwatchlogsconfig-status
	Status *Value `json:"Status,omitempty"`

	// StreamName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-cloudwatchlogsconfig.html#cfn-codebuild-project-cloudwatchlogsconfig-streamname
	StreamName *Value `json:"StreamName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_CloudWatchLogsConfig) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.CloudWatchLogsConfig"
}

func (r *AWSCodeBuildProject_CloudWatchLogsConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
