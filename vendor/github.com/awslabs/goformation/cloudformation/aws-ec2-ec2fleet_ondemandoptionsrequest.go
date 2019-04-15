package cloudformation

import (
	"encoding/json"
)

// AWSEC2EC2Fleet_OnDemandOptionsRequest AWS CloudFormation Resource (AWS::EC2::EC2Fleet.OnDemandOptionsRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-ondemandoptionsrequest.html
type AWSEC2EC2Fleet_OnDemandOptionsRequest struct {

	// AllocationStrategy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-ondemandoptionsrequest.html#cfn-ec2-ec2fleet-ondemandoptionsrequest-allocationstrategy
	AllocationStrategy *Value `json:"AllocationStrategy,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2EC2Fleet_OnDemandOptionsRequest) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.OnDemandOptionsRequest"
}

func (r *AWSEC2EC2Fleet_OnDemandOptionsRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
