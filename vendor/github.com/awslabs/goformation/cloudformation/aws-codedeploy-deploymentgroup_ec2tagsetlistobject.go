package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_EC2TagSetListObject AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.EC2TagSetListObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-ec2tagsetlistobject.html
type AWSCodeDeployDeploymentGroup_EC2TagSetListObject struct {

	// Ec2TagGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-ec2tagsetlistobject.html#cfn-codedeploy-deploymentgroup-ec2tagsetlistobject-ec2taggroup
	Ec2TagGroup []AWSCodeDeployDeploymentGroup_EC2TagFilter `json:"Ec2TagGroup,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_EC2TagSetListObject) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.EC2TagSetListObject"
}

func (r *AWSCodeDeployDeploymentGroup_EC2TagSetListObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
