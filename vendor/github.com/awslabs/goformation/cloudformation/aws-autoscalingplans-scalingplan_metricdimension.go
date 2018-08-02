package cloudformation

// AWSAutoScalingPlansScalingPlan_MetricDimension AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.MetricDimension)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-metricdimension.html
type AWSAutoScalingPlansScalingPlan_MetricDimension struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-metricdimension.html#cfn-autoscalingplans-scalingplan-metricdimension-name
	Name *StringIntrinsic `json:"Name,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-metricdimension.html#cfn-autoscalingplans-scalingplan-metricdimension-value
	Value *StringIntrinsic `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingPlansScalingPlan_MetricDimension) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.MetricDimension"
}
