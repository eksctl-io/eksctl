package cloudformation

import (
	"encoding/json"
)

// AWSEC2SpotFleet_ClassicLoadBalancer AWS CloudFormation Resource (AWS::EC2::SpotFleet.ClassicLoadBalancer)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-classicloadbalancer.html
type AWSEC2SpotFleet_ClassicLoadBalancer struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-classicloadbalancer.html#cfn-ec2-spotfleet-classicloadbalancer-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_ClassicLoadBalancer) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.ClassicLoadBalancer"
}

func (r *AWSEC2SpotFleet_ClassicLoadBalancer) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
