package cloudformation

import (
	"encoding/json"
)

// AWSCodeBuildProject_LogsConfig AWS CloudFormation Resource (AWS::CodeBuild::Project.LogsConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-logsconfig.html
type AWSCodeBuildProject_LogsConfig struct {

	// CloudWatchLogs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-logsconfig.html#cfn-codebuild-project-logsconfig-cloudwatchlogs
	CloudWatchLogs *AWSCodeBuildProject_CloudWatchLogsConfig `json:"CloudWatchLogs,omitempty"`

	// S3Logs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-logsconfig.html#cfn-codebuild-project-logsconfig-s3logs
	S3Logs *AWSCodeBuildProject_S3LogsConfig `json:"S3Logs,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_LogsConfig) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.LogsConfig"
}

func (r *AWSCodeBuildProject_LogsConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
