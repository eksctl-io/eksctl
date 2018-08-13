package cloudformation

import (
	"encoding/json"
)

// AWSAmazonMQBroker_LogList AWS CloudFormation Resource (AWS::AmazonMQ::Broker.LogList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-loglist.html
type AWSAmazonMQBroker_LogList struct {

	// Audit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-loglist.html#cfn-amazonmq-broker-loglist-audit
	Audit *Value `json:"Audit,omitempty"`

	// General AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-amazonmq-broker-loglist.html#cfn-amazonmq-broker-loglist-general
	General *Value `json:"General,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAmazonMQBroker_LogList) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::Broker.LogList"
}

func (r *AWSAmazonMQBroker_LogList) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
