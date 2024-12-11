package evidently

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Experiment_MetricGoalObject AWS CloudFormation Resource (AWS::Evidently::Experiment.MetricGoalObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-metricgoalobject.html
type Experiment_MetricGoalObject struct {

	// DesiredChange AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-metricgoalobject.html#cfn-evidently-experiment-metricgoalobject-desiredchange
	DesiredChange *types.Value `json:"DesiredChange,omitempty"`

	// EntityIdKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-metricgoalobject.html#cfn-evidently-experiment-metricgoalobject-entityidkey
	EntityIdKey *types.Value `json:"EntityIdKey,omitempty"`

	// EventPattern AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-metricgoalobject.html#cfn-evidently-experiment-metricgoalobject-eventpattern
	EventPattern *types.Value `json:"EventPattern,omitempty"`

	// MetricName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-metricgoalobject.html#cfn-evidently-experiment-metricgoalobject-metricname
	MetricName *types.Value `json:"MetricName,omitempty"`

	// UnitLabel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-metricgoalobject.html#cfn-evidently-experiment-metricgoalobject-unitlabel
	UnitLabel *types.Value `json:"UnitLabel,omitempty"`

	// ValueKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-experiment-metricgoalobject.html#cfn-evidently-experiment-metricgoalobject-valuekey
	ValueKey *types.Value `json:"ValueKey,omitempty"`

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
func (r *Experiment_MetricGoalObject) AWSCloudFormationType() string {
	return "AWS::Evidently::Experiment.MetricGoalObject"
}
