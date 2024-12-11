package iotwireless

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// TaskDefinition_LoRaWANGatewayVersion AWS CloudFormation Resource (AWS::IoTWireless::TaskDefinition.LoRaWANGatewayVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawangatewayversion.html
type TaskDefinition_LoRaWANGatewayVersion struct {

	// Model AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawangatewayversion.html#cfn-iotwireless-taskdefinition-lorawangatewayversion-model
	Model *types.Value `json:"Model,omitempty"`

	// PackageVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawangatewayversion.html#cfn-iotwireless-taskdefinition-lorawangatewayversion-packageversion
	PackageVersion *types.Value `json:"PackageVersion,omitempty"`

	// Station AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawangatewayversion.html#cfn-iotwireless-taskdefinition-lorawangatewayversion-station
	Station *types.Value `json:"Station,omitempty"`

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
func (r *TaskDefinition_LoRaWANGatewayVersion) AWSCloudFormationType() string {
	return "AWS::IoTWireless::TaskDefinition.LoRaWANGatewayVersion"
}
