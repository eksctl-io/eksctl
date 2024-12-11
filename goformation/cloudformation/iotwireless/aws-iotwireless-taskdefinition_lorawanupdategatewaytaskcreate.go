package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// TaskDefinition_LoRaWANUpdateGatewayTaskCreate AWS CloudFormation Resource (AWS::IoTWireless::TaskDefinition.LoRaWANUpdateGatewayTaskCreate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate.html
type TaskDefinition_LoRaWANUpdateGatewayTaskCreate struct {

	// CurrentVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate.html#cfn-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate-currentversion
	CurrentVersion *TaskDefinition_LoRaWANGatewayVersion `json:"CurrentVersion,omitempty"`

	// SigKeyCrc AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate.html#cfn-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate-sigkeycrc
	SigKeyCrc *types.Value `json:"SigKeyCrc,omitempty"`

	// UpdateSignature AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate.html#cfn-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate-updatesignature
	UpdateSignature *types.Value `json:"UpdateSignature,omitempty"`

	// UpdateVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate.html#cfn-iotwireless-taskdefinition-lorawanupdategatewaytaskcreate-updateversion
	UpdateVersion *TaskDefinition_LoRaWANGatewayVersion `json:"UpdateVersion,omitempty"`

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
func (r *TaskDefinition_LoRaWANUpdateGatewayTaskCreate) AWSCloudFormationType() string {
	return "AWS::IoTWireless::TaskDefinition.LoRaWANUpdateGatewayTaskCreate"
}
