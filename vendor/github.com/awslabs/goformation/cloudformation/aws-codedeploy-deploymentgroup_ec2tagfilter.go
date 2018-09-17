package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_EC2TagFilter AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.EC2TagFilter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-ec2tagfilters.html
type AWSCodeDeployDeploymentGroup_EC2TagFilter struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-ec2tagfilters.html#cfn-properties-codedeploy-deploymentgroup-ec2tagfilters-key
	Key *Value `json:"Key,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-ec2tagfilters.html#cfn-properties-codedeploy-deploymentgroup-ec2tagfilters-type
	Type *Value `json:"Type,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-ec2tagfilters.html#cfn-properties-codedeploy-deploymentgroup-ec2tagfilters-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_EC2TagFilter) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.EC2TagFilter"
}

func (r *AWSCodeDeployDeploymentGroup_EC2TagFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
