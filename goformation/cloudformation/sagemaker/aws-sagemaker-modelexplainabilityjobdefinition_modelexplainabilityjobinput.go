package sagemaker

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// ModelExplainabilityJobDefinition_ModelExplainabilityJobInput AWS CloudFormation Resource (AWS::SageMaker::ModelExplainabilityJobDefinition.ModelExplainabilityJobInput)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelexplainabilityjobdefinition-modelexplainabilityjobinput.html
type ModelExplainabilityJobDefinition_ModelExplainabilityJobInput struct {

	// EndpointInput AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-modelexplainabilityjobdefinition-modelexplainabilityjobinput.html#cfn-sagemaker-modelexplainabilityjobdefinition-modelexplainabilityjobinput-endpointinput
	EndpointInput *ModelExplainabilityJobDefinition_EndpointInput `json:"EndpointInput,omitempty"`

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
func (r *ModelExplainabilityJobDefinition_ModelExplainabilityJobInput) AWSCloudFormationType() string {
	return "AWS::SageMaker::ModelExplainabilityJobDefinition.ModelExplainabilityJobInput"
}
