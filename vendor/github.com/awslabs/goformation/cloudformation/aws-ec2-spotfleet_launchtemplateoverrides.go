package cloudformation

import (
	"encoding/json"
)

// AWSEC2SpotFleet_LaunchTemplateOverrides AWS CloudFormation Resource (AWS::EC2::SpotFleet.LaunchTemplateOverrides)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateoverrides.html
type AWSEC2SpotFleet_LaunchTemplateOverrides struct {

	// AvailabilityZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateoverrides.html#cfn-ec2-spotfleet-launchtemplateoverrides-availabilityzone
	AvailabilityZone *Value `json:"AvailabilityZone,omitempty"`

	// InstanceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateoverrides.html#cfn-ec2-spotfleet-launchtemplateoverrides-instancetype
	InstanceType *Value `json:"InstanceType,omitempty"`

	// SpotPrice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateoverrides.html#cfn-ec2-spotfleet-launchtemplateoverrides-spotprice
	SpotPrice *Value `json:"SpotPrice,omitempty"`

	// SubnetId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateoverrides.html#cfn-ec2-spotfleet-launchtemplateoverrides-subnetid
	SubnetId *Value `json:"SubnetId,omitempty"`

	// WeightedCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateoverrides.html#cfn-ec2-spotfleet-launchtemplateoverrides-weightedcapacity
	WeightedCapacity *Value `json:"WeightedCapacity,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_LaunchTemplateOverrides) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.LaunchTemplateOverrides"
}

func (r *AWSEC2SpotFleet_LaunchTemplateOverrides) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
