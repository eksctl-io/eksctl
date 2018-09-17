package cloudformation

import (
	"encoding/json"
)

// AWSAmazonMQBroker_ConfigurationId AWS CloudFormation Resource (AWS::AmazonMQ::Broker.ConfigurationId)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-configurationid.html
type AWSAmazonMQBroker_ConfigurationId struct {

	// Id AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-configurationid.html#cfn-amazonmq-broker-configurationid-id
	Id *Value `json:"Id,omitempty"`

	// Revision AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-configurationid.html#cfn-amazonmq-broker-configurationid-revision
	Revision *Value `json:"Revision,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAmazonMQBroker_ConfigurationId) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::Broker.ConfigurationId"
}

func (r *AWSAmazonMQBroker_ConfigurationId) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
