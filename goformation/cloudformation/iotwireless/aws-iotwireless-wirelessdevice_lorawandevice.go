package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WirelessDevice_LoRaWANDevice AWS CloudFormation Resource (AWS::IoTWireless::WirelessDevice.LoRaWANDevice)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-lorawandevice.html
type WirelessDevice_LoRaWANDevice struct {

	// AbpV10x AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-lorawandevice.html#cfn-iotwireless-wirelessdevice-lorawandevice-abpv10x
	AbpV10x *WirelessDevice_AbpV10x `json:"AbpV10x,omitempty"`

	// AbpV11 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-lorawandevice.html#cfn-iotwireless-wirelessdevice-lorawandevice-abpv11
	AbpV11 *WirelessDevice_AbpV11 `json:"AbpV11,omitempty"`

	// DevEui AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-lorawandevice.html#cfn-iotwireless-wirelessdevice-lorawandevice-deveui
	DevEui *types.Value `json:"DevEui,omitempty"`

	// DeviceProfileId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-lorawandevice.html#cfn-iotwireless-wirelessdevice-lorawandevice-deviceprofileid
	DeviceProfileId *types.Value `json:"DeviceProfileId,omitempty"`

	// OtaaV10x AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-lorawandevice.html#cfn-iotwireless-wirelessdevice-lorawandevice-otaav10x
	OtaaV10x *WirelessDevice_OtaaV10x `json:"OtaaV10x,omitempty"`

	// OtaaV11 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-lorawandevice.html#cfn-iotwireless-wirelessdevice-lorawandevice-otaav11
	OtaaV11 *WirelessDevice_OtaaV11 `json:"OtaaV11,omitempty"`

	// ServiceProfileId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessdevice-lorawandevice.html#cfn-iotwireless-wirelessdevice-lorawandevice-serviceprofileid
	ServiceProfileId *types.Value `json:"ServiceProfileId,omitempty"`

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
func (r *WirelessDevice_LoRaWANDevice) AWSCloudFormationType() string {
	return "AWS::IoTWireless::WirelessDevice.LoRaWANDevice"
}
