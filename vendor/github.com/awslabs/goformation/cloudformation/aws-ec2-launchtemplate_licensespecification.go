package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_LicenseSpecification AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.LicenseSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-licensespecification.html
type AWSEC2LaunchTemplate_LicenseSpecification struct {

	// LicenseConfigurationArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-licensespecification.html#cfn-ec2-launchtemplate-licensespecification-licenseconfigurationarn
	LicenseConfigurationArn *Value `json:"LicenseConfigurationArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_LicenseSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.LicenseSpecification"
}

func (r *AWSEC2LaunchTemplate_LicenseSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
