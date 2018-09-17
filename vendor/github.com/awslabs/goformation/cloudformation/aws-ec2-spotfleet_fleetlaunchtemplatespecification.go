package cloudformation

import (
	"encoding/json"
)

// AWSEC2SpotFleet_FleetLaunchTemplateSpecification AWS CloudFormation Resource (AWS::EC2::SpotFleet.FleetLaunchTemplateSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-fleetlaunchtemplatespecification.html
type AWSEC2SpotFleet_FleetLaunchTemplateSpecification struct {

	// LaunchTemplateId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-fleetlaunchtemplatespecification.html#cfn-ec2-spotfleet-fleetlaunchtemplatespecification-launchtemplateid
	LaunchTemplateId *Value `json:"LaunchTemplateId,omitempty"`

	// LaunchTemplateName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-fleetlaunchtemplatespecification.html#cfn-ec2-spotfleet-fleetlaunchtemplatespecification-launchtemplatename
	LaunchTemplateName *Value `json:"LaunchTemplateName,omitempty"`

	// Version AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-fleetlaunchtemplatespecification.html#cfn-ec2-spotfleet-fleetlaunchtemplatespecification-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_FleetLaunchTemplateSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.FleetLaunchTemplateSpecification"
}

func (r *AWSEC2SpotFleet_FleetLaunchTemplateSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
