package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// WirelessGateway_LoRaWANGateway AWS CloudFormation Resource (AWS::IoTWireless::WirelessGateway.LoRaWANGateway)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessgateway-lorawangateway.html
type WirelessGateway_LoRaWANGateway struct {

	// GatewayEui AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessgateway-lorawangateway.html#cfn-iotwireless-wirelessgateway-lorawangateway-gatewayeui
	GatewayEui *types.Value `json:"GatewayEui,omitempty"`

	// RfRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-wirelessgateway-lorawangateway.html#cfn-iotwireless-wirelessgateway-lorawangateway-rfregion
	RfRegion *types.Value `json:"RfRegion,omitempty"`

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
func (r *WirelessGateway_LoRaWANGateway) AWSCloudFormationType() string {
	return "AWS::IoTWireless::WirelessGateway.LoRaWANGateway"
}
