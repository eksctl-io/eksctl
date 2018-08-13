package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSApiGatewayAuthorizer AWS CloudFormation Resource (AWS::ApiGateway::Authorizer)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html
type AWSApiGatewayAuthorizer struct {

	// AuthType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-authtype
	AuthType *Value `json:"AuthType,omitempty"`

	// AuthorizerCredentials AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-authorizercredentials
	AuthorizerCredentials *Value `json:"AuthorizerCredentials,omitempty"`

	// AuthorizerResultTtlInSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-authorizerresultttlinseconds
	AuthorizerResultTtlInSeconds *Value `json:"AuthorizerResultTtlInSeconds,omitempty"`

	// AuthorizerUri AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-authorizeruri
	AuthorizerUri *Value `json:"AuthorizerUri,omitempty"`

	// IdentitySource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-identitysource
	IdentitySource *Value `json:"IdentitySource,omitempty"`

	// IdentityValidationExpression AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-identityvalidationexpression
	IdentityValidationExpression *Value `json:"IdentityValidationExpression,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-name
	Name *Value `json:"Name,omitempty"`

	// ProviderARNs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-providerarns
	ProviderARNs []*Value `json:"ProviderARNs,omitempty"`

	// RestApiId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-restapiid
	RestApiId *Value `json:"RestApiId,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigateway-authorizer.html#cfn-apigateway-authorizer-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayAuthorizer) AWSCloudFormationType() string {
	return "AWS::ApiGateway::Authorizer"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSApiGatewayAuthorizer) MarshalJSON() ([]byte, error) {
	type Properties AWSApiGatewayAuthorizer
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
func (r *AWSApiGatewayAuthorizer) UnmarshalJSON(b []byte) error {
	type Properties AWSApiGatewayAuthorizer
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
		*r = AWSApiGatewayAuthorizer(*res.Properties)
	}

	return nil
}

// GetAllAWSApiGatewayAuthorizerResources retrieves all AWSApiGatewayAuthorizer items from an AWS CloudFormation template
func (t *Template) GetAllAWSApiGatewayAuthorizerResources() map[string]AWSApiGatewayAuthorizer {
	results := map[string]AWSApiGatewayAuthorizer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSApiGatewayAuthorizer:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::ApiGateway::Authorizer" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSApiGatewayAuthorizer{}
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

// GetAWSApiGatewayAuthorizerWithName retrieves all AWSApiGatewayAuthorizer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSApiGatewayAuthorizerWithName(name string) (AWSApiGatewayAuthorizer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSApiGatewayAuthorizer:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::ApiGateway::Authorizer" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSApiGatewayAuthorizer{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSApiGatewayAuthorizer{}, errors.New("resource not found")
}
