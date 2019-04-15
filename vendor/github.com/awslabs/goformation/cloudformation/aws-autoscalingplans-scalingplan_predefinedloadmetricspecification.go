package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingPlansScalingPlan_PredefinedLoadMetricSpecification AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.PredefinedLoadMetricSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedloadmetricspecification.html
type AWSAutoScalingPlansScalingPlan_PredefinedLoadMetricSpecification struct {

	// PredefinedLoadMetricType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedloadmetricspecification.html#cfn-autoscalingplans-scalingplan-predefinedloadmetricspecification-predefinedloadmetrictype
	PredefinedLoadMetricType *Value `json:"PredefinedLoadMetricType,omitempty"`

	// ResourceLabel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-predefinedloadmetricspecification.html#cfn-autoscalingplans-scalingplan-predefinedloadmetricspecification-resourcelabel
	ResourceLabel *Value `json:"ResourceLabel,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingPlansScalingPlan_PredefinedLoadMetricSpecification) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.PredefinedLoadMetricSpecification"
}

func (r *AWSAutoScalingPlansScalingPlan_PredefinedLoadMetricSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
