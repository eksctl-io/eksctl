package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_AutoRollbackConfiguration AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.AutoRollbackConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-autorollbackconfiguration.html
type AWSCodeDeployDeploymentGroup_AutoRollbackConfiguration struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-autorollbackconfiguration.html#cfn-codedeploy-deploymentgroup-autorollbackconfiguration-enabled
	Enabled *Value `json:"Enabled,omitempty"`

	// Events AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-autorollbackconfiguration.html#cfn-codedeploy-deploymentgroup-autorollbackconfiguration-events
	Events []*Value `json:"Events,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_AutoRollbackConfiguration) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.AutoRollbackConfiguration"
}

func (r *AWSCodeDeployDeploymentGroup_AutoRollbackConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
