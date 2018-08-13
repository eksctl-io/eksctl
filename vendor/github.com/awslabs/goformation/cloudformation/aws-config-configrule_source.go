package cloudformation

import (
	"encoding/json"
)

// AWSConfigConfigRule_Source AWS CloudFormation Resource (AWS::Config::ConfigRule.Source)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-source.html
type AWSConfigConfigRule_Source struct {

	// Owner AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-source.html#cfn-config-configrule-source-owner
	Owner *Value `json:"Owner,omitempty"`

	// SourceDetails AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-source.html#cfn-config-configrule-source-sourcedetails
	SourceDetails []AWSConfigConfigRule_SourceDetail `json:"SourceDetails,omitempty"`

	// SourceIdentifier AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configrule-source.html#cfn-config-configrule-source-sourceidentifier
	SourceIdentifier *Value `json:"SourceIdentifier,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSConfigConfigRule_Source) AWSCloudFormationType() string {
	return "AWS::Config::ConfigRule.Source"
}

func (r *AWSConfigConfigRule_Source) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
