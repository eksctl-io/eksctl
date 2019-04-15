package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSAppStreamDirectoryConfig AWS CloudFormation Resource (AWS::AppStream::DirectoryConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-directoryconfig.html
type AWSAppStreamDirectoryConfig struct {

	// DirectoryName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-directoryconfig.html#cfn-appstream-directoryconfig-directoryname
	DirectoryName *Value `json:"DirectoryName,omitempty"`

	// OrganizationalUnitDistinguishedNames AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-directoryconfig.html#cfn-appstream-directoryconfig-organizationalunitdistinguishednames
	OrganizationalUnitDistinguishedNames []*Value `json:"OrganizationalUnitDistinguishedNames,omitempty"`

	// ServiceAccountCredentials AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-directoryconfig.html#cfn-appstream-directoryconfig-serviceaccountcredentials
	ServiceAccountCredentials *AWSAppStreamDirectoryConfig_ServiceAccountCredentials `json:"ServiceAccountCredentials,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamDirectoryConfig) AWSCloudFormationType() string {
	return "AWS::AppStream::DirectoryConfig"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSAppStreamDirectoryConfig) MarshalJSON() ([]byte, error) {
	type Properties AWSAppStreamDirectoryConfig
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
func (r *AWSAppStreamDirectoryConfig) UnmarshalJSON(b []byte) error {
	type Properties AWSAppStreamDirectoryConfig
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
		*r = AWSAppStreamDirectoryConfig(*res.Properties)
	}

	return nil
}

// GetAllAWSAppStreamDirectoryConfigResources retrieves all AWSAppStreamDirectoryConfig items from an AWS CloudFormation template
func (t *Template) GetAllAWSAppStreamDirectoryConfigResources() map[string]AWSAppStreamDirectoryConfig {
	results := map[string]AWSAppStreamDirectoryConfig{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSAppStreamDirectoryConfig:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppStream::DirectoryConfig" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppStreamDirectoryConfig{}
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

// GetAWSAppStreamDirectoryConfigWithName retrieves all AWSAppStreamDirectoryConfig items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSAppStreamDirectoryConfigWithName(name string) (AWSAppStreamDirectoryConfig, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSAppStreamDirectoryConfig:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppStream::DirectoryConfig" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppStreamDirectoryConfig{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSAppStreamDirectoryConfig{}, errors.New("resource not found")
}
