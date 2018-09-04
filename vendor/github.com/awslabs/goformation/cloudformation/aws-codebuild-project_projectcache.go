package cloudformation

// AWSCodeBuildProject_ProjectCache AWS CloudFormation Resource (AWS::CodeBuild::Project.ProjectCache)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projectcache.html
type AWSCodeBuildProject_ProjectCache struct {

	// Location AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projectcache.html#cfn-codebuild-project-projectcache-location
	Location *StringIntrinsic `json:"Location,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codebuild-project-projectcache.html#cfn-codebuild-project-projectcache-type
	Type *StringIntrinsic `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeBuildProject_ProjectCache) AWSCloudFormationType() string {
	return "AWS::CodeBuild::Project.ProjectCache"
}
