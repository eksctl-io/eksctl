package glue

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataCatalogEncryptionSettings_EncryptionAtRest AWS CloudFormation Resource (AWS::Glue::DataCatalogEncryptionSettings.EncryptionAtRest)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-datacatalogencryptionsettings-encryptionatrest.html
type DataCatalogEncryptionSettings_EncryptionAtRest struct {

	// CatalogEncryptionMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-datacatalogencryptionsettings-encryptionatrest.html#cfn-glue-datacatalogencryptionsettings-encryptionatrest-catalogencryptionmode
	CatalogEncryptionMode *types.Value `json:"CatalogEncryptionMode,omitempty"`

	// SseAwsKmsKeyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-datacatalogencryptionsettings-encryptionatrest.html#cfn-glue-datacatalogencryptionsettings-encryptionatrest-sseawskmskeyid
	SseAwsKmsKeyId *types.Value `json:"SseAwsKmsKeyId,omitempty"`

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
func (r *DataCatalogEncryptionSettings_EncryptionAtRest) AWSCloudFormationType() string {
	return "AWS::Glue::DataCatalogEncryptionSettings.EncryptionAtRest"
}
