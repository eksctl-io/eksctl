package iotsitewise

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Asset_AssetProperty AWS CloudFormation Resource (AWS::IoTSiteWise::Asset.AssetProperty)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-asset-assetproperty.html
type Asset_AssetProperty struct {

	// Alias AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-asset-assetproperty.html#cfn-iotsitewise-asset-assetproperty-alias
	Alias *types.Value `json:"Alias,omitempty"`

	// LogicalId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-asset-assetproperty.html#cfn-iotsitewise-asset-assetproperty-logicalid
	LogicalId *types.Value `json:"LogicalId,omitempty"`

	// NotificationState AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotsitewise-asset-assetproperty.html#cfn-iotsitewise-asset-assetproperty-notificationstate
	NotificationState *types.Value `json:"NotificationState,omitempty"`

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
func (r *Asset_AssetProperty) AWSCloudFormationType() string {
	return "AWS::IoTSiteWise::Asset.AssetProperty"
}
