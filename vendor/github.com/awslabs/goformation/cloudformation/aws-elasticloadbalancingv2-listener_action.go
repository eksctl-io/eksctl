package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingV2Listener_Action AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::Listener.Action)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-defaultactions.html
type AWSElasticLoadBalancingV2Listener_Action struct {

	// TargetGroupArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-defaultactions.html#cfn-elasticloadbalancingv2-listener-defaultactions-targetgrouparn
	TargetGroupArn *Value `json:"TargetGroupArn,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-defaultactions.html#cfn-elasticloadbalancingv2-listener-defaultactions-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingV2Listener_Action) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::Listener.Action"
}

func (r *AWSElasticLoadBalancingV2Listener_Action) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
