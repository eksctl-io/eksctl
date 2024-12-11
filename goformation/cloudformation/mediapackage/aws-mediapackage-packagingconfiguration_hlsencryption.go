package mediapackage

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// PackagingConfiguration_HlsEncryption AWS CloudFormation Resource (AWS::MediaPackage::PackagingConfiguration.HlsEncryption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsencryption.html
type PackagingConfiguration_HlsEncryption struct {

	// ConstantInitializationVector AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsencryption.html#cfn-mediapackage-packagingconfiguration-hlsencryption-constantinitializationvector
	ConstantInitializationVector *types.Value `json:"ConstantInitializationVector,omitempty"`

	// EncryptionMethod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsencryption.html#cfn-mediapackage-packagingconfiguration-hlsencryption-encryptionmethod
	EncryptionMethod *types.Value `json:"EncryptionMethod,omitempty"`

	// SpekeKeyProvider AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-packagingconfiguration-hlsencryption.html#cfn-mediapackage-packagingconfiguration-hlsencryption-spekekeyprovider
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
func (r *PackagingConfiguration_HlsEncryption) AWSCloudFormationType() string {
	return "AWS::MediaPackage::PackagingConfiguration.HlsEncryption"
}
