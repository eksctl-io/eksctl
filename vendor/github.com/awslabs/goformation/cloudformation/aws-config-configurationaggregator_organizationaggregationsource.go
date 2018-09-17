package cloudformation

import (
	"encoding/json"
)

// AWSConfigConfigurationAggregator_OrganizationAggregationSource AWS CloudFormation Resource (AWS::Config::ConfigurationAggregator.OrganizationAggregationSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-organizationaggregationsource.html
type AWSConfigConfigurationAggregator_OrganizationAggregationSource struct {

	// AllAwsRegions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-organizationaggregationsource.html#cfn-config-configurationaggregator-organizationaggregationsource-allawsregions
	AllAwsRegions *Value `json:"AllAwsRegions,omitempty"`

	// AwsRegions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-organizationaggregationsource.html#cfn-config-configurationaggregator-organizationaggregationsource-awsregions
	AwsRegions []*Value `json:"AwsRegions,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationaggregator-organizationaggregationsource.html#cfn-config-configurationaggregator-organizationaggregationsource-rolearn
	RoleArn *Value `json:"RoleArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSConfigConfigurationAggregator_OrganizationAggregationSource) AWSCloudFormationType() string {
	return "AWS::Config::ConfigurationAggregator.OrganizationAggregationSource"
}

func (r *AWSConfigConfigurationAggregator_OrganizationAggregationSource) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
