package cloudformation

// AWSEC2LaunchTemplate_SpotOptions AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.SpotOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-spotoptions.html
type AWSEC2LaunchTemplate_SpotOptions struct {

	// InstanceInterruptionBehavior AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-spotoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-spotoptions-instanceinterruptionbehavior
	InstanceInterruptionBehavior *StringIntrinsic `json:"InstanceInterruptionBehavior,omitempty"`

	// MaxPrice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-spotoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-spotoptions-maxprice
	MaxPrice *StringIntrinsic `json:"MaxPrice,omitempty"`

	// SpotInstanceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-spotoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-spotoptions-spotinstancetype
	SpotInstanceType *StringIntrinsic `json:"SpotInstanceType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_SpotOptions) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.SpotOptions"
}
