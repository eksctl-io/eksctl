package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSDLMLifecyclePolicy AWS CloudFormation Resource (AWS::DLM::LifecyclePolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dlm-lifecyclepolicy.html
type AWSDLMLifecyclePolicy struct {

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dlm-lifecyclepolicy.html#cfn-dlm-lifecyclepolicy-description
	Description *Value `json:"Description,omitempty"`

	// ExecutionRoleArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dlm-lifecyclepolicy.html#cfn-dlm-lifecyclepolicy-executionrolearn
	ExecutionRoleArn *Value `json:"ExecutionRoleArn,omitempty"`

	// PolicyDetails AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dlm-lifecyclepolicy.html#cfn-dlm-lifecyclepolicy-policydetails
	PolicyDetails *AWSDLMLifecyclePolicy_PolicyDetails `json:"PolicyDetails,omitempty"`

	// State AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dlm-lifecyclepolicy.html#cfn-dlm-lifecyclepolicy-state
	State *Value `json:"State,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSDLMLifecyclePolicy) AWSCloudFormationType() string {
	return "AWS::DLM::LifecyclePolicy"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSDLMLifecyclePolicy) MarshalJSON() ([]byte, error) {
	type Properties AWSDLMLifecyclePolicy
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
func (r *AWSDLMLifecyclePolicy) UnmarshalJSON(b []byte) error {
	type Properties AWSDLMLifecyclePolicy
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
		*r = AWSDLMLifecyclePolicy(*res.Properties)
	}

	return nil
}

// GetAllAWSDLMLifecyclePolicyResources retrieves all AWSDLMLifecyclePolicy items from an AWS CloudFormation template
func (t *Template) GetAllAWSDLMLifecyclePolicyResources() map[string]AWSDLMLifecyclePolicy {
	results := map[string]AWSDLMLifecyclePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSDLMLifecyclePolicy:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::DLM::LifecyclePolicy" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSDLMLifecyclePolicy{}
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

// GetAWSDLMLifecyclePolicyWithName retrieves all AWSDLMLifecyclePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSDLMLifecyclePolicyWithName(name string) (AWSDLMLifecyclePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSDLMLifecyclePolicy:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::DLM::LifecyclePolicy" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSDLMLifecyclePolicy{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSDLMLifecyclePolicy{}, errors.New("resource not found")
}
