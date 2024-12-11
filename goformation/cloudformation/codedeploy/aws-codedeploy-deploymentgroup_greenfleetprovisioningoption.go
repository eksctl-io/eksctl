package codedeploy

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DeploymentGroup_GreenFleetProvisioningOption AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.GreenFleetProvisioningOption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-greenfleetprovisioningoption.html
type DeploymentGroup_GreenFleetProvisioningOption struct {

	// Action AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-greenfleetprovisioningoption.html#cfn-codedeploy-deploymentgroup-bluegreendeploymentconfiguration-greenfleetprovisioningoption-action
	Action *types.Value `json:"Action,omitempty"`

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
func (r *DeploymentGroup_GreenFleetProvisioningOption) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.GreenFleetProvisioningOption"
}
