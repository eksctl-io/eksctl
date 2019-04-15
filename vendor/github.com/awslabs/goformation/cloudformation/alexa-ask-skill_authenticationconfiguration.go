package cloudformation

import (
	"encoding/json"
)

// AlexaASKSkill_AuthenticationConfiguration AWS CloudFormation Resource (Alexa::ASK::Skill.AuthenticationConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-authenticationconfiguration.html
type AlexaASKSkill_AuthenticationConfiguration struct {

	// ClientId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-authenticationconfiguration.html#cfn-ask-skill-authenticationconfiguration-clientid
	ClientId *Value `json:"ClientId,omitempty"`

	// ClientSecret AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-authenticationconfiguration.html#cfn-ask-skill-authenticationconfiguration-clientsecret
	ClientSecret *Value `json:"ClientSecret,omitempty"`

	// RefreshToken AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ask-skill-authenticationconfiguration.html#cfn-ask-skill-authenticationconfiguration-refreshtoken
	RefreshToken *Value `json:"RefreshToken,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AlexaASKSkill_AuthenticationConfiguration) AWSCloudFormationType() string {
	return "Alexa::ASK::Skill.AuthenticationConfiguration"
}

func (r *AlexaASKSkill_AuthenticationConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
