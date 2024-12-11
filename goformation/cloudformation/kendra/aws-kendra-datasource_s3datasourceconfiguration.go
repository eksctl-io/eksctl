package kendra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DataSource_S3DataSourceConfiguration AWS CloudFormation Resource (AWS::Kendra::DataSource.S3DataSourceConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-s3datasourceconfiguration.html
type DataSource_S3DataSourceConfiguration struct {

	// AccessControlListConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-s3datasourceconfiguration.html#cfn-kendra-datasource-s3datasourceconfiguration-accesscontrollistconfiguration
	AccessControlListConfiguration *DataSource_AccessControlListConfiguration `json:"AccessControlListConfiguration,omitempty"`

	// BucketName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-s3datasourceconfiguration.html#cfn-kendra-datasource-s3datasourceconfiguration-bucketname
	BucketName *types.Value `json:"BucketName,omitempty"`

	// DocumentsMetadataConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-s3datasourceconfiguration.html#cfn-kendra-datasource-s3datasourceconfiguration-documentsmetadataconfiguration
	DocumentsMetadataConfiguration *DataSource_DocumentsMetadataConfiguration `json:"DocumentsMetadataConfiguration,omitempty"`

	// ExclusionPatterns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-s3datasourceconfiguration.html#cfn-kendra-datasource-s3datasourceconfiguration-exclusionpatterns
	ExclusionPatterns *types.Value `json:"ExclusionPatterns,omitempty"`

	// InclusionPatterns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-s3datasourceconfiguration.html#cfn-kendra-datasource-s3datasourceconfiguration-inclusionpatterns
	InclusionPatterns *types.Value `json:"InclusionPatterns,omitempty"`

	// InclusionPrefixes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kendra-datasource-s3datasourceconfiguration.html#cfn-kendra-datasource-s3datasourceconfiguration-inclusionprefixes
	InclusionPrefixes *types.Value `json:"InclusionPrefixes,omitempty"`

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
func (r *DataSource_S3DataSourceConfiguration) AWSCloudFormationType() string {
	return "AWS::Kendra::DataSource.S3DataSourceConfiguration"
}
