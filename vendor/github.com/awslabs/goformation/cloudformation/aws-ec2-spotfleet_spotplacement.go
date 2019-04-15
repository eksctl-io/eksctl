package cloudformation

import (
	"encoding/json"
)

// AWSEC2SpotFleet_SpotPlacement AWS CloudFormation Resource (AWS::EC2::SpotFleet.SpotPlacement)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-spotfleetrequestconfigdata-launchspecifications-placement.html
type AWSEC2SpotFleet_SpotPlacement struct {

	// AvailabilityZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-spotfleetrequestconfigdata-launchspecifications-placement.html#cfn-ec2-spotfleet-spotplacement-availabilityzone
	AvailabilityZone *Value `json:"AvailabilityZone,omitempty"`

	// GroupName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-spotfleetrequestconfigdata-launchspecifications-placement.html#cfn-ec2-spotfleet-spotplacement-groupname
	GroupName *Value `json:"GroupName,omitempty"`

	// Tenancy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-spotfleetrequestconfigdata-launchspecifications-placement.html#cfn-ec2-spotfleet-spotplacement-tenancy
	Tenancy *Value `json:"Tenancy,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_SpotPlacement) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.SpotPlacement"
}

func (r *AWSEC2SpotFleet_SpotPlacement) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
