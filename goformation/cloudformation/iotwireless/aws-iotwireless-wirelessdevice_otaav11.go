package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WirelessDevice_OtaaV11 AWS CloudFormation Resource (AWS::IoTWireless::WirelessDevice.OtaaV11)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-otaav11.html
type WirelessDevice_OtaaV11 struct {

	// AppKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-otaav11.html#cfn-iotwireless-wirelessdevice-otaav11-appkey
	AppKey *types.Value `json:"AppKey,omitempty"`

	// JoinEui AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-otaav11.html#cfn-iotwireless-wirelessdevice-otaav11-joineui
	JoinEui *types.Value `json:"JoinEui,omitempty"`

	// NwkKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-otaav11.html#cfn-iotwireless-wirelessdevice-otaav11-nwkkey
	NwkKey *types.Value `json:"NwkKey,omitempty"`

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
func (r *WirelessDevice_OtaaV11) AWSCloudFormationType() string {
	return "AWS::IoTWireless::WirelessDevice.OtaaV11"
}
