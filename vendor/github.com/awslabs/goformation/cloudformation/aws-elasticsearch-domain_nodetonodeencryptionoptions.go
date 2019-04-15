package cloudformation

import (
	"encoding/json"
)

// AWSElasticsearchDomain_NodeToNodeEncryptionOptions AWS CloudFormation Resource (AWS::Elasticsearch::Domain.NodeToNodeEncryptionOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-nodetonodeencryptionoptions.html
type AWSElasticsearchDomain_NodeToNodeEncryptionOptions struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-nodetonodeencryptionoptions.html#cfn-elasticsearch-domain-nodetonodeencryptionoptions-enabled
	Enabled *Value `json:"Enabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticsearchDomain_NodeToNodeEncryptionOptions) AWSCloudFormationType() string {
	return "AWS::Elasticsearch::Domain.NodeToNodeEncryptionOptions"
}

func (r *AWSElasticsearchDomain_NodeToNodeEncryptionOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
