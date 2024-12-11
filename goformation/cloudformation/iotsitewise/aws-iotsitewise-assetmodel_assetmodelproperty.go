package iotsitewise

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AssetModel_AssetModelProperty AWS CloudFormation Resource (AWS::IoTSiteWise::AssetModel.AssetModelProperty)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-assetmodelproperty.html
type AssetModel_AssetModelProperty struct {

	// DataType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-assetmodelproperty.html#cfn-iotsitewise-assetmodel-assetmodelproperty-datatype
	DataType *types.Value `json:"DataType,omitempty"`

	// DataTypeSpec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-assetmodelproperty.html#cfn-iotsitewise-assetmodel-assetmodelproperty-datatypespec
	DataTypeSpec *types.Value `json:"DataTypeSpec,omitempty"`

	// LogicalId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-assetmodelproperty.html#cfn-iotsitewise-assetmodel-assetmodelproperty-logicalid
	LogicalId *types.Value `json:"LogicalId,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-assetmodelproperty.html#cfn-iotsitewise-assetmodel-assetmodelproperty-name
	Name *types.Value `json:"Name,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-assetmodelproperty.html#cfn-iotsitewise-assetmodel-assetmodelproperty-type
	Type *AssetModel_PropertyType `json:"Type,omitempty"`

	// Unit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-assetmodelproperty.html#cfn-iotsitewise-assetmodel-assetmodelproperty-unit
	Unit *types.Value `json:"Unit,omitempty"`

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
func (r *AssetModel_AssetModelProperty) AWSCloudFormationType() string {
	return "AWS::IoTSiteWise::AssetModel.AssetModelProperty"
}
