package cloudformation

// AWSElasticsearchDomain_EncryptionAtRestOptions AWS CloudFormation Resource (AWS::Elasticsearch::Domain.EncryptionAtRestOptions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-encryptionatrestoptions.html
type AWSElasticsearchDomain_EncryptionAtRestOptions struct {

	// Enabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-encryptionatrestoptions.html#cfn-elasticsearch-domain-encryptionatrestoptions-enabled
	Enabled bool `json:"Enabled,omitempty"`

	// KmsKeyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticsearch-domain-encryptionatrestoptions.html#cfn-elasticsearch-domain-encryptionatrestoptions-kmskeyid
	KmsKeyId *StringIntrinsic `json:"KmsKeyId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSElasticsearchDomain_EncryptionAtRestOptions) AWSCloudFormationType() string {
	return "AWS::Elasticsearch::Domain.EncryptionAtRestOptions"
}
