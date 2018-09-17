package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingPlansScalingPlan_MetricDimension AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.MetricDimension)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-metricdimension.html
type AWSAutoScalingPlansScalingPlan_MetricDimension struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-metricdimension.html#cfn-autoscalingplans-scalingplan-metricdimension-name
	Name *Value `json:"Name,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-metricdimension.html#cfn-autoscalingplans-scalingplan-metricdimension-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingPlansScalingPlan_MetricDimension) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.MetricDimension"
}

func (r *AWSAutoScalingPlansScalingPlan_MetricDimension) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
