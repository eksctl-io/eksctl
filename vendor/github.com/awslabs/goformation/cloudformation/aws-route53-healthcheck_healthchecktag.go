package cloudformation

import (
	"encoding/json"
)

// AWSRoute53HealthCheck_HealthCheckTag AWS CloudFormation Resource (AWS::Route53::HealthCheck.HealthCheckTag)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-healthcheck-healthchecktag.html
type AWSRoute53HealthCheck_HealthCheckTag struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-healthcheck-healthchecktag.html#cfn-route53-healthchecktags-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-healthcheck-healthchecktag.html#cfn-route53-healthchecktags-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoute53HealthCheck_HealthCheckTag) AWSCloudFormationType() string {
	return "AWS::Route53::HealthCheck.HealthCheckTag"
}

func (r *AWSRoute53HealthCheck_HealthCheckTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
