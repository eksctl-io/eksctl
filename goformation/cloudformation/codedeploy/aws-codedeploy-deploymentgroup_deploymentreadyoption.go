package codedeploy

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DeploymentGroup_DeploymentReadyOption AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.DeploymentReadyOption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deploymentreadyoption.html
type DeploymentGroup_DeploymentReadyOption struct {

	// ActionOnTimeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deploymentreadyoption.html#cfn-codedeploy-deploymentgroup-bluegreendeploymentconfiguration-deploymentreadyoption-actionontimeout
	ActionOnTimeout *types.Value `json:"ActionOnTimeout,omitempty"`

	// WaitTimeInMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-deploymentreadyoption.html#cfn-codedeploy-deploymentgroup-bluegreendeploymentconfiguration-deploymentreadyoption-waittimeinminutes
	WaitTimeInMinutes *types.Value `json:"WaitTimeInMinutes,omitempty"`

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
func (r *DeploymentGroup_DeploymentReadyOption) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.DeploymentReadyOption"
}
