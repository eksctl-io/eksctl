package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_Ipv6Add AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.Ipv6Add)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ipv6add.html
type AWSEC2LaunchTemplate_Ipv6Add struct {

	// Ipv6Address AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-ipv6add.html#cfn-ec2-launchtemplate-ipv6add-ipv6address
	Ipv6Address *Value `json:"Ipv6Address,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_Ipv6Add) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.Ipv6Add"
}

func (r *AWSEC2LaunchTemplate_Ipv6Add) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
