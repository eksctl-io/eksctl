package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ModelBiasJobDefinition_ModelBiasAppSpecification AWS CloudFormation Resource (AWS::SageMaker::ModelBiasJobDefinition.ModelBiasAppSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelbiasjobdefinition-modelbiasappspecification.html
type ModelBiasJobDefinition_ModelBiasAppSpecification struct {

	// ConfigUri AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelbiasjobdefinition-modelbiasappspecification.html#cfn-sagemaker-modelbiasjobdefinition-modelbiasappspecification-configuri
	ConfigUri *types.Value `json:"ConfigUri,omitempty"`

	// Environment AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelbiasjobdefinition-modelbiasappspecification.html#cfn-sagemaker-modelbiasjobdefinition-modelbiasappspecification-environment
	Environment map[string]*types.Value `json:"Environment,omitempty"`

	// ImageUri AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelbiasjobdefinition-modelbiasappspecification.html#cfn-sagemaker-modelbiasjobdefinition-modelbiasappspecification-imageuri
	ImageUri *types.Value `json:"ImageUri,omitempty"`

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
func (r *ModelBiasJobDefinition_ModelBiasAppSpecification) AWSCloudFormationType() string {
	return "AWS::SageMaker::ModelBiasJobDefinition.ModelBiasAppSpecification"
}
