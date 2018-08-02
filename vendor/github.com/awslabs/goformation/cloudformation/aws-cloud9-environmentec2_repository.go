package cloudformation

// AWSCloud9EnvironmentEC2_Repository AWS CloudFormation Resource (AWS::Cloud9::EnvironmentEC2.Repository)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloud9-environmentec2-repository.html
type AWSCloud9EnvironmentEC2_Repository struct {

	// PathComponent AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloud9-environmentec2-repository.html#cfn-cloud9-environmentec2-repository-pathcomponent
	PathComponent *StringIntrinsic `json:"PathComponent,omitempty"`

	// RepositoryUrl AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloud9-environmentec2-repository.html#cfn-cloud9-environmentec2-repository-repositoryurl
	RepositoryUrl *StringIntrinsic `json:"RepositoryUrl,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloud9EnvironmentEC2_Repository) AWSCloudFormationType() string {
	return "AWS::Cloud9::EnvironmentEC2.Repository"
}
