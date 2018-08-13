package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingPlansScalingPlan_TargetTrackingConfiguration AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.TargetTrackingConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-targettrackingconfiguration.html
type AWSAutoScalingPlansScalingPlan_TargetTrackingConfiguration struct {

	// CustomizedScalingMetricSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-targettrackingconfiguration.html#cfn-autoscalingplans-scalingplan-targettrackingconfiguration-customizedscalingmetricspecification
	CustomizedScalingMetricSpecification *AWSAutoScalingPlansScalingPlan_CustomizedScalingMetricSpecification `json:"CustomizedScalingMetricSpecification,omitempty"`

	// DisableScaleIn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-targettrackingconfiguration.html#cfn-autoscalingplans-scalingplan-targettrackingconfiguration-disablescalein
	DisableScaleIn *Value `json:"DisableScaleIn,omitempty"`

	// EstimatedInstanceWarmup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-targettrackingconfiguration.html#cfn-autoscalingplans-scalingplan-targettrackingconfiguration-estimatedinstancewarmup
	EstimatedInstanceWarmup *Value `json:"EstimatedInstanceWarmup,omitempty"`

	// PredefinedScalingMetricSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-targettrackingconfiguration.html#cfn-autoscalingplans-scalingplan-targettrackingconfiguration-predefinedscalingmetricspecification
	PredefinedScalingMetricSpecification *AWSAutoScalingPlansScalingPlan_PredefinedScalingMetricSpecification `json:"PredefinedScalingMetricSpecification,omitempty"`

	// ScaleInCooldown AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-targettrackingconfiguration.html#cfn-autoscalingplans-scalingplan-targettrackingconfiguration-scaleincooldown
	ScaleInCooldown *Value `json:"ScaleInCooldown,omitempty"`

	// ScaleOutCooldown AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-targettrackingconfiguration.html#cfn-autoscalingplans-scalingplan-targettrackingconfiguration-scaleoutcooldown
	ScaleOutCooldown *Value `json:"ScaleOutCooldown,omitempty"`

	// TargetValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-targettrackingconfiguration.html#cfn-autoscalingplans-scalingplan-targettrackingconfiguration-targetvalue
	TargetValue *Value `json:"TargetValue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingPlansScalingPlan_TargetTrackingConfiguration) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.TargetTrackingConfiguration"
}

func (r *AWSAutoScalingPlansScalingPlan_TargetTrackingConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
