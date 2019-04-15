package cloudformation

import (
	"encoding/json"
)

// AWSEventsRule_SqsParameters AWS CloudFormation Resource (AWS::Events::Rule.SqsParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-sqsparameters.html
type AWSEventsRule_SqsParameters struct {

	// MessageGroupId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-sqsparameters.html#cfn-events-rule-sqsparameters-messagegroupid
	MessageGroupId *Value `json:"MessageGroupId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEventsRule_SqsParameters) AWSCloudFormationType() string {
	return "AWS::Events::Rule.SqsParameters"
}

func (r *AWSEventsRule_SqsParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
