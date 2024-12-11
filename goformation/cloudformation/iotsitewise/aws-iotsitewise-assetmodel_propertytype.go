package iotsitewise

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AssetModel_PropertyType AWS CloudFormation Resource (AWS::IoTSiteWise::AssetModel.PropertyType)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-propertytype.html
type AssetModel_PropertyType struct {

	// Attribute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-propertytype.html#cfn-iotsitewise-assetmodel-propertytype-attribute
	Attribute *AssetModel_Attribute `json:"Attribute,omitempty"`

	// Metric AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-propertytype.html#cfn-iotsitewise-assetmodel-propertytype-metric
	Metric *AssetModel_Metric `json:"Metric,omitempty"`

	// Transform AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-propertytype.html#cfn-iotsitewise-assetmodel-propertytype-transform
	Transform *AssetModel_Transform `json:"Transform,omitempty"`

	// TypeName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-propertytype.html#cfn-iotsitewise-assetmodel-propertytype-typename
	TypeName *types.Value `json:"TypeName,omitempty"`

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
func (r *AssetModel_PropertyType) AWSCloudFormationType() string {
	return "AWS::IoTSiteWise::AssetModel.PropertyType"
}
