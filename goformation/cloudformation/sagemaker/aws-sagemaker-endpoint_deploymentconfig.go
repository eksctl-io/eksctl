package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Endpoint_DeploymentConfig AWS CloudFormation Resource (AWS::SageMaker::Endpoint.DeploymentConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpoint-deploymentconfig.html
type Endpoint_DeploymentConfig struct {

	// AutoRollbackConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpoint-deploymentconfig.html#cfn-sagemaker-endpoint-deploymentconfig-autorollbackconfiguration
	AutoRollbackConfiguration *Endpoint_AutoRollbackConfig `json:"AutoRollbackConfiguration,omitempty"`

	// BlueGreenUpdatePolicy AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-endpoint-deploymentconfig.html#cfn-sagemaker-endpoint-deploymentconfig-bluegreenupdatepolicy
	BlueGreenUpdatePolicy *Endpoint_BlueGreenUpdatePolicy `json:"BlueGreenUpdatePolicy,omitempty"`

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
func (r *Endpoint_DeploymentConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::Endpoint.DeploymentConfig"
}
