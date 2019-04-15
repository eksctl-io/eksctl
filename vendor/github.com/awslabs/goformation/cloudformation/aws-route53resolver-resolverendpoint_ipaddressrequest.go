package cloudformation

import (
	"encoding/json"
)

// AWSRoute53ResolverResolverEndpoint_IpAddressRequest AWS CloudFormation Resource (AWS::Route53Resolver::ResolverEndpoint.IpAddressRequest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-resolverendpoint-ipaddressrequest.html
type AWSRoute53ResolverResolverEndpoint_IpAddressRequest struct {

	// Ip AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-resolverendpoint-ipaddressrequest.html#cfn-route53resolver-resolverendpoint-ipaddressrequest-ip
	Ip *Value `json:"Ip,omitempty"`

	// SubnetId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53resolver-resolverendpoint-ipaddressrequest.html#cfn-route53resolver-resolverendpoint-ipaddressrequest-subnetid
	SubnetId *Value `json:"SubnetId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoute53ResolverResolverEndpoint_IpAddressRequest) AWSCloudFormationType() string {
	return "AWS::Route53Resolver::ResolverEndpoint.IpAddressRequest"
}

func (r *AWSRoute53ResolverResolverEndpoint_IpAddressRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
