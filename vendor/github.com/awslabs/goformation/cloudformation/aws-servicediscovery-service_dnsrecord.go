package cloudformation

// AWSServiceDiscoveryService_DnsRecord AWS CloudFormation Resource (AWS::ServiceDiscovery::Service.DnsRecord)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-service-dnsrecord.html
type AWSServiceDiscoveryService_DnsRecord struct {

	// TTL AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-service-dnsrecord.html#cfn-servicediscovery-service-dnsrecord-ttl
	TTL *StringIntrinsic `json:"TTL,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicediscovery-service-dnsrecord.html#cfn-servicediscovery-service-dnsrecord-type
	Type *StringIntrinsic `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServiceDiscoveryService_DnsRecord) AWSCloudFormationType() string {
	return "AWS::ServiceDiscovery::Service.DnsRecord"
}
