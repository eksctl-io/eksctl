package timestream

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ScheduledQuery_MultiMeasureAttributeMapping AWS CloudFormation Resource (AWS::Timestream::ScheduledQuery.MultiMeasureAttributeMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-multimeasureattributemapping.html
type ScheduledQuery_MultiMeasureAttributeMapping struct {

	// MeasureValueType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-multimeasureattributemapping.html#cfn-timestream-scheduledquery-multimeasureattributemapping-measurevaluetype
	MeasureValueType *types.Value `json:"MeasureValueType,omitempty"`

	// SourceColumn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-multimeasureattributemapping.html#cfn-timestream-scheduledquery-multimeasureattributemapping-sourcecolumn
	SourceColumn *types.Value `json:"SourceColumn,omitempty"`

	// TargetMultiMeasureAttributeName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-multimeasureattributemapping.html#cfn-timestream-scheduledquery-multimeasureattributemapping-targetmultimeasureattributename
	TargetMultiMeasureAttributeName *types.Value `json:"TargetMultiMeasureAttributeName,omitempty"`

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
func (r *ScheduledQuery_MultiMeasureAttributeMapping) AWSCloudFormationType() string {
	return "AWS::Timestream::ScheduledQuery.MultiMeasureAttributeMapping"
}
