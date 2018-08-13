package cloudformation

import (
	"encoding/json"
)

// AWSElasticLoadBalancingV2ListenerCertificate_Certificate AWS CloudFormation Resource (AWS::ElasticLoadBalancingV2::ListenerCertificate.Certificate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-certificates.html
type AWSElasticLoadBalancingV2ListenerCertificate_Certificate struct {

	// CertificateArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticloadbalancingv2-listener-certificates.html#cfn-elasticloadbalancingv2-listener-certificates-certificatearn
	CertificateArn *Value `json:"CertificateArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingV2ListenerCertificate_Certificate) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancingV2::ListenerCertificate.Certificate"
}

func (r *AWSElasticLoadBalancingV2ListenerCertificate_Certificate) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
