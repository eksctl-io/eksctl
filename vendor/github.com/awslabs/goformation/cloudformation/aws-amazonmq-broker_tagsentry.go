package cloudformation

import (
	"encoding/json"
)

// AWSAmazonMQBroker_TagsEntry AWS CloudFormation Resource (AWS::AmazonMQ::Broker.TagsEntry)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-tagsentry.html
type AWSAmazonMQBroker_TagsEntry struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-tagsentry.html#cfn-amazonmq-broker-tagsentry-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-tagsentry.html#cfn-amazonmq-broker-tagsentry-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAmazonMQBroker_TagsEntry) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::Broker.TagsEntry"
}

func (r *AWSAmazonMQBroker_TagsEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
