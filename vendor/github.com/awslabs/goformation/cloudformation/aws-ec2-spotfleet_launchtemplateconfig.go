package cloudformation

// AWSEC2SpotFleet_LaunchTemplateConfig AWS CloudFormation Resource (AWS::EC2::SpotFleet.LaunchTemplateConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateconfig.html
type AWSEC2SpotFleet_LaunchTemplateConfig struct {

	// LaunchTemplateSpecification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateconfig.html#cfn-ec2-spotfleet-launchtemplateconfig-launchtemplatespecification
	LaunchTemplateSpecification *AWSEC2SpotFleet_FleetLaunchTemplateSpecification `json:"LaunchTemplateSpecification,omitempty"`

	// Overrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-spotfleet-launchtemplateconfig.html#cfn-ec2-spotfleet-launchtemplateconfig-overrides
	Overrides []AWSEC2SpotFleet_LaunchTemplateOverrides `json:"Overrides,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2SpotFleet_LaunchTemplateConfig) AWSCloudFormationType() string {
	return "AWS::EC2::SpotFleet.LaunchTemplateConfig"
}
