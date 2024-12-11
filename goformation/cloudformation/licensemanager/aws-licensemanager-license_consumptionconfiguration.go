package licensemanager

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// License_ConsumptionConfiguration AWS CloudFormation Resource (AWS::LicenseManager::License.ConsumptionConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-license-consumptionconfiguration.html
type License_ConsumptionConfiguration struct {

	// BorrowConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-license-consumptionconfiguration.html#cfn-licensemanager-license-consumptionconfiguration-borrowconfiguration
	BorrowConfiguration *License_BorrowConfiguration `json:"BorrowConfiguration,omitempty"`

	// ProvisionalConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-license-consumptionconfiguration.html#cfn-licensemanager-license-consumptionconfiguration-provisionalconfiguration
	ProvisionalConfiguration *License_ProvisionalConfiguration `json:"ProvisionalConfiguration,omitempty"`

	// RenewType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-license-consumptionconfiguration.html#cfn-licensemanager-license-consumptionconfiguration-renewtype
	RenewType *types.Value `json:"RenewType,omitempty"`

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
func (r *License_ConsumptionConfiguration) AWSCloudFormationType() string {
	return "AWS::LicenseManager::License.ConsumptionConfiguration"
}
