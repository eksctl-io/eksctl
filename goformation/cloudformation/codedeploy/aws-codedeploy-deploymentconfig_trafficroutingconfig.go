package codedeploy

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DeploymentConfig_TrafficRoutingConfig AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentConfig.TrafficRoutingConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentconfig-trafficroutingconfig.html
type DeploymentConfig_TrafficRoutingConfig struct {

	// TimeBasedCanary AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentconfig-trafficroutingconfig.html#cfn-properties-codedeploy-deploymentconfig-trafficroutingconfig-timebasedcanary
	TimeBasedCanary *DeploymentConfig_TimeBasedCanary `json:"TimeBasedCanary,omitempty"`

	// TimeBasedLinear AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentconfig-trafficroutingconfig.html#cfn-properties-codedeploy-deploymentconfig-trafficroutingconfig-timebasedlinear
	TimeBasedLinear *DeploymentConfig_TimeBasedLinear `json:"TimeBasedLinear,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentconfig-trafficroutingconfig.html#cfn-properties-codedeploy-deploymentconfig-trafficroutingconfig-type
	Type *types.Value `json:"Type,omitempty"`

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
func (r *DeploymentConfig_TrafficRoutingConfig) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentConfig.TrafficRoutingConfig"
}
