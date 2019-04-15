package cloudformation

import (
	"encoding/json"
)

// AWSAmazonMQConfiguration_TagsEntry AWS CloudFormation Resource (AWS::AmazonMQ::Configuration.TagsEntry)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configuration-tagsentry.html
type AWSAmazonMQConfiguration_TagsEntry struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configuration-tagsentry.html#cfn-amazonmq-configuration-tagsentry-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-configuration-tagsentry.html#cfn-amazonmq-configuration-tagsentry-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAmazonMQConfiguration_TagsEntry) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::Configuration.TagsEntry"
}

func (r *AWSAmazonMQConfiguration_TagsEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
