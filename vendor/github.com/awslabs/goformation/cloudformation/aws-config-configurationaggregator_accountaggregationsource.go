package cloudformation

import (
	"encoding/json"
)

// AWSConfigConfigurationAggregator_AccountAggregationSource AWS CloudFormation Resource (AWS::Config::ConfigurationAggregator.AccountAggregationSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-accountaggregationsource.html
type AWSConfigConfigurationAggregator_AccountAggregationSource struct {

	// AccountIds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-accountaggregationsource.html#cfn-config-configurationaggregator-accountaggregationsource-accountids
	AccountIds []*Value `json:"AccountIds,omitempty"`

	// AllAwsRegions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-accountaggregationsource.html#cfn-config-configurationaggregator-accountaggregationsource-allawsregions
	AllAwsRegions *Value `json:"AllAwsRegions,omitempty"`

	// AwsRegions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-accountaggregationsource.html#cfn-config-configurationaggregator-accountaggregationsource-awsregions
	AwsRegions []*Value `json:"AwsRegions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSConfigConfigurationAggregator_AccountAggregationSource) AWSCloudFormationType() string {
	return "AWS::Config::ConfigurationAggregator.AccountAggregationSource"
}

func (r *AWSConfigConfigurationAggregator_AccountAggregationSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
