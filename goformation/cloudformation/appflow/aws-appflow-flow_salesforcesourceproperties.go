package appflow

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Flow_SalesforceSourceProperties AWS CloudFormation Resource (AWS::AppFlow::Flow.SalesforceSourceProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-salesforcesourceproperties.html
type Flow_SalesforceSourceProperties struct {

	// EnableDynamicFieldUpdate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-salesforcesourceproperties.html#cfn-appflow-flow-salesforcesourceproperties-enabledynamicfieldupdate
	EnableDynamicFieldUpdate *types.Value `json:"EnableDynamicFieldUpdate,omitempty"`

	// IncludeDeletedRecords AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-salesforcesourceproperties.html#cfn-appflow-flow-salesforcesourceproperties-includedeletedrecords
	IncludeDeletedRecords *types.Value `json:"IncludeDeletedRecords,omitempty"`

	// Object AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appflow-flow-salesforcesourceproperties.html#cfn-appflow-flow-salesforcesourceproperties-object
	Object *types.Value `json:"Object,omitempty"`

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
func (r *Flow_SalesforceSourceProperties) AWSCloudFormationType() string {
	return "AWS::AppFlow::Flow.SalesforceSourceProperties"
}
