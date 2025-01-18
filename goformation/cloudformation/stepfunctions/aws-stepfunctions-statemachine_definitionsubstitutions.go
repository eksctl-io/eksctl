package stepfunctions

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// StateMachine_DefinitionSubstitutions AWS CloudFormation Resource (AWS::StepFunctions::StateMachine.DefinitionSubstitutions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-stepfunctions-statemachine-definitionsubstitutions.html
type StateMachine_DefinitionSubstitutions struct {

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
func (r *StateMachine_DefinitionSubstitutions) AWSCloudFormationType() string {
	return "AWS::StepFunctions::StateMachine.DefinitionSubstitutions"
}
