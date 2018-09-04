package cloudformation

// AWSEC2LaunchTemplate_PrivateIpAdd AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.PrivateIpAdd)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-privateipadd.html
type AWSEC2LaunchTemplate_PrivateIpAdd struct {

	// Primary AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-privateipadd.html#cfn-ec2-launchtemplate-privateipadd-primary
	Primary bool `json:"Primary,omitempty"`

	// PrivateIpAddress AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-privateipadd.html#cfn-ec2-launchtemplate-privateipadd-privateipaddress
	PrivateIpAddress *StringIntrinsic `json:"PrivateIpAddress,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_PrivateIpAdd) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.PrivateIpAdd"
}
