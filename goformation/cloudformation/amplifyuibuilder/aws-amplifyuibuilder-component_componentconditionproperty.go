package amplifyuibuilder

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Component_ComponentConditionProperty AWS CloudFormation Resource (AWS::AmplifyUIBuilder::Component.ComponentConditionProperty)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentconditionproperty.html
type Component_ComponentConditionProperty struct {

	// Else AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentconditionproperty.html#cfn-amplifyuibuilder-component-componentconditionproperty-else
	Else *Component_ComponentProperty `json:"Else,omitempty"`

	// Field AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentconditionproperty.html#cfn-amplifyuibuilder-component-componentconditionproperty-field
	Field *types.Value `json:"Field,omitempty"`

	// Operand AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentconditionproperty.html#cfn-amplifyuibuilder-component-componentconditionproperty-operand
	Operand *types.Value `json:"Operand,omitempty"`

	// Operator AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentconditionproperty.html#cfn-amplifyuibuilder-component-componentconditionproperty-operator
	Operator *types.Value `json:"Operator,omitempty"`

	// Property AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentconditionproperty.html#cfn-amplifyuibuilder-component-componentconditionproperty-property
	Property *types.Value `json:"Property,omitempty"`

	// Then AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentconditionproperty.html#cfn-amplifyuibuilder-component-componentconditionproperty-then
	Then *Component_ComponentProperty `json:"Then,omitempty"`

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
func (r *Component_ComponentConditionProperty) AWSCloudFormationType() string {
	return "AWS::AmplifyUIBuilder::Component.ComponentConditionProperty"
}
