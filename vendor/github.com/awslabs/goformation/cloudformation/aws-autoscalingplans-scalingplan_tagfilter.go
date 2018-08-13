package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingPlansScalingPlan_TagFilter AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.TagFilter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-tagfilter.html
type AWSAutoScalingPlansScalingPlan_TagFilter struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-tagfilter.html#cfn-autoscalingplans-scalingplan-tagfilter-key
	Key *Value `json:"Key,omitempty"`

	// Values AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-tagfilter.html#cfn-autoscalingplans-scalingplan-tagfilter-values
	Values []*Value `json:"Values,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingPlansScalingPlan_TagFilter) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.TagFilter"
}

func (r *AWSAutoScalingPlansScalingPlan_TagFilter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
