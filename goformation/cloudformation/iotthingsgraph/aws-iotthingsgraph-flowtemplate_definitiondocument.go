package iotthingsgraph

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// FlowTemplate_DefinitionDocument AWS CloudFormation Resource (AWS::IoTThingsGraph::FlowTemplate.DefinitionDocument)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotthingsgraph-flowtemplate-definitiondocument.html
type FlowTemplate_DefinitionDocument struct {

	// Language AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotthingsgraph-flowtemplate-definitiondocument.html#cfn-iotthingsgraph-flowtemplate-definitiondocument-language
	Language *types.Value `json:"Language,omitempty"`

	// Text AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotthingsgraph-flowtemplate-definitiondocument.html#cfn-iotthingsgraph-flowtemplate-definitiondocument-text
	Text *types.Value `json:"Text,omitempty"`

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
func (r *FlowTemplate_DefinitionDocument) AWSCloudFormationType() string {
	return "AWS::IoTThingsGraph::FlowTemplate.DefinitionDocument"
}
