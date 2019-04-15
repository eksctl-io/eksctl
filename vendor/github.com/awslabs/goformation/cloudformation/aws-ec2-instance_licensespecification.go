package cloudformation

import (
	"encoding/json"
)

// AWSEC2Instance_LicenseSpecification AWS CloudFormation Resource (AWS::EC2::Instance.LicenseSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-licensespecification.html
type AWSEC2Instance_LicenseSpecification struct {

	// LicenseConfigurationArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-licensespecification.html#cfn-ec2-instance-licensespecification-licenseconfigurationarn
	LicenseConfigurationArn *Value `json:"LicenseConfigurationArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2Instance_LicenseSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::Instance.LicenseSpecification"
}

func (r *AWSEC2Instance_LicenseSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
