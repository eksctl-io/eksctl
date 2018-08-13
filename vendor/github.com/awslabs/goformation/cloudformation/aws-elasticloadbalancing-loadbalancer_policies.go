package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingLoadBalancer_Policies AWS CloudFormation Resource (AWS::ElasticLoadBalancing::LoadBalancer.Policies)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-policy.html
type AWSElasticLoadBalancingLoadBalancer_Policies struct {

	// Attributes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-policy.html#cfn-ec2-elb-policy-attributes
	Attributes []interface{} `json:"Attributes,omitempty"`

	// InstancePorts AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-policy.html#cfn-ec2-elb-policy-instanceports
	InstancePorts []*Value `json:"InstancePorts,omitempty"`

	// LoadBalancerPorts AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-policy.html#cfn-ec2-elb-policy-loadbalancerports
	LoadBalancerPorts []*Value `json:"LoadBalancerPorts,omitempty"`

	// PolicyName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-policy.html#cfn-ec2-elb-policy-policyname
	PolicyName *Value `json:"PolicyName,omitempty"`

	// PolicyType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-policy.html#cfn-ec2-elb-policy-policytype
	PolicyType *Value `json:"PolicyType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingLoadBalancer_Policies) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancing::LoadBalancer.Policies"
}

func (r *AWSElasticLoadBalancingLoadBalancer_Policies) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
