package cloudformation

// AWSCodeBuildProject_ProjectTriggers AWS CloudFormation Resource (AWS::CodeBuild::Project.ProjectTriggers)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projecttriggers.html
type AWSCodeBuildProject_ProjectTriggers struct {

	// Webhook AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projecttriggers.html#cfn-codebuild-project-projecttriggers-webhook
	Webhook bool `json:"Webhook,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_ProjectTriggers) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.ProjectTriggers"
}
