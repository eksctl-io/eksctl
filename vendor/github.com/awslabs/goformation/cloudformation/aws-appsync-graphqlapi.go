package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSAppSyncGraphQLApi AWS CloudFormation Resource (AWS::AppSync::GraphQLApi)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appsync-graphqlapi.html
type AWSAppSyncGraphQLApi struct {

	// AuthenticationType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appsync-graphqlapi.html#cfn-appsync-graphqlapi-authenticationtype
	AuthenticationType *StringIntrinsic `json:"AuthenticationType,omitempty"`

	// LogConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appsync-graphqlapi.html#cfn-appsync-graphqlapi-logconfig
	LogConfig *AWSAppSyncGraphQLApi_LogConfig `json:"LogConfig,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appsync-graphqlapi.html#cfn-appsync-graphqlapi-name
	Name *StringIntrinsic `json:"Name,omitempty"`

	// OpenIDConnectConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appsync-graphqlapi.html#cfn-appsync-graphqlapi-openidconnectconfig
	OpenIDConnectConfig *AWSAppSyncGraphQLApi_OpenIDConnectConfig `json:"OpenIDConnectConfig,omitempty"`

	// UserPoolConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appsync-graphqlapi.html#cfn-appsync-graphqlapi-userpoolconfig
	UserPoolConfig *AWSAppSyncGraphQLApi_UserPoolConfig `json:"UserPoolConfig,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppSyncGraphQLApi) AWSCloudFormationType() string {
	return "AWS::AppSync::GraphQLApi"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSAppSyncGraphQLApi) MarshalJSON() ([]byte, error) {
	type Properties AWSAppSyncGraphQLApi
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
func (r *AWSAppSyncGraphQLApi) UnmarshalJSON(b []byte) error {
	type Properties AWSAppSyncGraphQLApi
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
		*r = AWSAppSyncGraphQLApi(*res.Properties)
	}

	return nil
}

// GetAllAWSAppSyncGraphQLApiResources retrieves all AWSAppSyncGraphQLApi items from an AWS CloudFormation template
func (t *Template) GetAllAWSAppSyncGraphQLApiResources() map[string]AWSAppSyncGraphQLApi {
	results := map[string]AWSAppSyncGraphQLApi{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSAppSyncGraphQLApi:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppSync::GraphQLApi" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSAppSyncGraphQLApi
						if err := json.Unmarshal(b, &result); err == nil {
							results[name] = result
						}
					}
				}
			}
		}
	}
	return results
}

// GetAWSAppSyncGraphQLApiWithName retrieves all AWSAppSyncGraphQLApi items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSAppSyncGraphQLApiWithName(name string) (AWSAppSyncGraphQLApi, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSAppSyncGraphQLApi:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppSync::GraphQLApi" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSAppSyncGraphQLApi
						if err := json.Unmarshal(b, &result); err == nil {
							return result, nil
						}
					}
				}
			}
		}
	}
	return AWSAppSyncGraphQLApi{}, errors.New("resource not found")
}
