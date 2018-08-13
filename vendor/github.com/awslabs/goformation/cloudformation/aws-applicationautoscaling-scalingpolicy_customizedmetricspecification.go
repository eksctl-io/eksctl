package cloudformation

import (
	"encoding/json"
)

// AWSApplicationAutoScalingScalingPolicy_CustomizedMetricSpecification AWS CloudFormation Resource (AWS::ApplicationAutoScaling::ScalingPolicy.CustomizedMetricSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-customizedmetricspecification.html
type AWSApplicationAutoScalingScalingPolicy_CustomizedMetricSpecification struct {

	// Dimensions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-customizedmetricspecification.html#cfn-applicationautoscaling-scalingpolicy-customizedmetricspecification-dimensions
	Dimensions []AWSApplicationAutoScalingScalingPolicy_MetricDimension `json:"Dimensions,omitempty"`

	// MetricName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-customizedmetricspecification.html#cfn-applicationautoscaling-scalingpolicy-customizedmetricspecification-metricname
	MetricName *Value `json:"MetricName,omitempty"`

	// Namespace AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-customizedmetricspecification.html#cfn-applicationautoscaling-scalingpolicy-customizedmetricspecification-namespace
	Namespace *Value `json:"Namespace,omitempty"`

	// Statistic AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-customizedmetricspecification.html#cfn-applicationautoscaling-scalingpolicy-customizedmetricspecification-statistic
	Statistic *Value `json:"Statistic,omitempty"`

	// Unit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-applicationautoscaling-scalingpolicy-customizedmetricspecification.html#cfn-applicationautoscaling-scalingpolicy-customizedmetricspecification-unit
	Unit *Value `json:"Unit,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApplicationAutoScalingScalingPolicy_CustomizedMetricSpecification) AWSCloudFormationType() string {
	return "AWS::ApplicationAutoScaling::ScalingPolicy.CustomizedMetricSpecification"
}

func (r *AWSApplicationAutoScalingScalingPolicy_CustomizedMetricSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
