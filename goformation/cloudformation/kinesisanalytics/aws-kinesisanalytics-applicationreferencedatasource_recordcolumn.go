package kinesisanalytics

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ApplicationReferenceDataSource_RecordColumn AWS CloudFormation Resource (AWS::KinesisAnalytics::ApplicationReferenceDataSource.RecordColumn)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-recordcolumn.html
type ApplicationReferenceDataSource_RecordColumn struct {

	// Mapping AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalytics-applicationreferencedatasource-recordcolumn-mapping
	Mapping *types.Value `json:"Mapping,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalytics-applicationreferencedatasource-recordcolumn-name
	Name *types.Value `json:"Name,omitempty"`

	// SqlType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-kinesisanalytics-applicationreferencedatasource-recordcolumn.html#cfn-kinesisanalytics-applicationreferencedatasource-recordcolumn-sqltype
	SqlType *types.Value `json:"SqlType,omitempty"`

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
func (r *ApplicationReferenceDataSource_RecordColumn) AWSCloudFormationType() string {
	return "AWS::KinesisAnalytics::ApplicationReferenceDataSource.RecordColumn"
}
