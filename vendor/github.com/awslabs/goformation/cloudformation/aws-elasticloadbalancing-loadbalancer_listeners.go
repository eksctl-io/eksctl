package cloudformation

// AWSElasticLoadBalancingLoadBalancer_Listeners AWS CloudFormation Resource (AWS::ElasticLoadBalancing::LoadBalancer.Listeners)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-listener.html
type AWSElasticLoadBalancingLoadBalancer_Listeners struct {

	// InstancePort AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-listener.html#cfn-ec2-elb-listener-instanceport
	InstancePort *StringIntrinsic `json:"InstancePort,omitempty"`

	// InstanceProtocol AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-listener.html#cfn-ec2-elb-listener-instanceprotocol
	InstanceProtocol *StringIntrinsic `json:"InstanceProtocol,omitempty"`

	// LoadBalancerPort AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-listener.html#cfn-ec2-elb-listener-loadbalancerport
	LoadBalancerPort *StringIntrinsic `json:"LoadBalancerPort,omitempty"`

	// PolicyNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-listener.html#cfn-ec2-elb-listener-policynames
	PolicyNames []*StringIntrinsic `json:"PolicyNames,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-listener.html#cfn-ec2-elb-listener-protocol
	Protocol *StringIntrinsic `json:"Protocol,omitempty"`

	// SSLCertificateId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-elb-listener.html#cfn-ec2-elb-listener-sslcertificateid
	SSLCertificateId *StringIntrinsic `json:"SSLCertificateId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticLoadBalancingLoadBalancer_Listeners) AWSCloudFormationType() string {
	return "AWS::ElasticLoadBalancing::LoadBalancer.Listeners"
}
