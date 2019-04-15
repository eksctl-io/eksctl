package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingV2ListenerRule_FixedResponseConfig AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::ListenerRule.FixedResponseConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listenerrule-fixedresponseconfig.html
type AWSElasticLoadBalancingV2ListenerRule_FixedResponseConfig struct {

	// ContentType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listenerrule-fixedresponseconfig.html#cfn-elasticloadbalancingv2-listenerrule-fixedresponseconfig-contenttype
	ContentType *Value `json:"ContentType,omitempty"`

	// MessageBody AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listenerrule-fixedresponseconfig.html#cfn-elasticloadbalancingv2-listenerrule-fixedresponseconfig-messagebody
	MessageBody *Value `json:"MessageBody,omitempty"`

	// StatusCode AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listenerrule-fixedresponseconfig.html#cfn-elasticloadbalancingv2-listenerrule-fixedresponseconfig-statuscode
	StatusCode *Value `json:"StatusCode,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingV2ListenerRule_FixedResponseConfig) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::ListenerRule.FixedResponseConfig"
}

func (r *AWSElasticLoadBalancingV2ListenerRule_FixedResponseConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
