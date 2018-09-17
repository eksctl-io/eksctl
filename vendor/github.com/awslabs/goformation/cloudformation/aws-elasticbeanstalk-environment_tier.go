package cloudformation

import (
	"encoding/json"
)

// AWSElasticBeanstalkEnvironment_Tier AWS CloudFormation Resource (AWS::ElasticBeanstalk::Environment.Tier)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-beanstalk-environment-tier.html
type AWSElasticBeanstalkEnvironment_Tier struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-beanstalk-environment-tier.html#cfn-beanstalk-env-tier-name
	Name *Value `json:"Name,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-beanstalk-environment-tier.html#cfn-beanstalk-env-tier-type
	Type *Value `json:"Type,omitempty"`

	// Version AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-beanstalk-environment-tier.html#cfn-beanstalk-env-tier-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticBeanstalkEnvironment_Tier) AWSCloudFormationType() string {
	return "AWS::ElasticBeanstalk::Environment.Tier"
}

func (r *AWSElasticBeanstalkEnvironment_Tier) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
