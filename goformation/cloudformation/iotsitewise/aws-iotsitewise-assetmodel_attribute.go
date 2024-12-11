package iotsitewise

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AssetModel_Attribute AWS CloudFormation Resource (AWS::IoTSiteWise::AssetModel.Attribute)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-attribute.html
type AssetModel_Attribute struct {

	// DefaultValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-assetmodel-attribute.html#cfn-iotsitewise-assetmodel-attribute-defaultvalue
	DefaultValue *types.Value `json:"DefaultValue,omitempty"`

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
func (r *AssetModel_Attribute) AWSCloudFormationType() string {
	return "AWS::IoTSiteWise::AssetModel.Attribute"
}
