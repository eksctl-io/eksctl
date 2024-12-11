package customerprofiles

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Integration_Task AWS CloudFormation Resource (AWS::CustomerProfiles::Integration.Task)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-task.html
type Integration_Task struct {

	// ConnectorOperator AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-task.html#cfn-customerprofiles-integration-task-connectoroperator
	ConnectorOperator *Integration_ConnectorOperator `json:"ConnectorOperator,omitempty"`

	// DestinationField AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-task.html#cfn-customerprofiles-integration-task-destinationfield
	DestinationField *types.Value `json:"DestinationField,omitempty"`

	// SourceFields AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-task.html#cfn-customerprofiles-integration-task-sourcefields
	SourceFields *types.Value `json:"SourceFields,omitempty"`

	// TaskProperties AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-task.html#cfn-customerprofiles-integration-task-taskproperties
	TaskProperties []Integration_TaskPropertiesMap `json:"TaskProperties,omitempty"`

	// TaskType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-task.html#cfn-customerprofiles-integration-task-tasktype
	TaskType *types.Value `json:"TaskType,omitempty"`

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
func (r *Integration_Task) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::Integration.Task"
}
