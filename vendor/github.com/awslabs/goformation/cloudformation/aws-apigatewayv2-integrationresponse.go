package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSApiGatewayV2IntegrationResponse AWS CloudFormation Resource (AWS::ApiGatewayV2::IntegrationResponse)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html
type AWSApiGatewayV2IntegrationResponse struct {

	// ApiId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-apiid
	ApiId *Value `json:"ApiId,omitempty"`

	// ContentHandlingStrategy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-contenthandlingstrategy
	ContentHandlingStrategy *Value `json:"ContentHandlingStrategy,omitempty"`

	// IntegrationId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-integrationid
	IntegrationId *Value `json:"IntegrationId,omitempty"`

	// IntegrationResponseKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-integrationresponsekey
	IntegrationResponseKey *Value `json:"IntegrationResponseKey,omitempty"`

	// ResponseParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-responseparameters
	ResponseParameters interface{} `json:"ResponseParameters,omitempty"`

	// ResponseTemplates AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-responsetemplates
	ResponseTemplates interface{} `json:"ResponseTemplates,omitempty"`

	// TemplateSelectionExpression AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-apigatewayv2-integrationresponse.html#cfn-apigatewayv2-integrationresponse-templateselectionexpression
	TemplateSelectionExpression *Value `json:"TemplateSelectionExpression,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayV2IntegrationResponse) AWSCloudFormationType() string {
	return "AWS::ApiGatewayV2::IntegrationResponse"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSApiGatewayV2IntegrationResponse) MarshalJSON() ([]byte, error) {
	type Properties AWSApiGatewayV2IntegrationResponse
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
func (r *AWSApiGatewayV2IntegrationResponse) UnmarshalJSON(b []byte) error {
	type Properties AWSApiGatewayV2IntegrationResponse
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
		*r = AWSApiGatewayV2IntegrationResponse(*res.Properties)
	}

	return nil
}

// GetAllAWSApiGatewayV2IntegrationResponseResources retrieves all AWSApiGatewayV2IntegrationResponse items from an AWS CloudFormation template
func (t *Template) GetAllAWSApiGatewayV2IntegrationResponseResources() map[string]AWSApiGatewayV2IntegrationResponse {
	results := map[string]AWSApiGatewayV2IntegrationResponse{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSApiGatewayV2IntegrationResponse:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::ApiGatewayV2::IntegrationResponse" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSApiGatewayV2IntegrationResponse{}
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

// GetAWSApiGatewayV2IntegrationResponseWithName retrieves all AWSApiGatewayV2IntegrationResponse items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSApiGatewayV2IntegrationResponseWithName(name string) (AWSApiGatewayV2IntegrationResponse, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSApiGatewayV2IntegrationResponse:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::ApiGatewayV2::IntegrationResponse" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSApiGatewayV2IntegrationResponse{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSApiGatewayV2IntegrationResponse{}, errors.New("resource not found")
}
