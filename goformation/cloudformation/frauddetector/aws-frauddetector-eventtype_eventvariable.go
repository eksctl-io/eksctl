package frauddetector

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/cloudformation"
	"goformation/v4/cloudformation/policies"
)

// EventType_EventVariable AWS CloudFormation Resource (AWS::FraudDetector::EventType.EventVariable)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html
type EventType_EventVariable struct {

	// Arn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-arn
	Arn *types.Value `json:"Arn,omitempty"`

	// CreatedTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-createdtime
	CreatedTime *types.Value `json:"CreatedTime,omitempty"`

	// DataSource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-datasource
	DataSource *types.Value `json:"DataSource,omitempty"`

	// DataType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-datatype
	DataType *types.Value `json:"DataType,omitempty"`

	// DefaultValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-defaultvalue
	DefaultValue *types.Value `json:"DefaultValue,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-description
	Description *types.Value `json:"Description,omitempty"`

	// Inline AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-inline
	Inline *types.Value `json:"Inline,omitempty"`

	// LastUpdatedTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-lastupdatedtime
	LastUpdatedTime *types.Value `json:"LastUpdatedTime,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-name
	Name *types.Value `json:"Name,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-tags
	Tags []cloudformation.Tag `json:"Tags,omitempty"`

	// VariableType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-frauddetector-eventtype-eventvariable.html#cfn-frauddetector-eventtype-eventvariable-variabletype
	VariableType *types.Value `json:"VariableType,omitempty"`

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
func (r *EventType_EventVariable) AWSCloudFormationType() string {
	return "AWS::FraudDetector::EventType.EventVariable"
}
