package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSRoute53ResolverResolverRule AWS CloudFormation Resource (AWS::Route53Resolver::ResolverRule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverrule.html
type AWSRoute53ResolverResolverRule struct {

	// DomainName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverrule.html#cfn-route53resolver-resolverrule-domainname
	DomainName *Value `json:"DomainName,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverrule.html#cfn-route53resolver-resolverrule-name
	Name *Value `json:"Name,omitempty"`

	// ResolverEndpointId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverrule.html#cfn-route53resolver-resolverrule-resolverendpointid
	ResolverEndpointId *Value `json:"ResolverEndpointId,omitempty"`

	// RuleType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverrule.html#cfn-route53resolver-resolverrule-ruletype
	RuleType *Value `json:"RuleType,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverrule.html#cfn-route53resolver-resolverrule-tags
	Tags []Tag `json:"Tags,omitempty"`

	// TargetIps AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverrule.html#cfn-route53resolver-resolverrule-targetips
	TargetIps []AWSRoute53ResolverResolverRule_TargetAddress `json:"TargetIps,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoute53ResolverResolverRule) AWSCloudFormationType() string {
	return "AWS::Route53Resolver::ResolverRule"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSRoute53ResolverResolverRule) MarshalJSON() ([]byte, error) {
	type Properties AWSRoute53ResolverResolverRule
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
func (r *AWSRoute53ResolverResolverRule) UnmarshalJSON(b []byte) error {
	type Properties AWSRoute53ResolverResolverRule
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
		*r = AWSRoute53ResolverResolverRule(*res.Properties)
	}

	return nil
}

// GetAllAWSRoute53ResolverResolverRuleResources retrieves all AWSRoute53ResolverResolverRule items from an AWS CloudFormation template
func (t *Template) GetAllAWSRoute53ResolverResolverRuleResources() map[string]AWSRoute53ResolverResolverRule {
	results := map[string]AWSRoute53ResolverResolverRule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSRoute53ResolverResolverRule:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Route53Resolver::ResolverRule" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoute53ResolverResolverRule{}
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

// GetAWSRoute53ResolverResolverRuleWithName retrieves all AWSRoute53ResolverResolverRule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSRoute53ResolverResolverRuleWithName(name string) (AWSRoute53ResolverResolverRule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSRoute53ResolverResolverRule:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Route53Resolver::ResolverRule" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoute53ResolverResolverRule{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSRoute53ResolverResolverRule{}, errors.New("resource not found")
}
