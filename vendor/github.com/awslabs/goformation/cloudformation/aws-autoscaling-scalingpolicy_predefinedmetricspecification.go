package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingScalingPolicy_PredefinedMetricSpecification AWS CloudFormation Resource (AWS::AutoScaling::ScalingPolicy.PredefinedMetricSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predefinedmetricspecification.html
type AWSAutoScalingScalingPolicy_PredefinedMetricSpecification struct {

	// PredefinedMetricType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predefinedmetricspecification.html#cfn-autoscaling-scalingpolicy-predefinedmetricspecification-predefinedmetrictype
	PredefinedMetricType *Value `json:"PredefinedMetricType,omitempty"`

	// ResourceLabel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-scalingpolicy-predefinedmetricspecification.html#cfn-autoscaling-scalingpolicy-predefinedmetricspecification-resourcelabel
	ResourceLabel *Value `json:"ResourceLabel,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingScalingPolicy_PredefinedMetricSpecification) AWSCloudFormationType() string {
	return "AWS::AutoScaling::ScalingPolicy.PredefinedMetricSpecification"
}

func (r *AWSAutoScalingScalingPolicy_PredefinedMetricSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
