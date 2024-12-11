package timestream

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ScheduledQuery_MixedMeasureMapping AWS CloudFormation Resource (AWS::Timestream::ScheduledQuery.MixedMeasureMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-mixedmeasuremapping.html
type ScheduledQuery_MixedMeasureMapping struct {

	// MeasureName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-mixedmeasuremapping.html#cfn-timestream-scheduledquery-mixedmeasuremapping-measurename
	MeasureName *types.Value `json:"MeasureName,omitempty"`

	// MeasureValueType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-mixedmeasuremapping.html#cfn-timestream-scheduledquery-mixedmeasuremapping-measurevaluetype
	MeasureValueType *types.Value `json:"MeasureValueType,omitempty"`

	// MultiMeasureAttributeMappings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-mixedmeasuremapping.html#cfn-timestream-scheduledquery-mixedmeasuremapping-multimeasureattributemappings
	MultiMeasureAttributeMappings []ScheduledQuery_MultiMeasureAttributeMapping `json:"MultiMeasureAttributeMappings,omitempty"`

	// SourceColumn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-mixedmeasuremapping.html#cfn-timestream-scheduledquery-mixedmeasuremapping-sourcecolumn
	SourceColumn *types.Value `json:"SourceColumn,omitempty"`

	// TargetMeasureName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-timestream-scheduledquery-mixedmeasuremapping.html#cfn-timestream-scheduledquery-mixedmeasuremapping-targetmeasurename
	TargetMeasureName *types.Value `json:"TargetMeasureName,omitempty"`

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
func (r *ScheduledQuery_MixedMeasureMapping) AWSCloudFormationType() string {
	return "AWS::Timestream::ScheduledQuery.MixedMeasureMapping"
}
