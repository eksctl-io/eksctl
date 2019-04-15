package cloudformation

import (
	"encoding/json"
)

// AWSEC2SpotFleet_LoadBalancersConfig AWS CloudFormation Resource (AWS::EC2::SpotFleet.LoadBalancersConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-loadbalancersconfig.html
type AWSEC2SpotFleet_LoadBalancersConfig struct {

	// ClassicLoadBalancersConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-loadbalancersconfig.html#cfn-ec2-spotfleet-loadbalancersconfig-classicloadbalancersconfig
	ClassicLoadBalancersConfig *AWSEC2SpotFleet_ClassicLoadBalancersConfig `json:"ClassicLoadBalancersConfig,omitempty"`

	// TargetGroupsConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-loadbalancersconfig.html#cfn-ec2-spotfleet-loadbalancersconfig-targetgroupsconfig
	TargetGroupsConfig *AWSEC2SpotFleet_TargetGroupsConfig `json:"TargetGroupsConfig,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_LoadBalancersConfig) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.LoadBalancersConfig"
}

func (r *AWSEC2SpotFleet_LoadBalancersConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
