package cloudformation

import (
	"encoding/json"
)

// AWSElasticsearchDomain_VPCOptions AWS CloudFormation Resource (AWS::Elasticsearch::Domain.VPCOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-vpcoptions.html
type AWSElasticsearchDomain_VPCOptions struct {

	// SecurityGroupIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-vpcoptions.html#cfn-elasticsearch-domain-vpcoptions-securitygroupids
	SecurityGroupIds []*Value `json:"SecurityGroupIds,omitempty"`

	// SubnetIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-vpcoptions.html#cfn-elasticsearch-domain-vpcoptions-subnetids
	SubnetIds []*Value `json:"SubnetIds,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticsearchDomain_VPCOptions) AWSCloudFormationType() string {
	return "AWS::Elasticsearch::Domain.VPCOptions"
}

func (r *AWSElasticsearchDomain_VPCOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
