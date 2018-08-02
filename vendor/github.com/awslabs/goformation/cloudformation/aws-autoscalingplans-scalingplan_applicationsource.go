package cloudformation

// AWSAutoScalingPlansScalingPlan_ApplicationSource AWS CloudFormation Resource (AWS::AutoScalingPlans::ScalingPlan.ApplicationSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-applicationsource.html
type AWSAutoScalingPlansScalingPlan_ApplicationSource struct {

	// CloudFormationStackARN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-applicationsource.html#cfn-autoscalingplans-scalingplan-applicationsource-cloudformationstackarn
	CloudFormationStackARN *StringIntrinsic `json:"CloudFormationStackARN,omitempty"`

	// TagFilters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscalingplans-scalingplan-applicationsource.html#cfn-autoscalingplans-scalingplan-applicationsource-tagfilters
	TagFilters []AWSAutoScalingPlansScalingPlan_TagFilter `json:"TagFilters,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingPlansScalingPlan_ApplicationSource) AWSCloudFormationType() string {
	return "AWS::AutoScalingPlans::ScalingPlan.ApplicationSource"
}
