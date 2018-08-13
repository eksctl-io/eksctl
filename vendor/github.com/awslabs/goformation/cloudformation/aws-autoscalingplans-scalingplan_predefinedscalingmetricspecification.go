package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingPlansScalingPlan_PredefinedScalingMetricSpecification AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.PredefinedScalingMetricSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedscalingmetricspecification.html
type AWSAutoScalingPlansScalingPlan_PredefinedScalingMetricSpecification struct {

	// PredefinedScalingMetricType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedscalingmetricspecification.html#cfn-autoscalingplans-scalingplan-predefinedscalingmetricspecification-predefinedscalingmetrictype
	PredefinedScalingMetricType *Value `json:"PredefinedScalingMetricType,omitempty"`

	// ResourceLabel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedscalingmetricspecification.html#cfn-autoscalingplans-scalingplan-predefinedscalingmetricspecification-resourcelabel
	ResourceLabel *Value `json:"ResourceLabel,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingPlansScalingPlan_PredefinedScalingMetricSpecification) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.PredefinedScalingMetricSpecification"
}

func (r *AWSAutoScalingPlansScalingPlan_PredefinedScalingMetricSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
