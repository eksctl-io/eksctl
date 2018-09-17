package cloudformation

import (
	"encoding/json"
)

// AWSSNSTopic_Subscription AWS CloudFormation Resource (AWS::SNS::Topic.Subscription)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sns-subscription.html
type AWSSNSTopic_Subscription struct {

	// Endpoint AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sns-subscription.html#cfn-sns-topic-subscription-endpoint
	Endpoint *Value `json:"Endpoint,omitempty"`

	// Protocol AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sns-subscription.html#cfn-sns-topic-subscription-protocol
	Protocol *Value `json:"Protocol,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSNSTopic_Subscription) AWSCloudFormationType() string {
	return "AWS::SNS::Topic.Subscription"
}

func (r *AWSSNSTopic_Subscription) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
