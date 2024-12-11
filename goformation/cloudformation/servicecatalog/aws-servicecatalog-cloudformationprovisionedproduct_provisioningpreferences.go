package servicecatalog

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// CloudFormationProvisionedProduct_ProvisioningPreferences AWS CloudFormation Resource (AWS::ServiceCatalog::CloudFormationProvisionedProduct.ProvisioningPreferences)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences.html
type CloudFormationProvisionedProduct_ProvisioningPreferences struct {

	// StackSetAccounts AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences-stacksetaccounts
	StackSetAccounts *types.Value `json:"StackSetAccounts,omitempty"`

	// StackSetFailureToleranceCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences-stacksetfailuretolerancecount
	StackSetFailureToleranceCount *types.Value `json:"StackSetFailureToleranceCount,omitempty"`

	// StackSetFailureTolerancePercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences-stacksetfailuretolerancepercentage
	StackSetFailureTolerancePercentage *types.Value `json:"StackSetFailureTolerancePercentage,omitempty"`

	// StackSetMaxConcurrencyCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences-stacksetmaxconcurrencycount
	StackSetMaxConcurrencyCount *types.Value `json:"StackSetMaxConcurrencyCount,omitempty"`

	// StackSetMaxConcurrencyPercentage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences-stacksetmaxconcurrencypercentage
	StackSetMaxConcurrencyPercentage *types.Value `json:"StackSetMaxConcurrencyPercentage,omitempty"`

	// StackSetOperationType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences-stacksetoperationtype
	StackSetOperationType *types.Value `json:"StackSetOperationType,omitempty"`

	// StackSetRegions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences.html#cfn-servicecatalog-cloudformationprovisionedproduct-provisioningpreferences-stacksetregions
	StackSetRegions *types.Value `json:"StackSetRegions,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *CloudFormationProvisionedProduct_ProvisioningPreferences) AWSCloudFormationType() string {
	return "AWS::ServiceCatalog::CloudFormationProvisionedProduct.ProvisioningPreferences"
}
