package cloudformation

// AWSRoute53HostedZone_QueryLoggingConfig AWS CloudFormation Resource (AWS::Route53::HostedZone.QueryLoggingConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-hostedzone-queryloggingconfig.html
type AWSRoute53HostedZone_QueryLoggingConfig struct {

	// CloudWatchLogsLogGroupArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-hostedzone-queryloggingconfig.html#cfn-route53-hostedzone-queryloggingconfig-cloudwatchlogsloggrouparn
	CloudWatchLogsLogGroupArn *StringIntrinsic `json:"CloudWatchLogsLogGroupArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoute53HostedZone_QueryLoggingConfig) AWSCloudFormationType() string {
	return "AWS::Route53::HostedZone.QueryLoggingConfig"
}
