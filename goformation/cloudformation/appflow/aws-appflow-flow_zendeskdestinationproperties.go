package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Flow_ZendeskDestinationProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.ZendeskDestinationProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-zendeskdestinationproperties.html
type Flow_ZendeskDestinationProperties struct {

	// ErrorHandlingConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-zendeskdestinationproperties.html#cfn-appflow-flow-zendeskdestinationproperties-errorhandlingconfig
	ErrorHandlingConfig *Flow_ErrorHandlingConfig `json:"ErrorHandlingConfig,omitempty"`

	// IdFieldNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-zendeskdestinationproperties.html#cfn-appflow-flow-zendeskdestinationproperties-idfieldnames
	IdFieldNames *types.Value `json:"IdFieldNames,omitempty"`

	// Object AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-zendeskdestinationproperties.html#cfn-appflow-flow-zendeskdestinationproperties-object
	Object *types.Value `json:"Object,omitempty"`

	// WriteOperationType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-zendeskdestinationproperties.html#cfn-appflow-flow-zendeskdestinationproperties-writeoperationtype
	WriteOperationType *types.Value `json:"WriteOperationType,omitempty"`

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
func (r *Flow_ZendeskDestinationProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.ZendeskDestinationProperties"
}
