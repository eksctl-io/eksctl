package cloudformation

// AWSEC2LaunchTemplate_Monitoring AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.Monitoring)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-monitoring.html
type AWSEC2LaunchTemplate_Monitoring struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-monitoring.html#cfn-ec2-launchtemplate-launchtemplatedata-monitoring-enabled
	Enabled bool `json:"Enabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_Monitoring) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.Monitoring"
}
