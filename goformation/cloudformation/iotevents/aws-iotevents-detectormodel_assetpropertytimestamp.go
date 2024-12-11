package iotevents

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// DetectorModel_AssetPropertyTimestamp AWS CloudFormation Resource (AWS::IoTEvents::DetectorModel.AssetPropertyTimestamp)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-assetpropertytimestamp.html
type DetectorModel_AssetPropertyTimestamp struct {

	// OffsetInNanos AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-assetpropertytimestamp.html#cfn-iotevents-detectormodel-assetpropertytimestamp-offsetinnanos
	OffsetInNanos *types.Value `json:"OffsetInNanos,omitempty"`

	// TimeInSeconds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotevents-detectormodel-assetpropertytimestamp.html#cfn-iotevents-detectormodel-assetpropertytimestamp-timeinseconds
	TimeInSeconds *types.Value `json:"TimeInSeconds,omitempty"`

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
func (r *DetectorModel_AssetPropertyTimestamp) AWSCloudFormationType() string {
	return "AWS::IoTEvents::DetectorModel.AssetPropertyTimestamp"
}
