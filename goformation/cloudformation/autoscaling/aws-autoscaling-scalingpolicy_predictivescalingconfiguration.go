package autoscaling

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ScalingPolicy_PredictiveScalingConfiguration AWS CloudFormation Resource (AWS::AutoScaling::ScalingPolicy.PredictiveScalingConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predictivescalingconfiguration.html
type ScalingPolicy_PredictiveScalingConfiguration struct {

	// MaxCapacityBreachBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predictivescalingconfiguration.html#cfn-autoscaling-scalingpolicy-predictivescalingconfiguration-maxcapacitybreachbehavior
	MaxCapacityBreachBehavior *types.Value `json:"MaxCapacityBreachBehavior,omitempty"`

	// MaxCapacityBuffer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predictivescalingconfiguration.html#cfn-autoscaling-scalingpolicy-predictivescalingconfiguration-maxcapacitybuffer
	MaxCapacityBuffer *types.Value `json:"MaxCapacityBuffer,omitempty"`

	// MetricSpecifications AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predictivescalingconfiguration.html#cfn-autoscaling-scalingpolicy-predictivescalingconfiguration-metricspecifications
	MetricSpecifications []ScalingPolicy_PredictiveScalingMetricSpecification `json:"MetricSpecifications,omitempty"`

	// Mode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predictivescalingconfiguration.html#cfn-autoscaling-scalingpolicy-predictivescalingconfiguration-mode
	Mode *types.Value `json:"Mode,omitempty"`

	// SchedulingBufferTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predictivescalingconfiguration.html#cfn-autoscaling-scalingpolicy-predictivescalingconfiguration-schedulingbuffertime
	SchedulingBufferTime *types.Value `json:"SchedulingBufferTime,omitempty"`

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
func (r *ScalingPolicy_PredictiveScalingConfiguration) AWSCloudFormationType() string {
	return "AWS::AutoScaling::ScalingPolicy.PredictiveScalingConfiguration"
}
