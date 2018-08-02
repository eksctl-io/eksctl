package cloudformation

// AWSEC2LaunchTemplate_CreditSpecification AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.CreditSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-creditspecification.html
type AWSEC2LaunchTemplate_CreditSpecification struct {

	// CpuCredits AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-creditspecification.html#cfn-ec2-launchtemplate-launchtemplatedata-creditspecification-cpucredits
	CpuCredits *StringIntrinsic `json:"CpuCredits,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_CreditSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.CreditSpecification"
}
