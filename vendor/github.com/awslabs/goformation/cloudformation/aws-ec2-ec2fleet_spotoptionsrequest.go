package cloudformation

import (
	"encoding/json"
)

// AWSEC2EC2Fleet_SpotOptionsRequest AWS CloudFormation Resource (AWS::EC2::EC2Fleet.SpotOptionsRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html
type AWSEC2EC2Fleet_SpotOptionsRequest struct {

	// AllocationStrategy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-allocationstrategy
	AllocationStrategy *Value `json:"AllocationStrategy,omitempty"`

	// InstanceInterruptionBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-instanceinterruptionbehavior
	InstanceInterruptionBehavior *Value `json:"InstanceInterruptionBehavior,omitempty"`

	// InstancePoolsToUseCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-spotoptionsrequest.html#cfn-ec2-ec2fleet-spotoptionsrequest-instancepoolstousecount
	InstancePoolsToUseCount *Value `json:"InstancePoolsToUseCount,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2EC2Fleet_SpotOptionsRequest) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.SpotOptionsRequest"
}

func (r *AWSEC2EC2Fleet_SpotOptionsRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
