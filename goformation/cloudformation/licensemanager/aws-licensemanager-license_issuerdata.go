package licensemanager

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// License_IssuerData AWS CloudFormation Resource (AWS::LicenseManager::License.IssuerData)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-license-issuerdata.html
type License_IssuerData struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-license-issuerdata.html#cfn-licensemanager-license-issuerdata-name
	Name *types.Value `json:"Name,omitempty"`

	// SignKey AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-licensemanager-license-issuerdata.html#cfn-licensemanager-license-issuerdata-signkey
	SignKey *types.Value `json:"SignKey,omitempty"`

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
func (r *License_IssuerData) AWSCloudFormationType() string {
	return "AWS::LicenseManager::License.IssuerData"
}
