package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingLoadBalancer_HealthCheck AWS CloudFormation Resource (AWS::ElasticLoadBalancing::LoadBalancer.HealthCheck)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-health-check.html
type AWSElasticLoadBalancingLoadBalancer_HealthCheck struct {

	// HealthyThreshold AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-health-check.html#cfn-elb-healthcheck-healthythreshold
	HealthyThreshold *Value `json:"HealthyThreshold,omitempty"`

	// Interval AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-health-check.html#cfn-elb-healthcheck-interval
	Interval *Value `json:"Interval,omitempty"`

	// Target AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-health-check.html#cfn-elb-healthcheck-target
	Target *Value `json:"Target,omitempty"`

	// Timeout AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-health-check.html#cfn-elb-healthcheck-timeout
	Timeout *Value `json:"Timeout,omitempty"`

	// UnhealthyThreshold AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-health-check.html#cfn-elb-healthcheck-unhealthythreshold
	UnhealthyThreshold *Value `json:"UnhealthyThreshold,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingLoadBalancer_HealthCheck) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancing::LoadBalancer.HealthCheck"
}

func (r *AWSElasticLoadBalancingLoadBalancer_HealthCheck) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
