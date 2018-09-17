package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingLoadBalancer_ConnectionDrainingPolicy AWS CloudFormation Resource (AWS::ElasticLoadBalancing::LoadBalancer.ConnectionDrainingPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-connectiondrainingpolicy.html
type AWSElasticLoadBalancingLoadBalancer_ConnectionDrainingPolicy struct {

	// Enabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-connectiondrainingpolicy.html#cfn-elb-connectiondrainingpolicy-enabled
	Enabled *Value `json:"Enabled,omitempty"`

	// Timeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-connectiondrainingpolicy.html#cfn-elb-connectiondrainingpolicy-timeout
	Timeout *Value `json:"Timeout,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingLoadBalancer_ConnectionDrainingPolicy) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancing::LoadBalancer.ConnectionDrainingPolicy"
}

func (r *AWSElasticLoadBalancingLoadBalancer_ConnectionDrainingPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
