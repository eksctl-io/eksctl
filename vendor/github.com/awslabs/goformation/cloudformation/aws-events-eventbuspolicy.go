package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSEventsEventBusPolicy AWS CloudFormation Resource (AWS::Events::EventBusPolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-events-eventbuspolicy.html
type AWSEventsEventBusPolicy struct {

	// Action AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-events-eventbuspolicy.html#cfn-events-eventbuspolicy-action
	Action *Value `json:"Action,omitempty"`

	// Condition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-events-eventbuspolicy.html#cfn-events-eventbuspolicy-condition
	Condition *AWSEventsEventBusPolicy_Condition `json:"Condition,omitempty"`

	// Principal AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-events-eventbuspolicy.html#cfn-events-eventbuspolicy-principal
	Principal *Value `json:"Principal,omitempty"`

	// StatementId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-events-eventbuspolicy.html#cfn-events-eventbuspolicy-statementid
	StatementId *Value `json:"StatementId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEventsEventBusPolicy) AWSCloudFormationType() string {
	return "AWS::Events::EventBusPolicy"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSEventsEventBusPolicy) MarshalJSON() ([]byte, error) {
	type Properties AWSEventsEventBusPolicy
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
func (r *AWSEventsEventBusPolicy) UnmarshalJSON(b []byte) error {
	type Properties AWSEventsEventBusPolicy
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
		*r = AWSEventsEventBusPolicy(*res.Properties)
	}

	return nil
}

// GetAllAWSEventsEventBusPolicyResources retrieves all AWSEventsEventBusPolicy items from an AWS CloudFormation template
func (t *Template) GetAllAWSEventsEventBusPolicyResources() map[string]AWSEventsEventBusPolicy {
	results := map[string]AWSEventsEventBusPolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSEventsEventBusPolicy:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Events::EventBusPolicy" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEventsEventBusPolicy{}
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

// GetAWSEventsEventBusPolicyWithName retrieves all AWSEventsEventBusPolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSEventsEventBusPolicyWithName(name string) (AWSEventsEventBusPolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSEventsEventBusPolicy:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Events::EventBusPolicy" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEventsEventBusPolicy{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSEventsEventBusPolicy{}, errors.New("resource not found")
}
