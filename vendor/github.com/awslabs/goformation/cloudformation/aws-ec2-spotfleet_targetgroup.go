package cloudformation

import (
	"encoding/json"
)

// AWSEC2SpotFleet_TargetGroup AWS CloudFormation Resource (AWS::EC2::SpotFleet.TargetGroup)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-targetgroup.html
type AWSEC2SpotFleet_TargetGroup struct {

	// Arn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-targetgroup.html#cfn-ec2-spotfleet-targetgroup-arn
	Arn *Value `json:"Arn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_TargetGroup) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.TargetGroup"
}

func (r *AWSEC2SpotFleet_TargetGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
