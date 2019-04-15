package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingPlansScalingPlan_ScalingInstruction AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.ScalingInstruction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html
type AWSAutoScalingPlansScalingPlan_ScalingInstruction struct {

	// CustomizedLoadMetricSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-customizedloadmetricspecification
	CustomizedLoadMetricSpecification *AWSAutoScalingPlansScalingPlan_CustomizedLoadMetricSpecification `json:"CustomizedLoadMetricSpecification,omitempty"`

	// DisableDynamicScaling AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-disabledynamicscaling
	DisableDynamicScaling *Value `json:"DisableDynamicScaling,omitempty"`

	// MaxCapacity AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-maxcapacity
	MaxCapacity *Value `json:"MaxCapacity,omitempty"`

	// MinCapacity AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-mincapacity
	MinCapacity *Value `json:"MinCapacity,omitempty"`

	// PredefinedLoadMetricSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-predefinedloadmetricspecification
	PredefinedLoadMetricSpecification *AWSAutoScalingPlansScalingPlan_PredefinedLoadMetricSpecification `json:"PredefinedLoadMetricSpecification,omitempty"`

	// PredictiveScalingMaxCapacityBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-predictivescalingmaxcapacitybehavior
	PredictiveScalingMaxCapacityBehavior *Value `json:"PredictiveScalingMaxCapacityBehavior,omitempty"`

	// PredictiveScalingMaxCapacityBuffer AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-predictivescalingmaxcapacitybuffer
	PredictiveScalingMaxCapacityBuffer *Value `json:"PredictiveScalingMaxCapacityBuffer,omitempty"`

	// PredictiveScalingMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-predictivescalingmode
	PredictiveScalingMode *Value `json:"PredictiveScalingMode,omitempty"`

	// ResourceId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-resourceid
	ResourceId *Value `json:"ResourceId,omitempty"`

	// ScalableDimension AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-scalabledimension
	ScalableDimension *Value `json:"ScalableDimension,omitempty"`

	// ScalingPolicyUpdateBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-scalingpolicyupdatebehavior
	ScalingPolicyUpdateBehavior *Value `json:"ScalingPolicyUpdateBehavior,omitempty"`

	// ScheduledActionBufferTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-scheduledactionbuffertime
	ScheduledActionBufferTime *Value `json:"ScheduledActionBufferTime,omitempty"`

	// ServiceNamespace AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-servicenamespace
	ServiceNamespace *Value `json:"ServiceNamespace,omitempty"`

	// TargetTrackingConfigurations AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-scalinginstruction.html#cfn-autoscalingplans-scalingplan-scalinginstruction-targettrackingconfigurations
	TargetTrackingConfigurations []AWSAutoScalingPlansScalingPlan_TargetTrackingConfiguration `json:"TargetTrackingConfigurations,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingPlansScalingPlan_ScalingInstruction) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.ScalingInstruction"
}

func (r *AWSAutoScalingPlansScalingPlan_ScalingInstruction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
