package cloudformation

import (
	"encoding/json"
)

// AWSCodeBuildProject_WebhookFilter AWS CloudFormation Resource (AWS::CodeBuild::Project.WebhookFilter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-webhookfilter.html
type AWSCodeBuildProject_WebhookFilter struct {

	// ExcludeMatchedPattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-webhookfilter.html#cfn-codebuild-project-webhookfilter-excludematchedpattern
	ExcludeMatchedPattern *Value `json:"ExcludeMatchedPattern,omitempty"`

	// Pattern AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-webhookfilter.html#cfn-codebuild-project-webhookfilter-pattern
	Pattern *Value `json:"Pattern,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-webhookfilter.html#cfn-codebuild-project-webhookfilter-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_WebhookFilter) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.WebhookFilter"
}

func (r *AWSCodeBuildProject_WebhookFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
