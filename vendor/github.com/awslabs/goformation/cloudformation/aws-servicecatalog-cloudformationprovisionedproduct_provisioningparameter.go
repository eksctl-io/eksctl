package cloudformation

import (
	"encoding/json"
)

// AWSServiceCatalogCloudFormationProvisionedProduct_ProvisioningParameter AWS CloudFormation Resource (AWS::ServiceCatalog::CloudFormationProvisionedProduct.ProvisioningParameter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningparameter.html
type AWSServiceCatalogCloudFormationProvisionedProduct_ProvisioningParameter struct {

	// Key AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningparameter.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningparameter-key
	Key *Value `json:"Key,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningparameter.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningparameter-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServiceCatalogCloudFormationProvisionedProduct_ProvisioningParameter) AWSCloudFormationType() string {
	return "AWS::ServiceCatalog::CloudFormationProvisionedProduct.ProvisioningParameter"
}

func (r *AWSServiceCatalogCloudFormationProvisionedProduct_ProvisioningParameter) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
