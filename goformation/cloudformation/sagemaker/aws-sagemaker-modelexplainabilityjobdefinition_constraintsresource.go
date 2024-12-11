package sagemaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ModelExplainabilityJobDefinition_ConstraintsResource AWS CloudFormation Resource (AWS::SageMaker::ModelExplainabilityJobDefinition.ConstraintsResource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelexplainabilityjobdefinition-constraintsresource.html
type ModelExplainabilityJobDefinition_ConstraintsResource struct {

	// S3Uri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelexplainabilityjobdefinition-constraintsresource.html#cfn-sagemaker-modelexplainabilityjobdefinition-constraintsresource-s3uri
	S3Uri *types.Value `json:"S3Uri,omitempty"`

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
func (r *ModelExplainabilityJobDefinition_ConstraintsResource) AWSCloudFormationType() string {
	return "AWS::SageMaker::ModelExplainabilityJobDefinition.ConstraintsResource"
}
