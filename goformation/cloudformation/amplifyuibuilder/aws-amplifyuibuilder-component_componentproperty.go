package amplifyuibuilder

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Component_ComponentProperty AWS CloudFormation Resource (AWS::AmplifyUIBuilder::Component.ComponentProperty)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html
type Component_ComponentProperty struct {

	// BindingProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-bindingproperties
	BindingProperties *Component_ComponentPropertyBindingProperties `json:"BindingProperties,omitempty"`

	// Bindings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-bindings
	Bindings *Component_FormBindings `json:"Bindings,omitempty"`

	// CollectionBindingProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-collectionbindingproperties
	CollectionBindingProperties *Component_ComponentPropertyBindingProperties `json:"CollectionBindingProperties,omitempty"`

	// Concat AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-concat
	Concat []Component_ComponentProperty `json:"Concat,omitempty"`

	// Condition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-condition
	Condition *Component_ComponentConditionProperty `json:"Condition,omitempty"`

	// Configured AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-configured
	Configured *types.Value `json:"Configured,omitempty"`

	// DefaultValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-defaultvalue
	DefaultValue *types.Value `json:"DefaultValue,omitempty"`

	// Event AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-event
	Event *types.Value `json:"Event,omitempty"`

	// ImportedValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-importedvalue
	ImportedValue *types.Value `json:"ImportedValue,omitempty"`

	// Model AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-model
	Model *types.Value `json:"Model,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-type
	Type *types.Value `json:"Type,omitempty"`

	// UserAttribute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-userattribute
	UserAttribute *types.Value `json:"UserAttribute,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amplifyuibuilder-component-componentproperty.html#cfn-amplifyuibuilder-component-componentproperty-value
	Value *types.Value `json:"Value,omitempty"`

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
func (r *Component_ComponentProperty) AWSCloudFormationType() string {
	return "AWS::AmplifyUIBuilder::Component.ComponentProperty"
}
