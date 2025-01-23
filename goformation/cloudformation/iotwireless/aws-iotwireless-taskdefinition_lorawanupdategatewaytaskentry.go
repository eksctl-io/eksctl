package iotwireless

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// TaskDefinition_LoRaWANUpdateGatewayTaskEntry AWS CloudFormation Resource (AWS::IoTWireless::TaskDefinition.LoRaWANUpdateGatewayTaskEntry)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawanupdategatewaytaskentry.html
type TaskDefinition_LoRaWANUpdateGatewayTaskEntry struct {

	// CurrentVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawanupdategatewaytaskentry.html#cfn-iotwireless-taskdefinition-lorawanupdategatewaytaskentry-currentversion
	CurrentVersion *TaskDefinition_LoRaWANGatewayVersion `json:"CurrentVersion,omitempty"`

	// UpdateVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iotwireless-taskdefinition-lorawanupdategatewaytaskentry.html#cfn-iotwireless-taskdefinition-lorawanupdategatewaytaskentry-updateversion
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
func (r *TaskDefinition_LoRaWANUpdateGatewayTaskEntry) AWSCloudFormationType() string {
	return "AWS::IoTWireless::TaskDefinition.LoRaWANUpdateGatewayTaskEntry"
}
