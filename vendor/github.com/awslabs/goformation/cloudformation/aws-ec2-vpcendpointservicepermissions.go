package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSEC2VPCEndpointServicePermissions AWS CloudFormation Resource (AWS::EC2::VPCEndpointServicePermissions)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-vpcendpointservicepermissions.html
type AWSEC2VPCEndpointServicePermissions struct {

	// AllowedPrincipals AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-vpcendpointservicepermissions.html#cfn-ec2-vpcendpointservicepermissions-allowedprincipals
	AllowedPrincipals []*Value `json:"AllowedPrincipals,omitempty"`

	// ServiceId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-vpcendpointservicepermissions.html#cfn-ec2-vpcendpointservicepermissions-serviceid
	ServiceId *Value `json:"ServiceId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2VPCEndpointServicePermissions) AWSCloudFormationType() string {
	return "AWS::EC2::VPCEndpointServicePermissions"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSEC2VPCEndpointServicePermissions) MarshalJSON() ([]byte, error) {
	type Properties AWSEC2VPCEndpointServicePermissions
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
func (r *AWSEC2VPCEndpointServicePermissions) UnmarshalJSON(b []byte) error {
	type Properties AWSEC2VPCEndpointServicePermissions
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
		*r = AWSEC2VPCEndpointServicePermissions(*res.Properties)
	}

	return nil
}

// GetAllAWSEC2VPCEndpointServicePermissionsResources retrieves all AWSEC2VPCEndpointServicePermissions items from an AWS CloudFormation template
func (t *Template) GetAllAWSEC2VPCEndpointServicePermissionsResources() map[string]AWSEC2VPCEndpointServicePermissions {
	results := map[string]AWSEC2VPCEndpointServicePermissions{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSEC2VPCEndpointServicePermissions:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::EC2::VPCEndpointServicePermissions" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEC2VPCEndpointServicePermissions{}
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

// GetAWSEC2VPCEndpointServicePermissionsWithName retrieves all AWSEC2VPCEndpointServicePermissions items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSEC2VPCEndpointServicePermissionsWithName(name string) (AWSEC2VPCEndpointServicePermissions, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSEC2VPCEndpointServicePermissions:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::EC2::VPCEndpointServicePermissions" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEC2VPCEndpointServicePermissions{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSEC2VPCEndpointServicePermissions{}, errors.New("resource not found")
}
