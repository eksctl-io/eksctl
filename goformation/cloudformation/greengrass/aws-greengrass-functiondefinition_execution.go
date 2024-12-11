package greengrass

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// FunctionDefinition_Execution AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition.Execution)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-execution.html
type FunctionDefinition_Execution struct {

	// IsolationMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-execution.html#cfn-greengrass-functiondefinition-execution-isolationmode
	IsolationMode *types.Value `json:"IsolationMode,omitempty"`

	// RunAs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-greengrass-functiondefinition-execution.html#cfn-greengrass-functiondefinition-execution-runas
	RunAs *FunctionDefinition_RunAs `json:"RunAs,omitempty"`

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
func (r *FunctionDefinition_Execution) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition.Execution"
}
