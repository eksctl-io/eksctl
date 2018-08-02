package cloudformation

// AWSEC2LaunchTemplate_BlockDeviceMapping AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.BlockDeviceMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-blockdevicemapping.html
type AWSEC2LaunchTemplate_BlockDeviceMapping struct {

	// DeviceName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-blockdevicemapping.html#cfn-ec2-launchtemplate-blockdevicemapping-devicename
	DeviceName *StringIntrinsic `json:"DeviceName,omitempty"`

	// Ebs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-blockdevicemapping.html#cfn-ec2-launchtemplate-blockdevicemapping-ebs
	Ebs *AWSEC2LaunchTemplate_Ebs `json:"Ebs,omitempty"`

	// NoDevice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-blockdevicemapping.html#cfn-ec2-launchtemplate-blockdevicemapping-nodevice
	NoDevice *StringIntrinsic `json:"NoDevice,omitempty"`

	// VirtualName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-blockdevicemapping.html#cfn-ec2-launchtemplate-blockdevicemapping-virtualname
	VirtualName *StringIntrinsic `json:"VirtualName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_BlockDeviceMapping) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.BlockDeviceMapping"
}
