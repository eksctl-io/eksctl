package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSRoute53ResolverResolverRuleAssociation AWS CloudFormation Resource (AWS::Route53Resolver::ResolverRuleAssociation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverruleassociation.html
type AWSRoute53ResolverResolverRuleAssociation struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverruleassociation.html#cfn-route53resolver-resolverruleassociation-name
	Name *Value `json:"Name,omitempty"`

	// ResolverRuleId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverruleassociation.html#cfn-route53resolver-resolverruleassociation-resolverruleid
	ResolverRuleId *Value `json:"ResolverRuleId,omitempty"`

	// VPCId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverruleassociation.html#cfn-route53resolver-resolverruleassociation-vpcid
	VPCId *Value `json:"VPCId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoute53ResolverResolverRuleAssociation) AWSCloudFormationType() string {
	return "AWS::Route53Resolver::ResolverRuleAssociation"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSRoute53ResolverResolverRuleAssociation) MarshalJSON() ([]byte, error) {
	type Properties AWSRoute53ResolverResolverRuleAssociation
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
func (r *AWSRoute53ResolverResolverRuleAssociation) UnmarshalJSON(b []byte) error {
	type Properties AWSRoute53ResolverResolverRuleAssociation
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
		*r = AWSRoute53ResolverResolverRuleAssociation(*res.Properties)
	}

	return nil
}

// GetAllAWSRoute53ResolverResolverRuleAssociationResources retrieves all AWSRoute53ResolverResolverRuleAssociation items from an AWS CloudFormation template
func (t *Template) GetAllAWSRoute53ResolverResolverRuleAssociationResources() map[string]AWSRoute53ResolverResolverRuleAssociation {
	results := map[string]AWSRoute53ResolverResolverRuleAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSRoute53ResolverResolverRuleAssociation:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Route53Resolver::ResolverRuleAssociation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoute53ResolverResolverRuleAssociation{}
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

// GetAWSRoute53ResolverResolverRuleAssociationWithName retrieves all AWSRoute53ResolverResolverRuleAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSRoute53ResolverResolverRuleAssociationWithName(name string) (AWSRoute53ResolverResolverRuleAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSRoute53ResolverResolverRuleAssociation:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Route53Resolver::ResolverRuleAssociation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoute53ResolverResolverRuleAssociation{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSRoute53ResolverResolverRuleAssociation{}, errors.New("resource not found")
}
