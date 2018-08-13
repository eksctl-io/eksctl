package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_Alarm AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.Alarm)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-alarm.html
type AWSCodeDeployDeploymentGroup_Alarm struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-alarm.html#cfn-codedeploy-deploymentgroup-alarm-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_Alarm) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.Alarm"
}

func (r *AWSCodeDeployDeploymentGroup_Alarm) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
