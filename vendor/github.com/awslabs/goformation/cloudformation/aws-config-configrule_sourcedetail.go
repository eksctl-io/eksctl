package cloudformation

import (
	"encoding/json"
)

// AWSConfigConfigRule_SourceDetail AWS CloudFormation Resource (AWS::Config::ConfigRule.SourceDetail)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-source-sourcedetails.html
type AWSConfigConfigRule_SourceDetail struct {

	// EventSource AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-source-sourcedetails.html#cfn-config-configrule-source-sourcedetail-eventsource
	EventSource *Value `json:"EventSource,omitempty"`

	// MaximumExecutionFrequency AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-source-sourcedetails.html#cfn-config-configrule-sourcedetail-maximumexecutionfrequency
	MaximumExecutionFrequency *Value `json:"MaximumExecutionFrequency,omitempty"`

	// MessageType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-source-sourcedetails.html#cfn-config-configrule-source-sourcedetail-messagetype
	MessageType *Value `json:"MessageType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSConfigConfigRule_SourceDetail) AWSCloudFormationType() string {
	return "AWS::Config::ConfigRule.SourceDetail"
}

func (r *AWSConfigConfigRule_SourceDetail) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
