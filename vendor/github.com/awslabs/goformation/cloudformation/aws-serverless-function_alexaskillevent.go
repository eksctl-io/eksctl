package cloudformation

import (
	"encoding/json"
)

// AWSServerlessFunction_AlexaSkillEvent AWS CloudFormation Resource (AWS::Serverless::Function.AlexaSkillEvent)
// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#alexaskill
type AWSServerlessFunction_AlexaSkillEvent struct {

	// Variables AWS CloudFormation Property
	// Required: false
	// See: https://github.com/awslabs/serverless-application-model/blob/master/versions/2016-10-31.md#alexaskill
	Variables map[string]*Value `json:"Variables,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServerlessFunction_AlexaSkillEvent) AWSCloudFormationType() string {
	return "AWS::Serverless::Function.AlexaSkillEvent"
}

func (r *AWSServerlessFunction_AlexaSkillEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
