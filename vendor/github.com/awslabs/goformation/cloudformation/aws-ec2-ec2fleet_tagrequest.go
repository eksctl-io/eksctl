package cloudformation

import (
	"encoding/json"
)

// AWSEC2EC2Fleet_TagRequest AWS CloudFormation Resource (AWS::EC2::EC2Fleet.TagRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-tagrequest.html
type AWSEC2EC2Fleet_TagRequest struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-tagrequest.html#cfn-ec2-ec2fleet-tagrequest-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-tagrequest.html#cfn-ec2-ec2fleet-tagrequest-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2EC2Fleet_TagRequest) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.TagRequest"
}

func (r *AWSEC2EC2Fleet_TagRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
