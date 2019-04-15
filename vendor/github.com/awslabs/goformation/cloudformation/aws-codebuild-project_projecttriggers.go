package cloudformation

import (
	"encoding/json"
)

// AWSCodeBuildProject_ProjectTriggers AWS CloudFormation Resource (AWS::CodeBuild::Project.ProjectTriggers)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projecttriggers.html
type AWSCodeBuildProject_ProjectTriggers struct {

	// FilterGroups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projecttriggers.html#cfn-codebuild-project-projecttriggers-filtergroups
	FilterGroups []AWSCodeBuildProject_FilterGroup `json:"FilterGroups,omitempty"`

	// Webhook AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projecttriggers.html#cfn-codebuild-project-projecttriggers-webhook
	Webhook *Value `json:"Webhook,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_ProjectTriggers) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.ProjectTriggers"
}

func (r *AWSCodeBuildProject_ProjectTriggers) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
