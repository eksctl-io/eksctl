package cloudformation

import (
	"encoding/json"
)

// AWSOpsWorksLayer_AutoScalingThresholds AWS CloudFormation Resource (AWS::OpsWorks::Layer.AutoScalingThresholds)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-loadbasedautoscaling-autoscalingthresholds.html
type AWSOpsWorksLayer_AutoScalingThresholds struct {

	// CpuThreshold AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-loadbasedautoscaling-autoscalingthresholds.html#cfn-opsworks-layer-loadbasedautoscaling-autoscalingthresholds-cputhreshold
	CpuThreshold *Value `json:"CpuThreshold,omitempty"`

	// IgnoreMetricsTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-loadbasedautoscaling-autoscalingthresholds.html#cfn-opsworks-layer-loadbasedautoscaling-autoscalingthresholds-ignoremetricstime
	IgnoreMetricsTime *Value `json:"IgnoreMetricsTime,omitempty"`

	// InstanceCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-loadbasedautoscaling-autoscalingthresholds.html#cfn-opsworks-layer-loadbasedautoscaling-autoscalingthresholds-instancecount
	InstanceCount *Value `json:"InstanceCount,omitempty"`

	// LoadThreshold AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-loadbasedautoscaling-autoscalingthresholds.html#cfn-opsworks-layer-loadbasedautoscaling-autoscalingthresholds-loadthreshold
	LoadThreshold *Value `json:"LoadThreshold,omitempty"`

	// MemoryThreshold AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-loadbasedautoscaling-autoscalingthresholds.html#cfn-opsworks-layer-loadbasedautoscaling-autoscalingthresholds-memorythreshold
	MemoryThreshold *Value `json:"MemoryThreshold,omitempty"`

	// ThresholdsWaitTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-loadbasedautoscaling-autoscalingthresholds.html#cfn-opsworks-layer-loadbasedautoscaling-autoscalingthresholds-thresholdwaittime
	ThresholdsWaitTime *Value `json:"ThresholdsWaitTime,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSOpsWorksLayer_AutoScalingThresholds) AWSCloudFormationType() string {
	return "AWS::OpsWorks::Layer.AutoScalingThresholds"
}

func (r *AWSOpsWorksLayer_AutoScalingThresholds) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
