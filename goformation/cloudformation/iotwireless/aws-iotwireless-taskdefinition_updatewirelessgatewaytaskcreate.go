package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// TaskDefinition_UpdateWirelessGatewayTaskCreate AWS CloudFormation Resource (AWS::IoTWireless::TaskDefinition.UpdateWirelessGatewayTaskCreate)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-updatewirelessgatewaytaskcreate.html
type TaskDefinition_UpdateWirelessGatewayTaskCreate struct {

	// LoRaWAN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-updatewirelessgatewaytaskcreate.html#cfn-iotwireless-taskdefinition-updatewirelessgatewaytaskcreate-lorawan
	LoRaWAN *TaskDefinition_LoRaWANUpdateGatewayTaskCreate `json:"LoRaWAN,omitempty"`

	// UpdateDataRole AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-updatewirelessgatewaytaskcreate.html#cfn-iotwireless-taskdefinition-updatewirelessgatewaytaskcreate-updatedatarole
	UpdateDataRole *types.Value `json:"UpdateDataRole,omitempty"`

	// UpdateDataSource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-updatewirelessgatewaytaskcreate.html#cfn-iotwireless-taskdefinition-updatewirelessgatewaytaskcreate-updatedatasource
	UpdateDataSource *types.Value `json:"UpdateDataSource,omitempty"`

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
func (r *TaskDefinition_UpdateWirelessGatewayTaskCreate) AWSCloudFormationType() string {
	return "AWS::IoTWireless::TaskDefinition.UpdateWirelessGatewayTaskCreate"
}
