package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_Placement AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.Placement)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-placement.html
type AWSEC2LaunchTemplate_Placement struct {

	// Affinity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-placement.html#cfn-ec2-launchtemplate-launchtemplatedata-placement-affinity
	Affinity *Value `json:"Affinity,omitempty"`

	// AvailabilityZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-placement.html#cfn-ec2-launchtemplate-launchtemplatedata-placement-availabilityzone
	AvailabilityZone *Value `json:"AvailabilityZone,omitempty"`

	// GroupName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-placement.html#cfn-ec2-launchtemplate-launchtemplatedata-placement-groupname
	GroupName *Value `json:"GroupName,omitempty"`

	// HostId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-placement.html#cfn-ec2-launchtemplate-launchtemplatedata-placement-hostid
	HostId *Value `json:"HostId,omitempty"`

	// Tenancy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-placement.html#cfn-ec2-launchtemplate-launchtemplatedata-placement-tenancy
	Tenancy *Value `json:"Tenancy,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_Placement) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.Placement"
}

func (r *AWSEC2LaunchTemplate_Placement) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
