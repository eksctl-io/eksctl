package cloudformation

import (
	"encoding/json"
)

// AWSRoute53HostedZone_VPC AWS CloudFormation Resource (AWS::Route53::HostedZone.VPC)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53-hostedzone-hostedzonevpcs.html
type AWSRoute53HostedZone_VPC struct {

	// VPCId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53-hostedzone-hostedzonevpcs.html#cfn-route53-hostedzone-hostedzonevpcs-vpcid
	VPCId *Value `json:"VPCId,omitempty"`

	// VPCRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53-hostedzone-hostedzonevpcs.html#cfn-route53-hostedzone-hostedzonevpcs-vpcregion
	VPCRegion *Value `json:"VPCRegion,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoute53HostedZone_VPC) AWSCloudFormationType() string {
	return "AWS::Route53::HostedZone.VPC"
}

func (r *AWSRoute53HostedZone_VPC) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
