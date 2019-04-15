package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_OnPremisesTagSetListObject AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.OnPremisesTagSetListObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-onpremisestagsetlistobject.html
type AWSCodeDeployDeploymentGroup_OnPremisesTagSetListObject struct {

	// OnPremisesTagGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-onpremisestagsetlistobject.html#cfn-codedeploy-deploymentgroup-onpremisestagsetlistobject-onpremisestaggroup
	OnPremisesTagGroup []AWSCodeDeployDeploymentGroup_TagFilter `json:"OnPremisesTagGroup,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_OnPremisesTagSetListObject) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.OnPremisesTagSetListObject"
}

func (r *AWSCodeDeployDeploymentGroup_OnPremisesTagSetListObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
