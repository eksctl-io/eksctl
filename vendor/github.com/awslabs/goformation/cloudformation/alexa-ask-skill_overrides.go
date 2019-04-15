package cloudformation

import (
	"encoding/json"
)

// AlexaASKSkill_Overrides AWS CloudFormation Resource (Alexa::ASK::Skill.Overrides)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-overrides.html
type AlexaASKSkill_Overrides struct {

	// Manifest AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-overrides.html#cfn-ask-skill-overrides-manifest
	Manifest interface{} `json:"Manifest,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AlexaASKSkill_Overrides) AWSCloudFormationType() string {
	return "Alexa::ASK::Skill.Overrides"
}

func (r *AlexaASKSkill_Overrides) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
