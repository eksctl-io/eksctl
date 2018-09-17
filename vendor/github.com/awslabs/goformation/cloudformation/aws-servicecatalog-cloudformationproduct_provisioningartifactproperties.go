package cloudformation

import (
	"encoding/json"
)

// AWSServiceCatalogCloudFormationProduct_ProvisioningArtifactProperties AWS CloudFormation Resource (AWS::ServiceCatalog::CloudFormationProduct.ProvisioningArtifactProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationproduct-provisioningartifactproperties.html
type AWSServiceCatalogCloudFormationProduct_ProvisioningArtifactProperties struct {

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationproduct-provisioningartifactproperties.html#cfn-servicecatalog-cloudformationproduct-provisioningartifactproperties-description
	Description *Value `json:"Description,omitempty"`

	// Info AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationproduct-provisioningartifactproperties.html#cfn-servicecatalog-cloudformationproduct-provisioningartifactproperties-info
	Info interface{} `json:"Info,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationproduct-provisioningartifactproperties.html#cfn-servicecatalog-cloudformationproduct-provisioningartifactproperties-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSServiceCatalogCloudFormationProduct_ProvisioningArtifactProperties) AWSCloudFormationType() string {
	return "AWS::ServiceCatalog::CloudFormationProduct.ProvisioningArtifactProperties"
}

func (r *AWSServiceCatalogCloudFormationProduct_ProvisioningArtifactProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
