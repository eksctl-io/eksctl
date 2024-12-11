package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ModelBiasJobDefinition_VpcConfig AWS CloudFormation Resource (AWS::SageMaker::ModelBiasJobDefinition.VpcConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelbiasjobdefinition-vpcconfig.html
type ModelBiasJobDefinition_VpcConfig struct {

	// SecurityGroupIds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelbiasjobdefinition-vpcconfig.html#cfn-sagemaker-modelbiasjobdefinition-vpcconfig-securitygroupids
	SecurityGroupIds *types.Value `json:"SecurityGroupIds,omitempty"`

	// Subnets AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelbiasjobdefinition-vpcconfig.html#cfn-sagemaker-modelbiasjobdefinition-vpcconfig-subnets
	Subnets *types.Value `json:"Subnets,omitempty"`

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
func (r *ModelBiasJobDefinition_VpcConfig) AWSCloudFormationType() string {
	return "AWS::SageMaker::ModelBiasJobDefinition.VpcConfig"
}
