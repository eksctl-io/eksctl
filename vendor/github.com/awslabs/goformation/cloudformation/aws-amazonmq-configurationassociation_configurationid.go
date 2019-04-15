package cloudformation

import (
	"encoding/json"
)

// AWSAmazonMQConfigurationAssociation_ConfigurationId AWS CloudFormation Resource (AWS::AmazonMQ::ConfigurationAssociation.ConfigurationId)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configurationassociation-configurationid.html
type AWSAmazonMQConfigurationAssociation_ConfigurationId struct {

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configurationassociation-configurationid.html#cfn-amazonmq-configurationassociation-configurationid-id
	Id *Value `json:"Id,omitempty"`

	// Revision AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configurationassociation-configurationid.html#cfn-amazonmq-configurationassociation-configurationid-revision
	Revision *Value `json:"Revision,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAmazonMQConfigurationAssociation_ConfigurationId) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::ConfigurationAssociation.ConfigurationId"
}

func (r *AWSAmazonMQConfigurationAssociation_ConfigurationId) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
