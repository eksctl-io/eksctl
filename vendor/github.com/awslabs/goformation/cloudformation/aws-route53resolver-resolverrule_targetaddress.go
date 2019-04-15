package cloudformation

import (
	"encoding/json"
)

// AWSRoute53ResolverResolverRule_TargetAddress AWS CloudFormation Resource (AWS::Route53Resolver::ResolverRule.TargetAddress)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-resolverrule-targetaddress.html
type AWSRoute53ResolverResolverRule_TargetAddress struct {

	// Ip AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-resolverrule-targetaddress.html#cfn-route53resolver-resolverrule-targetaddress-ip
	Ip *Value `json:"Ip,omitempty"`

	// Port AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-resolverrule-targetaddress.html#cfn-route53resolver-resolverrule-targetaddress-port
	Port *Value `json:"Port,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoute53ResolverResolverRule_TargetAddress) AWSCloudFormationType() string {
	return "AWS::Route53Resolver::ResolverRule.TargetAddress"
}

func (r *AWSRoute53ResolverResolverRule_TargetAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
