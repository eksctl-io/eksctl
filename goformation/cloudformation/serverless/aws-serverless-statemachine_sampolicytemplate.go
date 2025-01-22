package serverless

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// StateMachine_SAMPolicyTemplate AWS CloudFormation Resource (AWS::Serverless::StateMachine.SAMPolicyTemplate)
// See: https://github.com/awslabs/serverless-application-model/blob/master/docs/policy_templates.rst
type StateMachine_SAMPolicyTemplate struct {

	// LambdaInvokePolicy AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/docs/policy_templates.rst
	LambdaInvokePolicy *StateMachine_FunctionSAMPT `json:"LambdaInvokePolicy,omitempty"`

	// StepFunctionsExecutionPolicy AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/docs/policy_templates.rst
	StepFunctionsExecutionPolicy *StateMachine_StateMachineSAMPT `json:"StepFunctionsExecutionPolicy,omitempty"`

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
func (r *StateMachine_SAMPolicyTemplate) AWSCloudFormationType() string {
	return "AWS::Serverless::StateMachine.SAMPolicyTemplate"
}
