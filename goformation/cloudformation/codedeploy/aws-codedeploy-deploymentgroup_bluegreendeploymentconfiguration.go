package codedeploy

import (
	"goformation/v4/cloudformation/policies"
)

// DeploymentGroup_BlueGreenDeploymentConfiguration AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.BlueGreenDeploymentConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-bluegreendeploymentconfiguration.html
type DeploymentGroup_BlueGreenDeploymentConfiguration struct {

	// DeploymentReadyOption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-bluegreendeploymentconfiguration.html#cfn-codedeploy-deploymentgroup-bluegreendeploymentconfiguration-deploymentreadyoption
	DeploymentReadyOption *DeploymentGroup_DeploymentReadyOption `json:"DeploymentReadyOption,omitempty"`

	// GreenFleetProvisioningOption AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-bluegreendeploymentconfiguration.html#cfn-codedeploy-deploymentgroup-bluegreendeploymentconfiguration-greenfleetprovisioningoption
	GreenFleetProvisioningOption *DeploymentGroup_GreenFleetProvisioningOption `json:"GreenFleetProvisioningOption,omitempty"`

	// TerminateBlueInstancesOnDeploymentSuccess AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-bluegreendeploymentconfiguration.html#cfn-codedeploy-deploymentgroup-bluegreendeploymentconfiguration-terminateblueinstancesondeploymentsuccess
	TerminateBlueInstancesOnDeploymentSuccess *DeploymentGroup_BlueInstanceTerminationOption `json:"TerminateBlueInstancesOnDeploymentSuccess,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *DeploymentGroup_BlueGreenDeploymentConfiguration) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.BlueGreenDeploymentConfiguration"
}
