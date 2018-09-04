package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSServiceCatalogPortfolioPrincipalAssociation AWS CloudFormation Resource (AWS::ServiceCatalog::PortfolioPrincipalAssociation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-servicecatalog-portfolioprincipalassociation.html
type AWSServiceCatalogPortfolioPrincipalAssociation struct {

	// AcceptLanguage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-servicecatalog-portfolioprincipalassociation.html#cfn-servicecatalog-portfolioprincipalassociation-acceptlanguage
	AcceptLanguage *StringIntrinsic `json:"AcceptLanguage,omitempty"`

	// PortfolioId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-servicecatalog-portfolioprincipalassociation.html#cfn-servicecatalog-portfolioprincipalassociation-portfolioid
	PortfolioId *StringIntrinsic `json:"PortfolioId,omitempty"`

	// PrincipalARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-servicecatalog-portfolioprincipalassociation.html#cfn-servicecatalog-portfolioprincipalassociation-principalarn
	PrincipalARN *StringIntrinsic `json:"PrincipalARN,omitempty"`

	// PrincipalType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-servicecatalog-portfolioprincipalassociation.html#cfn-servicecatalog-portfolioprincipalassociation-principaltype
	PrincipalType *StringIntrinsic `json:"PrincipalType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServiceCatalogPortfolioPrincipalAssociation) AWSCloudFormationType() string {
	return "AWS::ServiceCatalog::PortfolioPrincipalAssociation"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSServiceCatalogPortfolioPrincipalAssociation) MarshalJSON() ([]byte, error) {
	type Properties AWSServiceCatalogPortfolioPrincipalAssociation
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
func (r *AWSServiceCatalogPortfolioPrincipalAssociation) UnmarshalJSON(b []byte) error {
	type Properties AWSServiceCatalogPortfolioPrincipalAssociation
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
		*r = AWSServiceCatalogPortfolioPrincipalAssociation(*res.Properties)
	}

	return nil
}

// GetAllAWSServiceCatalogPortfolioPrincipalAssociationResources retrieves all AWSServiceCatalogPortfolioPrincipalAssociation items from an AWS CloudFormation template
func (t *Template) GetAllAWSServiceCatalogPortfolioPrincipalAssociationResources() map[string]AWSServiceCatalogPortfolioPrincipalAssociation {
	results := map[string]AWSServiceCatalogPortfolioPrincipalAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSServiceCatalogPortfolioPrincipalAssociation:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::ServiceCatalog::PortfolioPrincipalAssociation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSServiceCatalogPortfolioPrincipalAssociation
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

// GetAWSServiceCatalogPortfolioPrincipalAssociationWithName retrieves all AWSServiceCatalogPortfolioPrincipalAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSServiceCatalogPortfolioPrincipalAssociationWithName(name string) (AWSServiceCatalogPortfolioPrincipalAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSServiceCatalogPortfolioPrincipalAssociation:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::ServiceCatalog::PortfolioPrincipalAssociation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSServiceCatalogPortfolioPrincipalAssociation
						if err := json.Unmarshal(b, &result); err == nil {
							return result, nil
						}
					}
				}
			}
		}
	}
	return AWSServiceCatalogPortfolioPrincipalAssociation{}, errors.New("resource not found")
}
