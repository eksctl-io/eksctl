package cloudformation

import (
	"encoding/json"
)

// AWSCodeBuildProject_S3LogsConfig AWS CloudFormation Resource (AWS::CodeBuild::Project.S3LogsConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-s3logsconfig.html
type AWSCodeBuildProject_S3LogsConfig struct {

	// EncryptionDisabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-s3logsconfig.html#cfn-codebuild-project-s3logsconfig-encryptiondisabled
	EncryptionDisabled *Value `json:"EncryptionDisabled,omitempty"`

	// Location AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-s3logsconfig.html#cfn-codebuild-project-s3logsconfig-location
	Location *Value `json:"Location,omitempty"`

	// Status AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-s3logsconfig.html#cfn-codebuild-project-s3logsconfig-status
	Status *Value `json:"Status,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_S3LogsConfig) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.S3LogsConfig"
}

func (r *AWSCodeBuildProject_S3LogsConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
