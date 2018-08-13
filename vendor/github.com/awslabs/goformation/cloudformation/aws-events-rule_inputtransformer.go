package cloudformation

import (
	"encoding/json"
)

// AWSEventsRule_InputTransformer AWS CloudFormation Resource (AWS::Events::Rule.InputTransformer)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-inputtransformer.html
type AWSEventsRule_InputTransformer struct {

	// InputPathsMap AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-inputtransformer.html#cfn-events-rule-inputtransformer-inputpathsmap
	InputPathsMap map[string]*Value `json:"InputPathsMap,omitempty"`

	// InputTemplate AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-inputtransformer.html#cfn-events-rule-inputtransformer-inputtemplate
	InputTemplate *Value `json:"InputTemplate,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEventsRule_InputTransformer) AWSCloudFormationType() string {
	return "AWS::Events::Rule.InputTransformer"
}

func (r *AWSEventsRule_InputTransformer) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
