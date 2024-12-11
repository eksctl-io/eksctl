package iotevents

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DetectorModel_IotSiteWise AWS CloudFormation Resource (AWS::IoTEvents::DetectorModel.IotSiteWise)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-iotsitewise.html
type DetectorModel_IotSiteWise struct {

	// AssetId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-iotsitewise.html#cfn-iotevents-detectormodel-iotsitewise-assetid
	AssetId *types.Value `json:"AssetId,omitempty"`

	// EntryId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-iotsitewise.html#cfn-iotevents-detectormodel-iotsitewise-entryid
	EntryId *types.Value `json:"EntryId,omitempty"`

	// PropertyAlias AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-iotsitewise.html#cfn-iotevents-detectormodel-iotsitewise-propertyalias
	PropertyAlias *types.Value `json:"PropertyAlias,omitempty"`

	// PropertyId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-iotsitewise.html#cfn-iotevents-detectormodel-iotsitewise-propertyid
	PropertyId *types.Value `json:"PropertyId,omitempty"`

	// PropertyValue AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-iotsitewise.html#cfn-iotevents-detectormodel-iotsitewise-propertyvalue
	PropertyValue *DetectorModel_AssetPropertyValue `json:"PropertyValue,omitempty"`

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
func (r *DetectorModel_IotSiteWise) AWSCloudFormationType() string {
	return "AWS::IoTEvents::DetectorModel.IotSiteWise"
}
