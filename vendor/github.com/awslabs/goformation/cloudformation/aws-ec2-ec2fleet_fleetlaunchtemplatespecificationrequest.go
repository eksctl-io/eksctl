package cloudformation

import (
	"encoding/json"
)

// AWSEC2EC2Fleet_FleetLaunchTemplateSpecificationRequest AWS CloudFormation Resource (AWS::EC2::EC2Fleet.FleetLaunchTemplateSpecificationRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-fleetlaunchtemplatespecificationrequest.html
type AWSEC2EC2Fleet_FleetLaunchTemplateSpecificationRequest struct {

	// LaunchTemplateId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-fleetlaunchtemplatespecificationrequest.html#cfn-ec2-ec2fleet-fleetlaunchtemplatespecificationrequest-launchtemplateid
	LaunchTemplateId *Value `json:"LaunchTemplateId,omitempty"`

	// LaunchTemplateName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-fleetlaunchtemplatespecificationrequest.html#cfn-ec2-ec2fleet-fleetlaunchtemplatespecificationrequest-launchtemplatename
	LaunchTemplateName *Value `json:"LaunchTemplateName,omitempty"`

	// Version AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-fleetlaunchtemplatespecificationrequest.html#cfn-ec2-ec2fleet-fleetlaunchtemplatespecificationrequest-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2EC2Fleet_FleetLaunchTemplateSpecificationRequest) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.FleetLaunchTemplateSpecificationRequest"
}

func (r *AWSEC2EC2Fleet_FleetLaunchTemplateSpecificationRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
