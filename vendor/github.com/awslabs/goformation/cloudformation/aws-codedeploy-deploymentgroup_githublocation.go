package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_GitHubLocation AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.GitHubLocation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-githublocation.html
type AWSCodeDeployDeploymentGroup_GitHubLocation struct {

	// CommitId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-githublocation.html#cfn-properties-codedeploy-deploymentgroup-deployment-revision-githublocation-commitid
	CommitId *Value `json:"CommitId,omitempty"`

	// Repository AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deployment-revision-githublocation.html#cfn-properties-codedeploy-deploymentgroup-deployment-revision-githublocation-repository
	Repository *Value `json:"Repository,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_GitHubLocation) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.GitHubLocation"
}

func (r *AWSCodeDeployDeploymentGroup_GitHubLocation) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
