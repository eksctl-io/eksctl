package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingLoadBalancer_LBCookieStickinessPolicy AWS CloudFormation Resource (AWS::ElasticLoadBalancing::LoadBalancer.LBCookieStickinessPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-LBCookieStickinessPolicy.html
type AWSElasticLoadBalancingLoadBalancer_LBCookieStickinessPolicy struct {

	// CookieExpirationPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-LBCookieStickinessPolicy.html#cfn-elb-lbcookiestickinesspolicy-cookieexpirationperiod
	CookieExpirationPeriod *Value `json:"CookieExpirationPeriod,omitempty"`

	// PolicyName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-LBCookieStickinessPolicy.html#cfn-elb-lbcookiestickinesspolicy-policyname
	PolicyName *Value `json:"PolicyName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingLoadBalancer_LBCookieStickinessPolicy) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancing::LoadBalancer.LBCookieStickinessPolicy"
}

func (r *AWSElasticLoadBalancingLoadBalancer_LBCookieStickinessPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
