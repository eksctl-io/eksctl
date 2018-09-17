package cloudformation

import (
	"encoding/json"
)

// AWSApplicationAutoScalingScalingPolicy_StepScalingPolicyConfiguration AWS CloudFormation Resource (AWS::ApplicationAutoScaling::ScalingPolicy.StepScalingPolicyConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration.html
type AWSApplicationAutoScalingScalingPolicy_StepScalingPolicyConfiguration struct {

	// AdjustmentType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration.html#cfn-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-adjustmenttype
	AdjustmentType *Value `json:"AdjustmentType,omitempty"`

	// Cooldown AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration.html#cfn-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-cooldown
	Cooldown *Value `json:"Cooldown,omitempty"`

	// MetricAggregationType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration.html#cfn-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-metricaggregationtype
	MetricAggregationType *Value `json:"MetricAggregationType,omitempty"`

	// MinAdjustmentMagnitude AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration.html#cfn-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-minadjustmentmagnitude
	MinAdjustmentMagnitude *Value `json:"MinAdjustmentMagnitude,omitempty"`

	// StepAdjustments AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration.html#cfn-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-stepadjustments
	StepAdjustments []AWSApplicationAutoScalingScalingPolicy_StepAdjustment `json:"StepAdjustments,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApplicationAutoScalingScalingPolicy_StepScalingPolicyConfiguration) AWSCloudFormationType() string {
	return "AWS::ApplicationAutoScaling::ScalingPolicy.StepScalingPolicyConfiguration"
}

func (r *AWSApplicationAutoScalingScalingPolicy_StepScalingPolicyConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
