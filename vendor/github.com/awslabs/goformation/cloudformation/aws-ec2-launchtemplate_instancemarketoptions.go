package cloudformation

// AWSEC2LaunchTemplate_InstanceMarketOptions AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.InstanceMarketOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-instancemarketoptions.html
type AWSEC2LaunchTemplate_InstanceMarketOptions struct {

	// MarketType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-instancemarketoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-markettype
	MarketType *StringIntrinsic `json:"MarketType,omitempty"`

	// SpotOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplatedata-instancemarketoptions.html#cfn-ec2-launchtemplate-launchtemplatedata-instancemarketoptions-spotoptions
	SpotOptions *AWSEC2LaunchTemplate_SpotOptions `json:"SpotOptions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_InstanceMarketOptions) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.InstanceMarketOptions"
}
