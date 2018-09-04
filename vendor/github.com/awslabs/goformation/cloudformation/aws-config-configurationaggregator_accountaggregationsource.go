package cloudformation

// AWSConfigConfigurationAggregator_AccountAggregationSource AWS CloudFormation Resource (AWS::Config::ConfigurationAggregator.AccountAggregationSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-accountaggregationsource.html
type AWSConfigConfigurationAggregator_AccountAggregationSource struct {

	// AccountIds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-accountaggregationsource.html#cfn-config-configurationaggregator-accountaggregationsource-accountids
	AccountIds []*StringIntrinsic `json:"AccountIds,omitempty"`

	// AllAwsRegions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-accountaggregationsource.html#cfn-config-configurationaggregator-accountaggregationsource-allawsregions
	AllAwsRegions bool `json:"AllAwsRegions,omitempty"`

	// AwsRegions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-accountaggregationsource.html#cfn-config-configurationaggregator-accountaggregationsource-awsregions
	AwsRegions []*StringIntrinsic `json:"AwsRegions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSConfigConfigurationAggregator_AccountAggregationSource) AWSCloudFormationType() string {
	return "AWS::Config::ConfigurationAggregator.AccountAggregationSource"
}
