package cloudformation

// AWSEC2Instance_CreditSpecification AWS CloudFormation Resource (AWS::EC2::Instance.CreditSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-creditspecification.html
type AWSEC2Instance_CreditSpecification struct {

	// CPUCredits AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-creditspecification.html#cfn-ec2-instance-creditspecification-cpucredits
	CPUCredits *StringIntrinsic `json:"CPUCredits,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2Instance_CreditSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::Instance.CreditSpecification"
}
