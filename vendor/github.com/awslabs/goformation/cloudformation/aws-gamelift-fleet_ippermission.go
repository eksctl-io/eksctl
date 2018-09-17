package cloudformation

import (
	"encoding/json"
)

// AWSGameLiftFleet_IpPermission AWS CloudFormation Resource (AWS::GameLift::Fleet.IpPermission)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-ec2inboundpermission.html
type AWSGameLiftFleet_IpPermission struct {

	// FromPort AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-ec2inboundpermission.html#cfn-gamelift-fleet-ec2inboundpermissions-fromport
	FromPort *Value `json:"FromPort,omitempty"`

	// IpRange AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-ec2inboundpermission.html#cfn-gamelift-fleet-ec2inboundpermissions-iprange
	IpRange *Value `json:"IpRange,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-ec2inboundpermission.html#cfn-gamelift-fleet-ec2inboundpermissions-protocol
	Protocol *Value `json:"Protocol,omitempty"`

	// ToPort AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-gamelift-fleet-ec2inboundpermission.html#cfn-gamelift-fleet-ec2inboundpermissions-toport
	ToPort *Value `json:"ToPort,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGameLiftFleet_IpPermission) AWSCloudFormationType() string {
	return "AWS::GameLift::Fleet.IpPermission"
}

func (r *AWSGameLiftFleet_IpPermission) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
