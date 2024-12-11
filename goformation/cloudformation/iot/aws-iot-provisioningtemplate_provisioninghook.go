package iot

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ProvisioningTemplate_ProvisioningHook AWS CloudFormation Resource (AWS::IoT::ProvisioningTemplate.ProvisioningHook)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-provisioningtemplate-provisioninghook.html
type ProvisioningTemplate_ProvisioningHook struct {

	// PayloadVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-provisioningtemplate-provisioninghook.html#cfn-iot-provisioningtemplate-provisioninghook-payloadversion
	PayloadVersion *types.Value `json:"PayloadVersion,omitempty"`

	// TargetArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-provisioningtemplate-provisioninghook.html#cfn-iot-provisioningtemplate-provisioninghook-targetarn
	TargetArn *types.Value `json:"TargetArn,omitempty"`

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
func (r *ProvisioningTemplate_ProvisioningHook) AWSCloudFormationType() string {
	return "AWS::IoT::ProvisioningTemplate.ProvisioningHook"
}
