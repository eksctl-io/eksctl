package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSRoute53ResolverResolverEndpoint AWS CloudFormation Resource (AWS::Route53Resolver::ResolverEndpoint)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverendpoint.html
type AWSRoute53ResolverResolverEndpoint struct {

	// Direction AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverendpoint.html#cfn-route53resolver-resolverendpoint-direction
	Direction *Value `json:"Direction,omitempty"`

	// IpAddresses AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverendpoint.html#cfn-route53resolver-resolverendpoint-ipaddresses
	IpAddresses []AWSRoute53ResolverResolverEndpoint_IpAddressRequest `json:"IpAddresses,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverendpoint.html#cfn-route53resolver-resolverendpoint-name
	Name *Value `json:"Name,omitempty"`

	// SecurityGroupIds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverendpoint.html#cfn-route53resolver-resolverendpoint-securitygroupids
	SecurityGroupIds []*Value `json:"SecurityGroupIds,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-route53resolver-resolverendpoint.html#cfn-route53resolver-resolverendpoint-tags
	Tags []Tag `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoute53ResolverResolverEndpoint) AWSCloudFormationType() string {
	return "AWS::Route53Resolver::ResolverEndpoint"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSRoute53ResolverResolverEndpoint) MarshalJSON() ([]byte, error) {
	type Properties AWSRoute53ResolverResolverEndpoint
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
func (r *AWSRoute53ResolverResolverEndpoint) UnmarshalJSON(b []byte) error {
	type Properties AWSRoute53ResolverResolverEndpoint
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
		*r = AWSRoute53ResolverResolverEndpoint(*res.Properties)
	}

	return nil
}

// GetAllAWSRoute53ResolverResolverEndpointResources retrieves all AWSRoute53ResolverResolverEndpoint items from an AWS CloudFormation template
func (t *Template) GetAllAWSRoute53ResolverResolverEndpointResources() map[string]AWSRoute53ResolverResolverEndpoint {
	results := map[string]AWSRoute53ResolverResolverEndpoint{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSRoute53ResolverResolverEndpoint:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Route53Resolver::ResolverEndpoint" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoute53ResolverResolverEndpoint{}
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

// GetAWSRoute53ResolverResolverEndpointWithName retrieves all AWSRoute53ResolverResolverEndpoint items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSRoute53ResolverResolverEndpointWithName(name string) (AWSRoute53ResolverResolverEndpoint, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSRoute53ResolverResolverEndpoint:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Route53Resolver::ResolverEndpoint" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoute53ResolverResolverEndpoint{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSRoute53ResolverResolverEndpoint{}, errors.New("resource not found")
}
