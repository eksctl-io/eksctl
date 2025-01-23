package mediapackage

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// PackagingConfiguration_MssEncryption AWS CloudFormation Resource (AWS::MediaPackage::PackagingConfiguration.MssEncryption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-mssencryption.html
type PackagingConfiguration_MssEncryption struct {

	// SpekeKeyProvider AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-mssencryption.html#cfn-mediapackage-packagingconfiguration-mssencryption-spekekeyprovider
	SpekeKeyProvider *PackagingConfiguration_SpekeKeyProvider `json:"SpekeKeyProvider,omitempty"`

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
func (r *PackagingConfiguration_MssEncryption) AWSCloudFormationType() string {
	return "AWS::MediaPackage::PackagingConfiguration.MssEncryption"
}
