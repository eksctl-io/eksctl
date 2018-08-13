package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingV2TargetGroup_TargetGroupAttribute AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::TargetGroup.TargetGroupAttribute)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-targetgroup-targetgroupattribute.html
type AWSElasticLoadBalancingV2TargetGroup_TargetGroupAttribute struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-targetgroup-targetgroupattribute.html#cfn-elasticloadbalancingv2-targetgroup-targetgroupattribute-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-targetgroup-targetgroupattribute.html#cfn-elasticloadbalancingv2-targetgroup-targetgroupattribute-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingV2TargetGroup_TargetGroupAttribute) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::TargetGroup.TargetGroupAttribute"
}

func (r *AWSElasticLoadBalancingV2TargetGroup_TargetGroupAttribute) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
