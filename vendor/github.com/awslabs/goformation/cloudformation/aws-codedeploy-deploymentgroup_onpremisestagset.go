package cloudformation

import (
	"encoding/json"
)

// AWSCodeDeployDeploymentGroup_OnPremisesTagSet AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.OnPremisesTagSet)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-onpremisestagset.html
type AWSCodeDeployDeploymentGroup_OnPremisesTagSet struct {

	// OnPremisesTagSetList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-onpremisestagset.html#cfn-codedeploy-deploymentgroup-onpremisestagset-onpremisestagsetlist
	OnPremisesTagSetList []AWSCodeDeployDeploymentGroup_OnPremisesTagSetListObject `json:"OnPremisesTagSetList,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_OnPremisesTagSet) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.OnPremisesTagSet"
}

func (r *AWSCodeDeployDeploymentGroup_OnPremisesTagSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
