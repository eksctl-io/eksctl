package cloudformation

import (
	"encoding/json"
)

// AWSEC2EC2Fleet_FleetLaunchTemplateConfigRequest AWS CloudFormation Resource (AWS::EC2::EC2Fleet.FleetLaunchTemplateConfigRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-fleetlaunchtemplateconfigrequest.html
type AWSEC2EC2Fleet_FleetLaunchTemplateConfigRequest struct {

	// LaunchTemplateSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-fleetlaunchtemplateconfigrequest.html#cfn-ec2-ec2fleet-fleetlaunchtemplateconfigrequest-launchtemplatespecification
	LaunchTemplateSpecification *AWSEC2EC2Fleet_FleetLaunchTemplateSpecificationRequest `json:"LaunchTemplateSpecification,omitempty"`

	// Overrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-fleetlaunchtemplateconfigrequest.html#cfn-ec2-ec2fleet-fleetlaunchtemplateconfigrequest-overrides
	Overrides []AWSEC2EC2Fleet_FleetLaunchTemplateOverridesRequest `json:"Overrides,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2EC2Fleet_FleetLaunchTemplateConfigRequest) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.FleetLaunchTemplateConfigRequest"
}

func (r *AWSEC2EC2Fleet_FleetLaunchTemplateConfigRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
