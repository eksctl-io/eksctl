package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WirelessDevice_SessionKeysAbpV10x AWS CloudFormation Resource (AWS::IoTWireless::WirelessDevice.SessionKeysAbpV10x)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-sessionkeysabpv10x.html
type WirelessDevice_SessionKeysAbpV10x struct {

	// AppSKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-sessionkeysabpv10x.html#cfn-iotwireless-wirelessdevice-sessionkeysabpv10x-appskey
	AppSKey *types.Value `json:"AppSKey,omitempty"`

	// NwkSKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-sessionkeysabpv10x.html#cfn-iotwireless-wirelessdevice-sessionkeysabpv10x-nwkskey
	NwkSKey *types.Value `json:"NwkSKey,omitempty"`

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
func (r *WirelessDevice_SessionKeysAbpV10x) AWSCloudFormationType() string {
	return "AWS::IoTWireless::WirelessDevice.SessionKeysAbpV10x"
}
