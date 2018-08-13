package cloudformation

import (
	"encoding/json"
)

// AWSApplicationAutoScalingScalingPolicy_StepAdjustment AWS CloudFormation Resource (AWS::ApplicationAutoScaling::ScalingPolicy.StepAdjustment)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-stepadjustment.html
type AWSApplicationAutoScalingScalingPolicy_StepAdjustment struct {

	// MetricIntervalLowerBound AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-stepadjustment.html#cfn-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-stepadjustment-metricintervallowerbound
	MetricIntervalLowerBound *Value `json:"MetricIntervalLowerBound,omitempty"`

	// MetricIntervalUpperBound AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-stepadjustment.html#cfn-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-stepadjustment-metricintervalupperbound
	MetricIntervalUpperBound *Value `json:"MetricIntervalUpperBound,omitempty"`

	// ScalingAdjustment AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-stepadjustment.html#cfn-applicationautoscaling-scalingpolicy-stepscalingpolicyconfiguration-stepadjustment-scalingadjustment
	ScalingAdjustment *Value `json:"ScalingAdjustment,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApplicationAutoScalingScalingPolicy_StepAdjustment) AWSCloudFormationType() string {
	return "AWS::ApplicationAutoScaling::ScalingPolicy.StepAdjustment"
}

func (r *AWSApplicationAutoScalingScalingPolicy_StepAdjustment) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
