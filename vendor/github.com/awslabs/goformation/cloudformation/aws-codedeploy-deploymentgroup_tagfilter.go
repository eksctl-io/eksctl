package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_TagFilter AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.TagFilter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-tagfilter.html
type AWSCodeDeployDeploymentGroup_TagFilter struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-tagfilter.html#cfn-codedeploy-deploymentgroup-tagfilter-key
	Key *Value `json:"Key,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-tagfilter.html#cfn-codedeploy-deploymentgroup-tagfilter-type
	Type *Value `json:"Type,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-tagfilter.html#cfn-codedeploy-deploymentgroup-tagfilter-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_TagFilter) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.TagFilter"
}

func (r *AWSCodeDeployDeploymentGroup_TagFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
