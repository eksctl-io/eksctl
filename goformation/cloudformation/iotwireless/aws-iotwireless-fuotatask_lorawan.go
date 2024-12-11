package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// FuotaTask_LoRaWAN AWS CloudFormation Resource (AWS::IoTWireless::FuotaTask.LoRaWAN)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-fuotatask-lorawan.html
type FuotaTask_LoRaWAN struct {

	// RfRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-fuotatask-lorawan.html#cfn-iotwireless-fuotatask-lorawan-rfregion
	RfRegion *types.Value `json:"RfRegion,omitempty"`

	// StartTime AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-fuotatask-lorawan.html#cfn-iotwireless-fuotatask-lorawan-starttime
	StartTime *types.Value `json:"StartTime,omitempty"`

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
func (r *FuotaTask_LoRaWAN) AWSCloudFormationType() string {
	return "AWS::IoTWireless::FuotaTask.LoRaWAN"
}
