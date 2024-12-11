package dynamodb

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GlobalTable_TimeToLiveSpecification AWS CloudFormation Resource (AWS::DynamoDB::GlobalTable.TimeToLiveSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-timetolivespecification.html
type GlobalTable_TimeToLiveSpecification struct {

	// AttributeName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-timetolivespecification.html#cfn-dynamodb-globaltable-timetolivespecification-attributename
	AttributeName *types.Value `json:"AttributeName,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-dynamodb-globaltable-timetolivespecification.html#cfn-dynamodb-globaltable-timetolivespecification-enabled
	Enabled *types.Value `json:"Enabled"`

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
func (r *GlobalTable_TimeToLiveSpecification) AWSCloudFormationType() string {
	return "AWS::DynamoDB::GlobalTable.TimeToLiveSpecification"
}
