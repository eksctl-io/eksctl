package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AlexaASKSkill AWS CloudFormation Resource (Alexa::ASK::Skill)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ask-skill.html
type AlexaASKSkill struct {

	// AuthenticationConfiguration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ask-skill.html#cfn-ask-skill-authenticationconfiguration
	AuthenticationConfiguration *AlexaASKSkill_AuthenticationConfiguration `json:"AuthenticationConfiguration,omitempty"`

	// SkillPackage AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ask-skill.html#cfn-ask-skill-skillpackage
	SkillPackage *AlexaASKSkill_SkillPackage `json:"SkillPackage,omitempty"`

	// VendorId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ask-skill.html#cfn-ask-skill-vendorid
	VendorId *Value `json:"VendorId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AlexaASKSkill) AWSCloudFormationType() string {
	return "Alexa::ASK::Skill"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AlexaASKSkill) MarshalJSON() ([]byte, error) {
	type Properties AlexaASKSkill
	return json.Marshal(&struct {
		Type       string
		Properties Properties
	}{
		Type:       r.AWSCloudFormationType(),
		Properties: (Properties)(*r),
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that strips the outer
// AWS CloudFormation resource object, and just keeps the 'Properties' field.
func (r *AlexaASKSkill) UnmarshalJSON(b []byte) error {
	type Properties AlexaASKSkill
	res := &struct {
		Type       string
		Properties *Properties
	}{}
	if err := json.Unmarshal(b, &res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		*r = AlexaASKSkill(*res.Properties)
	}

	return nil
}

// GetAllAlexaASKSkillResources retrieves all AlexaASKSkill items from an AWS CloudFormation template
func (t *Template) GetAllAlexaASKSkillResources() map[string]AlexaASKSkill {
	results := map[string]AlexaASKSkill{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AlexaASKSkill:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "Alexa::ASK::Skill" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AlexaASKSkill{}
						if err := result.UnmarshalJSON(b); err == nil {
							results[name] = *result
						}
					}
				}
			}
		}
	}
	return results
}

// GetAlexaASKSkillWithName retrieves all AlexaASKSkill items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAlexaASKSkillWithName(name string) (AlexaASKSkill, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AlexaASKSkill:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "Alexa::ASK::Skill" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AlexaASKSkill{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AlexaASKSkill{}, errors.New("resource not found")
}
