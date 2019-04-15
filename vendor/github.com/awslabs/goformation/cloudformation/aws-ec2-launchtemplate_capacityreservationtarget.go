package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_CapacityReservationTarget AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.CapacityReservationTarget)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-capacityreservationtarget.html
type AWSEC2LaunchTemplate_CapacityReservationTarget struct {

	// CapacityReservationId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-capacityreservationtarget.html#cfn-ec2-launchtemplate-capacityreservationtarget-capacityreservationid
	CapacityReservationId *Value `json:"CapacityReservationId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_CapacityReservationTarget) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.CapacityReservationTarget"
}

func (r *AWSEC2LaunchTemplate_CapacityReservationTarget) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
