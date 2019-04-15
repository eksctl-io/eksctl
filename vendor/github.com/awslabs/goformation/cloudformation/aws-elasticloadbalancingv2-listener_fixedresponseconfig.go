package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingV2Listener_FixedResponseConfig AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::Listener.FixedResponseConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-fixedresponseconfig.html
type AWSElasticLoadBalancingV2Listener_FixedResponseConfig struct {

	// ContentType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-fixedresponseconfig.html#cfn-elasticloadbalancingv2-listener-fixedresponseconfig-contenttype
	ContentType *Value `json:"ContentType,omitempty"`

	// MessageBody AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-fixedresponseconfig.html#cfn-elasticloadbalancingv2-listener-fixedresponseconfig-messagebody
	MessageBody *Value `json:"MessageBody,omitempty"`

	// StatusCode AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-fixedresponseconfig.html#cfn-elasticloadbalancingv2-listener-fixedresponseconfig-statuscode
	StatusCode *Value `json:"StatusCode,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingV2Listener_FixedResponseConfig) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::Listener.FixedResponseConfig"
}

func (r *AWSElasticLoadBalancingV2Listener_FixedResponseConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
